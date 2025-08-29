package authentication

import (
	"github.com/hosseinasadian/chat-application/adapter/redis"
	"github.com/hosseinasadian/chat-application/pkg/httpserver"
	authService "github.com/hosseinasadian/chat-application/service/authentication/service"
	"time"
)

type Config struct {
	TotalShutdownTimeout time.Duration      `koanf:"total_shutdown_timeout"`
	HTTPServer           httpserver.Config  `koanf:"http_server"`
	AuthService          authService.Config `koanf:"auth_service"`
	Redis                redis.Config       `koanf:"redis"`
}
