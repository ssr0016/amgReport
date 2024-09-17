package main

import (
	"amg/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"amg/internal/server"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	cfg.Logger.Info("Starting server", zap.String("port", cfg.Port))

	s := server.NewServer(cfg)

	go func() {
		err := s.Start()
		if err != nil {
			log.Fatalf("Server failed to start: %v", err)
			cfg.Logger.Error("Server failed to start", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Stop()
	if err != nil {
		log.Fatalf("Server failed to shutdown: %v", err)
	}

	log.Println("Server shutdown gracefully!")
}
