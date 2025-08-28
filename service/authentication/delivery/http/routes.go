package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/hosseinasadian/chat-application/pkg/ratelimiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func (h Handler) Routes(client *redis.Client) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		privateRateLimiter := ratelimiter.New(client, 10, time.Minute)
		privateRateLimitMiddleware := ratelimiter.RateLimitMiddleware(*privateRateLimiter, "authenticatePrivate")
		r.Use(privateRateLimitMiddleware)

		r.Post("/send-otp", h.SendOtpHandler)
		r.Post("/verify-otp", h.VerifyOtpHandler)
		r.Post("/refresh-token", h.RefreshTokenHandler)
	})

	r.Route("/me", func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		r.Get("/", h.MeHandler)
		r.Post("/logout", h.LogoutHandler)
	})

	return r
}
