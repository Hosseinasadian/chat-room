package httpresponse

import (
	"encoding/json"
	"log"
	"net/http"
)

func SetJsonContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func SetStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
func SetMessage(w http.ResponseWriter, message any) {
	body, mErr := json.Marshal(message)
	if mErr != nil {
		log.Printf("Error encoding JSON: %v", mErr)
		return
	}

	_, wErr := w.Write(body)
	if wErr != nil {
		return
	}
}
