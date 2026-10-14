package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/massodo1993/service-example/assembly/internal/config"
	kafkaConverter "github.com/massodo1993/service-example/assembly/internal/converter/kafka"
	"github.com/massodo1993/service-example/assembly/internal/converter/kafka/decoder"
	"github.com/massodo1993/service-example/assembly/internal/service"
	orderConsumer "github.com/massodo1993/service-example/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/massodo1993/service-example/assembly/internal/service/producer/order_producer"
	"github.com/massodo1993/service-example/platform/closer"
	"github.com/massodo1993/service-example/platform/logger"
	wrappedKafka "github.com/massodo1993/service-example/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/massodo1993/service-example/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/massodo1993/service-example/platform/pkg/kafka/producer"
	kafkaMiddleware "github.com/massodo1993/service-example/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderProducerService service.OrderProducerService
	consumerService      service.ConsumerService

	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer wrappedKafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder

	syncProducer           sarama.SyncProducer
	orderAssembledProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderAssembledProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) ConsumerService() service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = orderConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.OrderProducerService(),
		)
	}

	return d.consumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{config.AppConfig().OrderPaidConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderAssembledProducer() wrappedKafka.Producer {
	if d.orderAssembledProducer == nil {
		d.orderAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderAssembledProducer
}
