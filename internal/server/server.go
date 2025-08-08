package server

import (
	"CryptoPriceCollection/internal/types"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Server struct {
	config *types.ConfigApp
	server *http.Server
	router *gin.Engine
}

func NewServer(config *types.ConfigApp, router *gin.Engine) *Server {
	return &Server{
		config: config,
		router: router,
	}
}

// Start запуск сервера
func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.HTTPServer.Port),
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.HTTPServer.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.HTTPServer.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.HTTPServer.IdleTimeout) * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server starting on :%d", s.config.HTTPServer.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Server received shutdown signal")
		return s.Shutdown()
	case err := <-serverErr:
		return err
	}
}

// Shutdown красивое затыкание сервера, ахах
func (s *Server) Shutdown() error {
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.HTTPServer.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server shutdown complete")
	return nil
}
