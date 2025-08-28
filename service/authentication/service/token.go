package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hosseinasadian/chat-application/pkg/richerror"
	"strings"
	"time"
)

func (s Service) issueAccess(phone, deviceID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   phone,
		"did":   deviceID,
		"scope": "access",
		"exp":   time.Now().Add(s.config.AccessTokenTTL).Unix(),
		"iat":   time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.config.AccessTokenSecret))
}

func (s Service) issueRefresh(phone, deviceID string) (tokenString string, jti string, err error) {
	jti = uuid.NewString()
	claims := jwt.MapClaims{
		"sub": phone,
		"did": deviceID,
		"jti": jti,
		"scp": "refresh",
		"exp": time.Now().Add(s.config.RefreshTokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = t.SignedString([]byte(s.config.RefreshTokenSecret))
	return
}

func (s Service) saveRefresh(phone, deviceID, refresh string) error {
	key := fmt.Sprintf("refresh:%s:%s", phone, deviceID)
	redisAdapter := s.otpRepo.Adapter()
	return redisAdapter.Client().Set(redisAdapter.Context(), key, refresh, s.config.RefreshTokenTTL).Err()
}

func (s Service) getRefresh(phone, deviceID string) (string, error) {
	key := fmt.Sprintf("refresh:%s:%s", phone, deviceID)
	redisAdapter := s.otpRepo.Adapter()
	return redisAdapter.Client().Get(redisAdapter.Context(), key).Result()
}

func (s Service) deleteRefresh(phone, deviceID string) {
	key := fmt.Sprintf("refresh:%s:%s", phone, deviceID)
	redisAdapter := s.otpRepo.Adapter()
	_ = redisAdapter.Client().Del(redisAdapter.Context(), key).Err()
}

func (s Service) markUsedOnce(jti string) (bool, error) {
	key := "used:" + jti
	redisAdapter := s.otpRepo.Adapter()
	ok, err := redisAdapter.Client().SetNX(redisAdapter.Context(), key, 1, s.config.RefreshTokenTTL).Result()
	return ok, err
}

func (s Service) blacklistJTI(jti string) {
	redisAdapter := s.otpRepo.Adapter()
	_ = redisAdapter.Client().Set(redisAdapter.Context(), "revoked:"+jti, 1, s.config.RefreshTokenTTL).Err()
}

func (s Service) isBlacklisted(jti string) bool {
	redisAdapter := s.otpRepo.Adapter()
	_, err := redisAdapter.Client().Get(redisAdapter.Context(), "revoked:"+jti).Result()
	return err == nil
}

func (s Service) ParseToken(bearerToken string) (jwt.MapClaims, error) {
	const op = "authentication/service.ParseToken"

	if bearerToken == "" {
		return nil, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Missing token")
	}

	tokenString := strings.TrimPrefix(bearerToken, "Bearer ")
	if tokenString == bearerToken {
		return nil, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid token format")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.AccessTokenSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, richerror.New(op).WithKind(richerror.KindUnauthorized).WithMessage("Invalid claims")
}
