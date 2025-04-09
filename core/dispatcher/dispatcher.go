package dispatcher

import (
	"fmt"
	"notification/core/history"
	"notification/core/producer"
	"notification/pkg/logger"
)

// DispatchMap is a struct that contains the output channels
type DispatchMap struct {
	EmailCh   chan producer.Notification
	SMSCh     chan producer.Notification
	WebhookCh chan producer.Notification
	PushCh    chan producer.Notification
}

// Route is a function that routes the notification to the appropriate channel
func Route(notificationInput <-chan producer.Notification, errs chan<- error, historyCh chan<- history.Entry) DispatchMap {
	emailCh := make(chan producer.Notification)
	smsCh := make(chan producer.Notification)
	webhookCh := make(chan producer.Notification)
	pushCh := make(chan producer.Notification)

	go func() {
		defer close(emailCh)
		defer close(smsCh)
		defer close(webhookCh)
		defer close(pushCh)

		for n := range notificationInput {
			switch n.Type {
			case "Email":
				emailCh <- n
			case "SMS":
				smsCh <- n
			case "Webhook":
				webhookCh <- n
			case "Push":
				pushCh <- n
			default:
				errs <- fmt.Errorf("[Dispatcher Error] Unsupported notification type: %s, ID: %d", n.Type, n.ID)
				logger.Error(fmt.Sprintf("[Dispatcher Error] Unsupported notification type: %s, ID: %d", n.Type, n.ID))
			}
		}

	}()

	return DispatchMap{
		EmailCh:   emailCh,
		SMSCh:     smsCh,
		WebhookCh: webhookCh,
		PushCh:    pushCh,
	}
}
