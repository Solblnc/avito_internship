package server

import (
	_ "avito_internship/docs"
	"context"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Service interface {
	Create(segment string, percent uint) (int, error)
	Delete(segment string) error
	AddUser(segmentsAdd []string, segmentsDelete []string, userId int) error
	AddSegmentDeadline(ttl int, segmentName string, userId int) error
	GetActualSegments(userId int) []string
	CreateUser() error
	GetHistory(year, month int) (string, error)
}

type Server struct {
	Router  *mux.Router
	Service Service
	Server  *http.Server
}

func NewServer(service Service) *Server {
	h := &Server{
		Service: service,
	}

	h.Router = mux.NewRouter()
	h.mapRoutes()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeoutMiddleware)

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}
	return h
}

func (h *Server) mapRoutes() {
	h.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/segment/create_segment", h.CreateSegment).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/segment/delete_segment", h.DeleteSegment).Methods(http.MethodDelete)
	h.Router.HandleFunc("/api/v1/user/get_segments", h.GetSegments).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/user/create_user", h.CreateUsers).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/segment/get_history", h.GetHistory).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/segment/add_user_segment", h.AddUserToSegment).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/segment/add_user_deadline", h.AddUserDeadline).Methods(http.MethodPost)
}

func (h *Server) Serve() error {
	log.Println("Server is running on port :8080")
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("shut down gracefully")
	return nil
}
