package main

import (
	"context"
	"linkshortener/internal/config"
	postgresconnection "linkshortener/internal/database/postgres_connection"
	redisconnection "linkshortener/internal/database/redis_connection"
	"linkshortener/internal/logger"
	postgresrepo "linkshortener/internal/repository/postgres"
	redisrepo "linkshortener/internal/repository/redis"
	"linkshortener/internal/services"
	httptrans "linkshortener/internal/transport/http_server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	//logger
	//error
}

func init() {
	if err := godotenv.Load(""); err != nil {
		log.Print("Failed to load from env")
	}
}

// Парсинг флагов, соединение с бд, запуск сервера, логгер,
func main() {
	logger := logger.NewLogger("info")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.NewConfig()

	rdb, err := redisconnection.CreateRedisClient(ctx, cfg, logger)
	if err != nil {
		logger.Error("Failed to create redis client")
		os.Exit(1)
	}
	defer rdb.Close()

	pool, err := postgresconnection.PostgresCreatePool(ctx, cfg, logger)
	if err != nil {
		logger.Error("Failed to create pool")
		os.Exit(1)
	}
	defer pool.Close()

	linkRepo := postgresrepo.NewPostgresRepository(pool)
	cacheRepo := redisrepo.NewRedisRepository(rdb)

	linkService := services.NewLinkService(linkRepo, cacheRepo, logger)

	mux := http.NewServeMux()

	handler := httptrans.NewHandler(logger, linkService)
	handler.Route(mux)

	server := &http.Server{
		Addr:         ":" + cfg.App.AppPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server on port", "port", cfg.App.AppPort)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server closed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown fail", "error", err)
	}

	logger.Info("server stopped")

}
