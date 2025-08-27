package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hosseinasadian/chat-application/pkg/httpmsg"
	"github.com/hosseinasadian/chat-application/pkg/httpresponse"
	"github.com/hosseinasadian/chat-application/service/authentication/service"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	AuthSvc service.Service
}

func New(authSvc service.Service) Handler {
	return Handler{
		AuthSvc: authSvc,
	}
}

func (h Handler) SendOtpHandler(w http.ResponseWriter, r *http.Request) {
	var req service.SendOtpRequest
	if dErr := json.NewDecoder(r.Body).Decode(&req); dErr != nil {
		msg, code := httpmsg.Error(dErr)
		httpresponse.SetJsonContentType(w)
		httpresponse.SetStatus(w, code)
		httpresponse.SetMessage(w, map[string]string{
			"error": msg,
		})
		return
	}

	res, sErr := h.AuthSvc.SendOtp(req)
	if sErr != nil {
		msg, code := httpmsg.Error(sErr)
		httpresponse.SetJsonContentType(w)
		httpresponse.SetStatus(w, code)
		httpresponse.SetMessage(w, map[string]string{
			"error": msg,
		})
		return
	}

	httpresponse.SetJsonContentType(w)
	httpresponse.SetMessage(w, res)
}

func (h Handler) VerifyOtpHandler(w http.ResponseWriter, r *http.Request) {
	var req service.VerifyOtpRequest
	if fErr := json.NewDecoder(r.Body).Decode(&req); fErr != nil {
		msg, code := httpmsg.Error(fErr)
		httpresponse.SetJsonContentType(w)
		httpresponse.SetStatus(w, code)
		httpresponse.SetMessage(w, map[string]string{
			"error": msg,
		})
		return
	}

	res, vErr := h.AuthSvc.VerifyOtp(req)
	if vErr != nil {
		msg, code := httpmsg.Error(vErr)
		httpresponse.SetJsonContentType(w)
		httpresponse.SetStatus(w, code)
		httpresponse.SetMessage(w, map[string]string{
			"error": msg,
		})
		return
	}

	httpresponse.SetJsonContentType(w)
	httpresponse.SetMessage(w, res)
}

func (h Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req service.RefreshRequest
	if dErr := json.NewDecoder(r.Body).Decode(&req); dErr != nil || req.RefreshToken == "" {
		msg, code := httpmsg.Error(dErr)
		httpresponse.SetStatus(w, code)
		httpresponse.SetMessage(w, map[string]string{
			"error": msg,
		})
		return
	}

	res, rErr := h.AuthSvc.RefreshToken(req)
	if rErr != nil {
		msg, code := httpmsg.Error(rErr)
		httpresponse.SetStatus(w, code)
		httpresponse.SetMessage(w, map[string]string{
			"error": msg,
		})
		return
	}

	httpresponse.SetJsonContentType(w)
	httpresponse.SetMessage(w, res)
}

func (h Handler) MeHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.AuthSvc.Me(service.MeRequest{
		Claims: r.Context().Value("claims"),
	})

	if err != nil {
		msg, code := httpmsg.Error(err)
		httpresponse.SetMessage(w, map[string]string{
			"message": msg,
		})
		httpresponse.SetStatus(w, code)
		return
	}

	httpresponse.SetMessage(w, res)
	httpresponse.SetStatus(w, http.StatusOK)

}

func (h Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.AuthSvc.Logout(service.LogoutRequest{
		Claims: r.Context().Value("claims"),
	})

	if err != nil {
		msg, code := httpmsg.Error(err)
		httpresponse.SetStatus(w, code)
		httpresponse.SetMessage(w, map[string]string{
			"error": msg,
		})
		return
	}

	httpresponse.SetJsonContentType(w)
	httpresponse.SetMessage(w, res)
}

func (h Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpresponse.SetStatus(w, http.StatusUnauthorized)
			httpresponse.SetMessage(w, map[string]string{
				"error": "Missing token",
			})
			return
		}

		// Expect format: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			httpresponse.SetStatus(w, http.StatusUnauthorized)
			httpresponse.SetMessage(w, map[string]string{
				"error": "Invalid token format",
			})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.AuthSvc.Config.AccessTokenSecret), nil
		})
		if err != nil || !token.Valid {
			httpresponse.SetStatus(w, http.StatusUnauthorized)
			httpresponse.SetMessage(w, map[string]string{
				"error": "Invalid token",
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		httpresponse.SetStatus(w, http.StatusUnauthorized)
		httpresponse.SetMessage(w, map[string]string{
			"error": "Invalid claims",
		})
	})
}
