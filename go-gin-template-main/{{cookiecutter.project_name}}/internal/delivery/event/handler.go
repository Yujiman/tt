package event

import (
	"context"

	"git-ffd.kz/fmobile/events/goevents"
	"git-ffd.kz/pkg/golog"
	watermilllog "git-ffd.kz/pkg/golog/contrib/watermill"
	"git-ffd.kz/pkg/gosentry"
	"git-ffd.kz/pkg/gosentry/contrib/gosentry_watermill"
	"git-ffd.kz/pkg/gotags/contrib/gotags_watermill"
	"git-ffd.kz/pkg/gowatermill/producer"
	"git-ffd.kz/pkg/gowatermill/subscriber"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	consumer "{{cookiecutter.package_name}}/internal/delivery/event/consumer"
	"{{cookiecutter.package_name}}/internal/service"
)

type Handler struct {
	subscriber subscriber.Subscriber
	services   *service.Services
	logger     golog.ContextLogger

	router    *message.Router
	producer  producer.Producer
	consumers *consumer.Consumer
}

func NewHandler(
	subscriber subscriber.Subscriber,
	services *service.Services,
	producer producer.Producer,
	logger golog.ContextLogger,
) (*Handler, error) {
	handler := &Handler{
		subscriber: subscriber,
		services:   services,
		logger:     logger,
		router:     nil,
		producer:   producer,
	}

	if err := handler.Init(); err != nil {
		return nil, err
	}

	return handler, nil
}

func (h *Handler) Init() (err error) {
	log := watermilllog.NewWatermillAdapter(h.logger)

	h.router, err = message.NewRouter(message.RouterConfig{}, log)
	if err != nil {
		return err
	}

	h.router.AddMiddleware(
		gotags_watermill.GoTags,
		middleware.Recoverer,
		gosentry_watermill.Sentry(func(msg *message.Message) *gosentry.Info {
			eventInfo, _ := goevents.GetEventInfo(msg.Payload)
			eventMap, _ := goevents.RawEventToMap(msg.Payload)

			if eventInfo.EventId == "" {
				eventInfo.EventId = msg.UUID
			}

			extras := make(map[string]interface{})
			if eventMap != nil {
				extras["event"] = eventMap
			}

			return &gosentry.Info{
				EventID: eventInfo.EventId,
				Extras:  extras,
			}
		}),
	)

	if err = h.initConsumers(h.router); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Run(ctx context.Context) error {
	return h.router.Run(ctx)
}

func (h *Handler) Close() error {
	if err := h.consumers.Close(); err != nil {
		return err
	}

	if h.router.IsRunning() {
		if err := h.router.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) initConsumers(router *message.Router) error {
	h.consumers = consumer.NewConsumer(
		h.services,
		h.producer,
		h.subscriber,
		h.logger,
	)

	return h.consumers.Init(router)
}
