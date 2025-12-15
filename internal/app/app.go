package app

import (
	"context"
	"fmt"
	"net/http"

	"s3/internal/config"
	"s3/internal/handler"
	"s3/internal/infrastructure/logger"
	"s3/internal/s3"
)

func AppStart(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Ошибка загрузки конфигурации: %w", err)
	}

	s3Client, err := s3.New(cfg)
	if err != nil {
		return fmt.Errorf("Ошибка создания S3-клиента: %w", err)
	}

	if err := s3Client.Ping(ctx); err != nil {
		return fmt.Errorf("Не удалось подключиться к S3: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.IndexHandler(s3Client))
	mux.HandleFunc("/upload", handler.UploadHandler(s3Client))
	mux.HandleFunc("/move", handler.MoveHandler(s3Client))

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	go func() {
		logger.Info("HTTP-сервер запущен", "port", ":"+cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Сервер завершился с ошибкой", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Получен сигнал завершения — начинаем graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Ошибка graceful shutdown", "error", err)
		return fmt.Errorf("Ошибка graceful shutdown HTTP-сервера: %w", err)
	}

	logger.Info("Сервер корректно остановлен")
	return nil
}
