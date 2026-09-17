package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/massodo1993/service-example/inventory/internal/app"
	"github.com/massodo1993/service-example/inventory/internal/config"
	"github.com/massodo1993/service-example/platform/closer"
	"github.com/massodo1993/service-example/platform/logger"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("не удалось загрузить конфиг: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	app, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "не удалось создать приложение", zap.Error(err))
		return
	}

	err = app.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "ошибка в работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "ошибка при завершении работы", zap.Error(err))
	}
}
