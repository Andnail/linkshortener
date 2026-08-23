package main

import (
	"context"
	"linkshortener/internal/config"
	postgresconnection "linkshortener/internal/database/postgres_connection"
	"linkshortener/internal/logger"
	"linkshortener/internal/repository/postgres"
	"linkshortener/internal/services"
	httptrans "linkshortener/internal/transport/http_trans"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	//_ "github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	//logger
	//error
}

// Парсинг флагов, соединение с бд, запуск сервера, логгер,
func main() {
	logg := logger.NewLogger("info")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := "postgres://web:database64@localhost:5432/links"

	cfg := config.NewConfig(dsn, 2, 10, time.Hour)

	pool, err := postgresconnection.PostgresCreatePool(ctx, cfg, logg)
	if err != nil {
		logg.Error("Failed to create pool")
	}

	defer pool.Close()

	linkRepo := postgres.NewPostgresRepository(pool)

	linkService := services.NewLinkService(linkRepo, logg)

	mux := http.NewServeMux()

	handler := httptrans.NewHandler(logg, linkService)
	handler.Route(mux)

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logg.Info("Starting server on port", "port", cfg.AppPort)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Error("Server closed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logg.Error("server shutdown fail", "error", err)
	}

	logg.Info("server stopped")

}
