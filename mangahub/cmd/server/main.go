package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"mangahub/internal/config"
	"mangahub/internal/db"
	grpcservice "mangahub/internal/grpc"
	"mangahub/internal/realtime"
	"mangahub/internal/router"
	"mangahub/internal/tcp"
	"mangahub/internal/udp"
)

func main() {
	// CLI flags
	migrate := flag.Bool("migrate", false, "run database migrations and exit")
	seed    := flag.Bool("seed", false, "seed database with sample data and exit")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	log.Printf("[main] starting MangaHub (env=%s)", cfg.AppEnv)

	// Connect to PostgreSQL
	sqlDB, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("database connect failed: %v", err)
	}
	defer db.Close()

	// Run migrations
	migrationsPath := filepath.Join("data", "migrations")
	if err := db.Migrate(sqlDB, migrationsPath); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	if *migrate {
		log.Println("[main] migrations done, exiting")
		os.Exit(0)
	}

	// Seed data
	if *seed {
		if err := seedData(sqlDB); err != nil {
			log.Fatalf("seed failed: %v", err)
		}
		log.Println("[main] seed done, exiting")
		os.Exit(0)
	}

	// Start WebSocket hub
	hub := realtime.NewHub()

	// Start TCP sync server
	tcpServer := tcp.NewServer(cfg.TCPPort)
	if err := tcpServer.Start(); err != nil {
		log.Printf("[main] TCP server start warning: %v", err)
	}

	// Start UDP notification server
	udpServer := udp.NewServer(cfg.UDPPort)
	if err := udpServer.Start(); err != nil {
		log.Printf("[main] UDP server start warning: %v", err)
	}

	// Start gRPC server
	if err := grpcservice.StartGRPCServer(cfg.GRPCPort, sqlDB); err != nil {
		log.Printf("[main] gRPC server start warning: %v", err)
	}

	// Build HTTP router
	r := router.New(sqlDB, hub)

	// HTTP server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[main] HTTP server listening on http://localhost%s", addr)
		log.Printf("[main] TCP  sync server  on port %s", cfg.TCPPort)
		log.Printf("[main] UDP  notify server on port %s", cfg.UDPPort)
		log.Printf("[main] gRPC service       on port %s", cfg.GRPCPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	log.Println("[main] shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("[main] goodbye")
}
