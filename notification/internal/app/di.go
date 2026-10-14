package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/massodo1993/service-example/notification/internal/client/http"
	telegramClient "github.com/massodo1993/service-example/notification/internal/client/http/telegram"
	"github.com/massodo1993/service-example/notification/internal/config"
	kafkaConverter "github.com/massodo1993/service-example/notification/internal/converter/kafka"
	"github.com/massodo1993/service-example/notification/internal/converter/kafka/decoder"
	"github.com/massodo1993/service-example/notification/internal/service"
	assembledConsumer "github.com/massodo1993/service-example/notification/internal/service/consumer/order_assembled_consumer"
	paidConsumer "github.com/massodo1993/service-example/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/massodo1993/service-example/notification/internal/service/telegram"
	"github.com/massodo1993/service-example/platform/closer"
	"github.com/massodo1993/service-example/platform/logger"
	wrappedKafka "github.com/massodo1993/service-example/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/massodo1993/service-example/platform/pkg/kafka/consumer"
	kafkaMiddleware "github.com/massodo1993/service-example/platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramBot     *bot.Bot
	telegramClient  httpClient.TelegramClient
	telegramService service.TelegramService

	orderPaidConsumerService      service.ConsumerService
	orderAssembledConsumerService service.ConsumerService

	orderPaidConsumerGroup      sarama.ConsumerGroup
	orderAssembledConsumerGroup sarama.ConsumerGroup

	orderPaidConsumer      wrappedKafka.Consumer
	orderAssembledConsumer wrappedKafka.Consumer

	orderPaidDecoder     kafkaConverter.OrderPaidDecoder
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramBot(_ context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(d.TelegramClient(ctx))
	}

	return d.telegramService
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = paidConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = assembledConsumer.NewService(
			d.OrderAssembledConsumer(),
			d.ShipAssembledDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order paid consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka order paid consumer group", func(ctx context.Context) error {
			return group.Close()
		})

		d.orderPaidConsumerGroup = group
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order assembled consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka order assembled consumer group", func(ctx context.Context) error {
			return group.Close()
		})

		d.orderAssembledConsumerGroup = group
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{config.AppConfig().OrderPaidConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{config.AppConfig().OrderAssembledConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) ShipAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewShipAssembledDecoder()
	}

	return d.shipAssembledDecoder
}
