package httpserver

import (
	"log"
	"net"
	"net/http"
)

type ServerConfig struct {
	Host string
	Port string
}

type Server struct {
	server *http.Server
	notify chan error
}

func NewServer(handler http.Handler, cfg *ServerConfig) *Server {
	httpServer := &http.Server{
		Handler: handler,
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
	}
	return &Server{
		server: httpServer,
		notify: make(chan error, 1),
	}
}

func (s *Server) Start() {
	go func() {
		log.Printf("HTTP server listening on %s", s.server.Addr)
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.server.Close()
}
