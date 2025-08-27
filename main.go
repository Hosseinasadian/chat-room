package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	redisAdapter "github.com/hosseinasadian/chat-application/adapter/redis"
	authHttp "github.com/hosseinasadian/chat-application/service/authentication/delivery/http"
	authRepository "github.com/hosseinasadian/chat-application/service/authentication/repository"
	authService "github.com/hosseinasadian/chat-application/service/authentication/service"
)

var (
	accessSecret  = "super-secret-access-key"
	refreshSecret = "super-secret-refresh-key"
	accessTTL     = 15 * time.Minute
	refreshTTL    = 7 * 24 * time.Hour
	ctx           = context.Background()
)

func main() {
	rdAdapter, rdErr := redisAdapter.New(ctx, redisAdapter.Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	})

	if rdErr != nil {
		log.Fatal(rdErr)
	}

	authCache := authRepository.New(*rdAdapter)
	authConfig := authService.Config{
		AccessTokenSecret:  accessSecret,
		AccessTokenTTL:     accessTTL,
		RefreshTokenSecret: refreshSecret,
		RefreshTokenTTL:    refreshTTL,
	}
	authSvc := authService.New(authConfig, authCache)
	authHandler := authHttp.New(authSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	log.Printf("Hi, Server started on port 8080")
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // React dev server
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			return
		}
	})

	r.Mount("/auth", authHandler.Routes())

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
