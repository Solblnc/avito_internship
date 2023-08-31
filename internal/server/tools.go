package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

func jsonRespond(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application-json; charset=UTF-8")

		next.ServeHTTP(writer, request)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logrus.WithFields(
			logrus.Fields{
				"method": request.Method,
				"path":   request.URL.Path,
			}).Info("handled request")

		next.ServeHTTP(writer, request)
	})
}

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithTimeout(request.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
