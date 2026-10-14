package app

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/massodo1993/service-example/assembly/internal/config"
	"github.com/massodo1993/service-example/platform/closer"
	"github.com/massodo1993/service-example/platform/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("консьюмер завершился с ошибкой: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "получен сигнал завершения")
		return nil
	case err := <-errCh:
		logger.Error(ctx, "компонент завершился с ошибкой, останавливаем приложение", zap.Error(err))
		return err
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "Kafka-консьюмер OrderPaid запущен")

	return a.diContainer.ConsumerService().RunConsumer(ctx)
}
