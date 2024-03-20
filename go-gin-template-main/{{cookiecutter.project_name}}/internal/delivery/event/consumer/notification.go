// Пример консьюмера. Удалите, если в проекте не будут использоваться обработка событий.

package event

// import (
// 	"time"
//
// 	"git-ffd.kz/fmobile/events/goevents"
// 	"git-ffd.kz/pkg/gobackoff"
// 	"git-ffd.kz/pkg/goerr"
// 	"github.com/ThreeDotsLabs/watermill/message"
// )
//
// func (c *Consumer) initNotificationConsumer(router *message.Router) error {
// 	if err := c.subscriber.AddNoPublishHandler(
// 		router,
// 		c.SendSmsHandler,
// 		goevents.Notification_SendSms_SubjectName,
// 		"sendSms",
// 		gobackoff.ExponentialBackOff{
// 			InitialInterval: time.Millisecond * 300,
// 			MaxInterval:     time.Second * 20,
// 			MaxElapsedTime:  time.Minute * 2,
// 		},
// 	); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (c *Consumer) SendSmsHandler(msg *message.Message) error {
// 	ctx, _ := c.logger.FromContext(msg.Context(), "handler", "SendSmsHandler", "event_id", msg.UUID)
//
// 	eventName, err := goevents.GetEventName(msg.Payload)
// 	if err != nil {
// 		return goerr.Wrap(err).WithCtx(msg.Context())
// 	}
//
// 	switch eventName {
// 	case goevents.Notification_SendSms_v1_EventName:
// 		event, err := goevents.NewSendSmsV1_FromRaw(msg.Payload)
// 		if err != nil {
// 			return goerr.Wrap(err).WithCtx(ctx)
// 		}
//
// 		err = c.services.Notification.SendSms(ctx, event.Data.Phone, event.Data.Text)
// 		if err != nil {
// 			return goerr.Wrap(err).WithCtx(ctx)
// 		}
//
// 		return nil
// 	default:
// 		return nil
// 	}
// }
