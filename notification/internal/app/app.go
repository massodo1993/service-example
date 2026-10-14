package app

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/massodo1993/service-example/notification/internal/config"
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

	errCh := make(chan error, 2)

	go func() {
		if err := a.diContainer.OrderPaidConsumerService(ctx).RunConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("консьюмер order paid завершился с ошибкой: %w", err)
		}
	}()

	go func() {
		if err := a.diContainer.OrderAssembledConsumerService(ctx).RunConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("консьюмер order assembled завершился с ошибкой: %w", err)
		}
	}()

	go func() {
		logger.Info(ctx, "🤖 Telegram-бот запущен")
		a.diContainer.TelegramBot(ctx).Start(ctx)
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
		a.initTelegramBot,
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

func (a *App) initTelegramBot(ctx context.Context) error {
	a.diContainer.TelegramBot(ctx).RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			logger.Info(ctx, "chat id", zap.Int64("chat_id", update.Message.Chat.ID))

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "🚀 Привет! Бот активирован — вы будете получать уведомления об оплате и сборке заказов.",
			})
			if err != nil {
				logger.Error(ctx, "не удалось отправить приветственное сообщение", zap.Error(err))
			}
		},
	)

	return nil
}
