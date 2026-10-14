package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/massodo1993/service-example/notification/internal/client/http"
	"github.com/massodo1993/service-example/notification/internal/model"
	def "github.com/massodo1993/service-example/notification/internal/service"
	"github.com/massodo1993/service-example/platform/logger"
)

// chatID захардкожен по условию задания. Свой ID можно узнать из лога после /start.
const chatID int64 = 149162978

//go:embed templates/*.tmpl
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.tmpl"))

var _ def.TelegramService = (*service)(nil)

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{telegramClient: telegramClient}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error {
	return s.send(ctx, "paid_notification.tmpl", event)
}

func (s *service) SendShipAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error {
	return s.send(ctx, "assembled_notification.tmpl", event)
}

func (s *service) send(ctx context.Context, tmpl string, data any) error {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, tmpl, data); err != nil {
		return err
	}

	if err := s.telegramClient.SendMessage(ctx, chatID, buf.String()); err != nil {
		return err
	}

	logger.Info(ctx, "сообщение отправлено в Telegram", zap.Int64("chat_id", chatID), zap.String("template", tmpl))

	return nil
}
