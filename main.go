package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/karandeepbhardwaj/pixl.ink/internal/config"
	"github.com/karandeepbhardwaj/pixl.ink/internal/handler"
	"github.com/karandeepbhardwaj/pixl.ink/internal/server"
	"github.com/karandeepbhardwaj/pixl.ink/internal/storage"
)

func main() {
	cfg := config.Load()

	db, err := storage.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	disk, err := storage.NewDiskStore(cfg.UploadDir)
	if err != nil {
		log.Fatalf("Failed to initialize disk storage: %v", err)
	}

	h := handler.New(cfg, db, disk)
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      server.New(cfg, h),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
