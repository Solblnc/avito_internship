package server

import (
	"avito_internship/internal/database"
	"log"
	"net/http"
)

type Server struct {
	Port string `mapstructure:"PORT"`
	DB   *database.DataBase
}

func NewServer(port string, db *database.DataBase) (*Server, error) {
	return &Server{Port: port, DB: db}, nil
}

func (s *Server) setupGetHandlers() {
	http.HandleFunc("/api/v1/get_segments", s.GetSegments)
}

func (s *Server) setupPostHandlers() {
	http.HandleFunc("/api/v1/create_segment", loggerMiddleware(s.CreateSegment))
	http.HandleFunc("/api/v1/delete_segment", loggerMiddleware(s.DeleteSegment))
	http.HandleFunc("/api/v1/add_user_segment", loggerMiddleware(s.AddUserToSegment))
}

func (s *Server) SetupHandlers() {
	s.setupGetHandlers()
	s.setupPostHandlers()
}

func (s *Server) Run() {
	log.Println("Server is running on port :8080")
	log.Fatal(http.ListenAndServe(s.Port, nil))
}
