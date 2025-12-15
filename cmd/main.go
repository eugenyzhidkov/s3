package main

import (
	"context"
	"os/signal"
	"s3/internal/app"
	"syscall"

	"s3/internal/infrastructure/logger"
)

func main() {
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM) // Можно использовать os.Interrupt, а не SIGINT?

	/*os.Interrupt — кроссплатформенная константа. На всех ОС (Linux, macOS, Windows) она означает "прерывание от пользователя" (Ctrl+C).
	  syscall.SIGINT — низкоуровневая, и на Windows может вести себя иначе.
	  Официальная документация Go и все примеры используют os.Interrupt.
	*/
	defer stop()

	if err := app.AppStart(ctx); err != nil {
		logger.Fatal("Сервис завершился с ошибкой", "error", err)
	}
}
