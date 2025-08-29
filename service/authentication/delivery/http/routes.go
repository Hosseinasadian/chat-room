package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"time"
)

func (h Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(10, time.Minute))

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
