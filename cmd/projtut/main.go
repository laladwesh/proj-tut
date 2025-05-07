package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/laladwesh/proj-tut/internal/config"
	"github.com/laladwesh/proj-tut/internal/http/handlers/student"
	"github.com/laladwesh/proj-tut/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String("path", cfg.StoragePath), slog.String("env", cfg.Env))
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.Getlist(storage))

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Starting server...", slog.String("address", cfg.HTTPServer.Addr))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()

	<-done
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", "error", err.Error())
	} else {
		slog.Info("Server shutdown gracefully")
	}
}
