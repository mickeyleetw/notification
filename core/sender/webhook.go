package sender

import (
	"fmt"
	"notification/core/history"
	"notification/core/producer"
	"notification/pkg/logger"
	"time"
)

// SendWebhook send webhook notifications
func SendWebhook(notificationInput <-chan producer.Notification, errs chan<- error, historyCh chan<- history.Entry) <-chan string {
	notificationOutput := make(chan string)

	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			logger.Info(fmt.Sprintf("Processing webhook notification ID: %d", n.ID))

			if n.Message == "" {
				errs <- fmt.Errorf("[Sender Error] Webhook message is empty. ID: %d", n.ID)
				continue
			}

			// simulate webhook sending
			time.Sleep(2 * time.Second)
			history.StoreNotification(n, "webhook")
			msg := fmt.Sprintf("[Webhook Sent] ID: %d, Message: %s\n", n.ID, n.Message)
			logger.Info(fmt.Sprintf("Webhook sent successfully: %s", msg))

			// Add to history
			historyCh <- history.Entry{
				Stage:        "Webhook Sender",
				Notification: n,
				Message:      "Webhook sent successfully",
			}
			notificationOutput <- msg
		}
	}()

	return notificationOutput
}
