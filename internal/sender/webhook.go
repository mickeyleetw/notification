package sender

import (
	"fmt"
	"notification/internal/producer"
	"notification/pkg/logger"
	"time"
)

// SendWebhook send webhook notifications
func SendWebhook(notificationInput <-chan producer.Notification) <-chan string {
	notificationOutput := make(chan string)

	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			logger.Info(fmt.Sprintf("Processing webhook notification ID: %d", n.ID))
			// simulate webhook sending
			time.Sleep(2 * time.Second)
			msg := fmt.Sprintf("[Webhook Sent] ID: %d, Message: %s\n", n.ID, n.Message)
			logger.Info(fmt.Sprintf("Webhook sent successfully: %s", msg))
			notificationOutput <- msg
		}
	}()

	return notificationOutput
}
