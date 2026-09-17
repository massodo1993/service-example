//go:build integration

package integration

import (
	"context"

	"go.uber.org/zap"

	"github.com/massodo1993/service-example/platform/logger"
)

// teardownTestEnvironment — освобождает все ресурсы тестового окружения
func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	logger.Info(ctx, "🧹 Очистка тестового окружения...")

	cleanupTestEnvironment(ctx, env)

	logger.Info(ctx, "✅ Тестовое окружение успешно очищено")
}

// cleanupTestEnvironment — вспомогательная функция для освобождения ресурсов
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер приложения", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер приложения остановлен")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер MongoDB", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер MongoDB остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "не удалось удалить сеть", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Сеть удалена")
		}
	}
}
