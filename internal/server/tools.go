package server

import (
	"log"
	"net/http"
)

func jsonRespond(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s \n", r.Method, r.URL)
		next(w, r)
	}
}
