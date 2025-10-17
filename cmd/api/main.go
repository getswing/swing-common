package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"getswing.app/player-service/internal/config"
	"getswing.app/player-service/internal/db"
	"getswing.app/player-service/internal/models"
	"getswing.app/player-service/internal/repository"
	"getswing.app/player-service/internal/shared"

	service "getswing.app/player-service/internal/services"
	pb "getswing.app/player-service/proto"

	"google.golang.org/grpc"
)

// Global context (shared across the application)
var (
	AppCtx    context.Context
	AppCancel context.CancelFunc
)

func init() {
	AppCtx, AppCancel = context.WithCancel(context.Background())
}

func main() {
	defer AppCancel()

	cfg := config.Load()

	log.Printf("[api] starting with port=%s log_level=%s", cfg.HTTPPort, cfg.LogLevel)

	// DB connect
	gormDB, sqlDB, err := db.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("[api] db connect error: %v", err)
	}
	defer sqlDB.Close()

	// Auto migrate
	gormDB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err := gormDB.AutoMigrate(&models.Player{}); err != nil {
		log.Fatalf("[api] db migrate error: %v", err)
	}

	// Setup Echo
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Use(shared.WithConfig(&cfg))

	// Run HTTP server concurrently
	serverErr := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%s", cfg.HTTPPort)
		serverErr <- e.Start(addr)
	}()

	// Start gRPC client/server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	playerRepository := repository.NewPlayerRepository(gormDB)

	pb.RegisterPlayerServiceServer(grpcServer, service.NewPlayerService(playerRepository))

	log.Println("✅ gRPC server running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Graceful shutdown
	select {
	case err := <-serverErr:
		if err != nil {
			log.Printf("[api] server error: %v", err)
		}
	case <-AppCtx.Done():
		log.Println("[api] context canceled, shutting down gracefully...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("[api] shutdown error: %v", err)
	}
}
