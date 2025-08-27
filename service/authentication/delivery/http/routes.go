package http

import (
	"github.com/go-chi/chi/v5"
)

func (h Handler) Routes() chi.Router {
	r := chi.NewRouter()
	
	r.Post("/send-otp", h.SendOtpHandler)
	r.Post("/verify-otp", h.VerifyOtpHandler)
	r.Post("/refresh-token", h.RefreshTokenHandler)

	r.Route("/me", func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		r.Get("/", h.MeHandler)
		r.Post("/logout", h.LogoutHandler)
	})

	return r
}
