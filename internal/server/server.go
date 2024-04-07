package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type APIServer struct {
	httpServer *http.Server
}

func NewAPIServer(port string, handler http.Handler) *APIServer {
	return &APIServer{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20, //1 MB
			ReadTimeout:    time.Second * 10,
			WriteTimeout:   time.Second * 10,
		},
	}
}

func (s *APIServer) Run() error {
	logrus.Printf("Listening on %s \n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
