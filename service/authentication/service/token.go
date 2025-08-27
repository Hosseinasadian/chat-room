package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func (s Service) issueAccess(phone, deviceID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   phone,
		"did":   deviceID,
		"scope": "access",
		"exp":   time.Now().Add(s.Config.AccessTokenTTL).Unix(),
		"iat":   time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.Config.AccessTokenSecret))
}

func (s Service) issueRefresh(phone, deviceID string) (tokenString string, jti string, err error) {
	jti = uuid.NewString()
	claims := jwt.MapClaims{
		"sub": phone,
		"did": deviceID,
		"jti": jti,
		"scp": "refresh",
		"exp": time.Now().Add(s.Config.RefreshTokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = t.SignedString([]byte(s.Config.RefreshTokenSecret))
	return
}

func (s Service) saveRefresh(phone, deviceID, refresh string) error {
	key := fmt.Sprintf("refresh:%s:%s", phone, deviceID)
	redisAdapter := s.Cache.Adapter()
	return redisAdapter.Client().Set(redisAdapter.Context(), key, refresh, s.Config.RefreshTokenTTL).Err()
}

func (s Service) getRefresh(phone, deviceID string) (string, error) {
	key := fmt.Sprintf("refresh:%s:%s", phone, deviceID)
	redisAdapter := s.Cache.Adapter()
	return redisAdapter.Client().Get(redisAdapter.Context(), key).Result()
}

func (s Service) deleteRefresh(phone, deviceID string) {
	key := fmt.Sprintf("refresh:%s:%s", phone, deviceID)
	redisAdapter := s.Cache.Adapter()
	_ = redisAdapter.Client().Del(redisAdapter.Context(), key).Err()
}

func (s Service) markUsedOnce(jti string) (bool, error) {
	key := "used:" + jti
	redisAdapter := s.Cache.Adapter()
	ok, err := redisAdapter.Client().SetNX(redisAdapter.Context(), key, 1, s.Config.RefreshTokenTTL).Result()
	return ok, err
}

func (s Service) blacklistJTI(jti string) {
	redisAdapter := s.Cache.Adapter()
	_ = redisAdapter.Client().Set(redisAdapter.Context(), "revoked:"+jti, 1, s.Config.RefreshTokenTTL).Err()
}

func (s Service) isBlacklisted(jti string) bool {
	redisAdapter := s.Cache.Adapter()
	_, err := redisAdapter.Client().Get(redisAdapter.Context(), "revoked:"+jti).Result()
	return err == nil
}
