package event

import (
	"git-ffd.kz/pkg/golog"
	"git-ffd.kz/pkg/gowatermill/producer"
	"git-ffd.kz/pkg/gowatermill/subscriber"
	"github.com/ThreeDotsLabs/watermill/message"

	"{{cookiecutter.package_name}}/internal/service"
)

type Consumer struct {
	services   *service.Services
	producer   producer.Producer
	subscriber subscriber.Subscriber
	logger     golog.ContextLogger
}

func NewConsumer(
	services *service.Services,
	producer producer.Producer,
	subscriber subscriber.Subscriber,
	logger golog.ContextLogger,
) *Consumer {
	return &Consumer{
		services:   services,
		producer:   producer,
		subscriber: subscriber,
		logger:     logger,
	}
}

func (c *Consumer) Init(router *message.Router) error {
	// REPLACE initNotificationConsumer with your consumer!
	// if err := c.initNotificationConsumer(router); err != nil {
	//	return fmt.Errorf("failed initialize notification consumer: %w", err)
	// }

	return nil
}

func (c *Consumer) Close() error {
	return c.subscriber.Close()
}
