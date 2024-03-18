package server

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	http.Server
}

func NewServer(myHandler http.Handler) *Server {
	return &Server{
		http.Server{
			Addr:              ":8080",
			Handler:           myHandler,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	log.Println("Start server on :8080")
	return s.ListenAndServe()
}
