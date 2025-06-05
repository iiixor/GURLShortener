package main

import (
	"URLShortener/internal/bot"
	"URLShortener/internal/config"
	httpHandler "URLShortener/internal/handler/http"
	"URLShortener/internal/repository/memory"
	"URLShortener/internal/service"
	"URLShortener/pkg/logger"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// --- Initialization ---
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)
	defer log.Sync()

	log.Info("starting application", zap.String("env", cfg.Env))

	// --- Dependency Injection ---
	storage := memory.New()
	urlShortenerService := service.NewURLShortener(storage, &cfg.URLShortener)

	// --- Running Services ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run Telegram Bot
	telegramBot, err := bot.New(cfg, log, urlShortenerService)
	if err != nil {
		log.Fatal("failed to initialize bot", zap.Error(err))
	}
	telegramBot.Start(ctx) // Start runs in a goroutine now

	// Run HTTP Server for redirects
	httpServer := httpHandler.NewServer(cfg.HTTPServer.Address, log, storage)
	go func() {
		log.Info("starting http server", zap.String("address", cfg.Address))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed to start", zap.Error(err))
			cancel() // Stop other services if http server fails
		}
	}()

	// --- Graceful Shutdown ---
	// Wait for interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Info("received shutdown signal", zap.String("signal", sig.String()))
	case <-ctx.Done():
		log.Info("context cancelled, shutting down")
	}

	// Shutdown server
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Error("http server shutdown error", zap.Error(err))
	}

	log.Info("application shut down gracefully")
}
