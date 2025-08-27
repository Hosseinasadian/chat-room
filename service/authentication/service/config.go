package service

import "time"

type Config struct {
	AccessTokenSecret  string        `koanf:"access_token_secret"`
	AccessTokenTTL     time.Duration `koanf:"access_token_ttl"`
	RefreshTokenSecret string        `koanf:"refresh_token_secret"`
	RefreshTokenTTL    time.Duration `koanf:"refresh_token_ttl"`
}
