package dispatcher

import (
	"notification/internal/producer"
)

// OutputChannels is a struct that contains the output channels
type OutputChannels struct {
	EmailCh   chan producer.Notification
	SMSCh     chan producer.Notification
	WebhookCh chan producer.Notification
	PushCh    chan producer.Notification
}

// Route is a function that routes the notification to the appropriate channel
func Route(notificationInput <-chan producer.Notification) OutputChannels {
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
			}
		}
	}()

	return OutputChannels{
		EmailCh:   emailCh,
		SMSCh:     smsCh,
		WebhookCh: webhookCh,
		PushCh:    pushCh,
	}
}
