package ratelimiter

import (
	"fmt"
	"github.com/hosseinasadian/chat-application/pkg/httpmsg"
	"github.com/hosseinasadian/chat-application/pkg/httpresponse"
	"net/http"
	"strings"
)

func RateLimitMiddleware(limiter RateLimiter, key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if err := limiter.Allow(r.Context(), fmt.Sprintf("%s:%s", key, clientIP)); err != nil {
				msg, code := httpmsg.Error(err)
				httpresponse.SetStatus(w, code)
				httpresponse.SetMessage(w, map[string]string{
					"error": msg,
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}
	parts := strings.Split(r.RemoteAddr, ":")
	return parts[0]
}
