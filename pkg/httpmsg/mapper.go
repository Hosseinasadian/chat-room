package httpmsg

import (
	"errors"
	"github.com/hosseinasadian/chat-application/pkg/richerror"
	"net/http"
)

func Error(err error) (message string, code int) {
	var richError richerror.RichError
	switch {
	case errors.As(err, &richError):
		var re richerror.RichError
		errors.As(err, &re)
		message := re.Message()
		code := MapKindToHTTPStatusCode(re.Kind())

		return message, code
	default:
		return err.Error(), http.StatusBadRequest
	}
}

func MapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindBadRequest:
		return http.StatusBadRequest
	case richerror.KindTooManyRequests:
		return http.StatusTooManyRequests
	case richerror.KindGone:
		return http.StatusGone
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindUnauthorized:
		return http.StatusUnauthorized
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
