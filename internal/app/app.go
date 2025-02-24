package app

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"news-rest-api/internal/config"
	"news-rest-api/internal/delivery"
	"news-rest-api/internal/middleware"
	"news-rest-api/internal/pkg/logger"
	"news-rest-api/internal/repository"
	"news-rest-api/internal/service"

	"github.com/gofiber/fiber/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

const (
	exitCodeError = 1
)

func Run(cfg *config.Config) {
	l := logger.InitLogger(slog.LevelDebug)

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		l.Error("error opening database", slog.String("err", err.Error()))
		os.Exit(exitCodeError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		l.Error("Database connection failed", slog.String("err", err.Error()))
		os.Exit(exitCodeError)
	}
	l.Info("Connected to database")

	reformDB := reform.NewDB(db, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))

	m := middleware.NewMiddleware(cfg.SecretKey, l)

	repo := repository.NewRepository(reformDB, l)
	service := service.NewService(repo)
	handler := delivery.NewHandler(service)

	app := fiber.New()
	app.Use(m.LoggerMiddleware())
	app.Use(m.AuthMiddleware())

	app.Get("/list", handler.GetNewsList)
	app.Post("/edit/:id", handler.EditNews)

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		l.Info("API server started successfully", "address", cfg.Server)
		if err = app.Listen(cfg.Server); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("Failed to start the server", slog.String("err", err.Error()))
			stop()
		}
	}()

	<-sigCtx.Done()
	l.Info("Received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err = app.ShutdownWithContext(shutdownCtx); err != nil {
		l.Error("Server shutdown error", slog.String("err", err.Error()))
	}
	shutdownCancel()

	if err = db.Close(); err != nil {
		l.Error("Error closing database", slog.String("err", err.Error()))
	}

	l.Info("Database connection closed")
	l.Info("Server stopped")
}
