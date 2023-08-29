package server

import (
	"avito_internship/internal/database"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Port   string `mapstructure:"PORT"`
	DB     *database.DataBase
	Server *http.Server
}

func NewServer(port string, db *database.DataBase) (*Server, error) {
	return &Server{Port: port, DB: db, Server: &http.Server{Addr: port, Handler: nil}}, nil
}

func (s *Server) setupGetHandlers() {
	http.HandleFunc("/api/v1/get_segments", loggerMiddleware(s.GetSegments))
	http.HandleFunc("/api/v1/create_user", loggerMiddleware(s.CreateUsers))
	http.HandleFunc("/api/v1/get_history", loggerMiddleware(s.GetHistory))
}

func (s *Server) setupPostHandlers() {
	http.HandleFunc("/api/v1/create_segment", loggerMiddleware(s.CreateSegment))
	http.HandleFunc("/api/v1/delete_segment", loggerMiddleware(s.DeleteSegment))
	http.HandleFunc("/api/v1/add_user_segment", loggerMiddleware(s.AddUserToSegment))
	http.HandleFunc("/api/v1/add_user_deadline", loggerMiddleware(s.AddUserDeadline))
}

func (s *Server) SetupHandlers() {
	s.setupGetHandlers()
	s.setupPostHandlers()
}

func (s *Server) Run() error {

	go func() {
		log.Println("Server is running on port :8080")
		if err := s.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.Server.Shutdown(ctx)

	log.Println("shut down gracefully")
	return nil
}
