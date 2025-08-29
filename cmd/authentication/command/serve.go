package command

import (
	"context"
	"fmt"
	"github.com/go-chi/httprate"
	redisAdapter "github.com/hosseinasadian/chat-application/adapter/redis"
	"github.com/hosseinasadian/chat-application/pkg/configloader"
	"github.com/hosseinasadian/chat-application/pkg/httpserver"
	authHttp "github.com/hosseinasadian/chat-application/service/authentication/delivery/http"
	authRepository "github.com/hosseinasadian/chat-application/service/authentication/repository"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/hosseinasadian/chat-application/service/authentication"
	authService "github.com/hosseinasadian/chat-application/service/authentication/service"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the authentication service",
	Long:  `This command starts the authentication service.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	var cfg *authentication.Config
	//var cf authService.Config
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v", err)
	}

	yamlPath := os.Getenv("CONFIG_PATH")
	if yamlPath == "" {
		yamlPath = filepath.Join(workingDir, "deploy", "authentication", "development", "config.yaml")
	}

	options := configloader.Option{
		Prefix:       "AUTHENTICATION_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: yamlPath,
		CallbackEnv:  nil,
	}

	if err := configloader.Load(options, &cfg); err != nil {
		log.Fatalf("Failed to load food config: %v", err)
	}

	rdAdapter, rdErr := redisAdapter.New(context.Background(), cfg.Redis)

	if rdErr != nil {
		log.Fatal(rdErr)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	authCache := authRepository.New(*rdAdapter)
	authSvc := authService.New(cfg.AuthService, authCache)

	loginRateLimiter := httprate.NewRateLimiter(5, 15*time.Minute)
	authHandler := authHttp.New(authSvc, loginRateLimiter)

	server := httpserver.New(cfg.HTTPServer, authHandler)

	svc := authentication.Setup(logger, *cfg, server)
	svc.Start()

}

func init() {
	RootCommand.AddCommand(serveCmd)
}
