package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/hosseinasadian/chat-application/pkg/richerror"
	"github.com/hosseinasadian/chat-application/service/authentication/repository"
)

type Service struct {
	Cache  repository.Cache
	Config Config
}

func New(config Config, cache repository.Cache) Service {
	return Service{Cache: cache, Config: config}
}

func (s Service) SendOtp(req SendOtpRequest) (SendOtpResponse, error) {
	const op = "authentication.service.SendOtp"

	otp := fmt.Sprintf("%06d", rand.IntN(1000000))

	redisAdapter := s.Cache.Adapter()
	if rsErr := redisAdapter.Client().Set(redisAdapter.Context(), "otp:"+req.Phone, otp, 5*time.Minute).Err(); rsErr != nil {
		return SendOtpResponse{}, richerror.New(op).WithWrapper(rsErr)
	}

	// TODO: send SMS using gateway
	log.Printf("Phone: %s, OTP: %s\n", req.Phone, otp)

	return SendOtpResponse{Message: "OTP sent"}, nil
}

func (s Service) VerifyOtp(req VerifyOtpRequest) (VerifyOtpResponse, error) {
	const op = "authentication.service.VerifyOtp"

	redisAdapter := s.Cache.Adapter()

	attemptKey := "otp_attempts:" + req.Phone
	attempts, _ := redisAdapter.Client().Get(redisAdapter.Context(), attemptKey).Int()
	if attempts >= 5 { // Max 5 attempts
		return VerifyOtpResponse{}, richerror.New(op).WithKind(richerror.KindTooManyRequests).WithMessage("Too many failed attempts. Please request a new OTP")
	}

	// Validate OTP format (should be 6 digits)
	if len(req.Otp) != 6 {
		s.incrementFailedAttempts(req.Phone)
		return VerifyOtpResponse{}, richerror.New(op).WithKind(richerror.KindBadRequest).WithMessage("OTP must be 6 digits")
	}

	stored, err := redisAdapter.Client().Get(redisAdapter.Context(), "otp:"+req.Phone).Result()
	if errors.Is(err, redis.Nil) {
		s.incrementFailedAttempts(req.Phone)
		return VerifyOtpResponse{}, richerror.New(op).WithKind(richerror.KindGone).WithMessage("OTP has expired")
	} else if err != nil {
		return VerifyOtpResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithMessage(http.StatusText(http.StatusInternalServerError))
	}

	if stored != req.Otp {
		s.incrementFailedAttempts(req.Phone)
		return VerifyOtpResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessage("Invalid OTP code")
	}

	// Clear failed attempts on success
	_ = redisAdapter.Client().Del(redisAdapter.Context(), attemptKey).Err()

	// assign or accept deviceId
	deviceID := req.DeviceID
	if deviceID == "" {
		deviceID = uuid.NewString()
	}

	access, iaErr := s.issueAccess(req.Phone, deviceID)
	if iaErr != nil {
		return VerifyOtpResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithMessage("Failed to generate access token")
	}

	refresh, jti, irErr := s.issueRefresh(req.Phone, deviceID)
	if irErr != nil {
		return VerifyOtpResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithMessage("Failed to generate refresh token")
	}

	// Save latest refresh for this device
	if sErr := s.saveRefresh(req.Phone, deviceID, refresh); sErr != nil {
		return VerifyOtpResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithMessage("Failed to persist session")
	}

	// Clean up OTP
	_ = redisAdapter.Client().Del(redisAdapter.Context(), "otp:"+req.Phone).Err()

	// Optional: store meta
	_ = redisAdapter.Client().Set(redisAdapter.Context(), "refresh-meta:"+jti, fmt.Sprintf(`{"phone":"%s","deviceId":"%s"}`, req.Phone, deviceID), s.Config.AccessTokenTTL).Err()

	return VerifyOtpResponse{AccessToken: access, RefreshToken: refresh, DeviceID: deviceID}, nil

}

func (s Service) RefreshToken(req RefreshRequest) (RefreshResponse, error) {
	const op = "authentication.service.RefreshToken"
	redisAdapter := s.Cache.Adapter()

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return s.Config.AccessTokenSecret, nil
	})
	if err != nil || !token.Valid {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid refresh token")
	}

	phone, _ := claims["sub"].(string)
	deviceID, _ := claims["did"].(string)
	jti, _ := claims["jti"].(string)

	if phone == "" || deviceID == "" || jti == "" {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid refresh token")
	}

	if s.isBlacklisted(jti) {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid refresh token")
	}

	storedToken, err := s.getRefresh(phone, deviceID)
	if errors.Is(err, redis.Nil) {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid refresh token")
	} else if err != nil {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapper(err)
	}
	if storedToken != req.RefreshToken {
		s.blacklistJTI(jti)
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid refresh token")
	}

	firstUse, err := s.markUsedOnce(jti)
	if err != nil {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapper(err)
	}
	if !firstUse {
		s.deleteRefresh(phone, deviceID)
		s.blacklistJTI(jti)
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid refresh token")
	}

	access, err := s.issueAccess(phone, deviceID)
	if err != nil {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapper(err)
	}

	newRefresh, newJTI, err := s.issueRefresh(phone, deviceID)
	if err != nil {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapper(err)
	}

	if err := s.saveRefresh(phone, deviceID, newRefresh); err != nil {
		return RefreshResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapper(err)
	}

	// Optional: store new meta
	_ = redisAdapter.Client().Set(redisAdapter.Context(), "refresh-meta:"+newJTI, fmt.Sprintf(`{"phone":"%s","deviceId":"%s"}`, phone, deviceID), s.Config.RefreshTokenTTL).Err()

	return RefreshResponse{AccessToken: access, RefreshToken: newRefresh, DeviceID: deviceID}, nil
}

func (s Service) Me(req MeRequest) (MeResponse, error) {
	const op = "authentication.service.Me"

	claims, ok := req.Claims.(jwt.MapClaims)
	if !ok {
		return MeResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage(http.StatusText(http.StatusUnauthorized))
	}

	phone := claims["sub"].(string) // or int if you use that type

	// todo get user from db with phone and fill MeResponse with that user information
	return MeResponse{
		ID:       1,
		UserName: "hossein",
		Avatar:   "https://avatar.iran.liara.run/public/8",
		Phone:    phone,
	}, nil

}

func (s Service) Logout(req LogoutRequest) (LogoutResponse, error) {
	const op = "authentication.service.Logout"

	claims, ok := req.Claims.(jwt.MapClaims)
	if !ok {
		return LogoutResponse{}, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid authentication context")
	}

	phone, _ := claims["sub"].(string)
	deviceID, _ := claims["did"].(string)

	if phone != "" && deviceID != "" {
		s.deleteRefresh(phone, deviceID)
	}

	return LogoutResponse{
		Message: "Logged out successfully",
	}, nil

}

func (s Service) incrementFailedAttempts(phone string) {
	redisAdapter := s.Cache.Adapter()

	attemptKey := "otp_attempts:" + phone
	pipe := redisAdapter.Client().Pipeline()
	pipe.Incr(redisAdapter.Context(), attemptKey)
	pipe.Expire(redisAdapter.Context(), attemptKey, 15*time.Minute) // Reset after 15 minutes
	_, err := pipe.Exec(redisAdapter.Context())
	if err != nil {
		return
	}
}
