package sender

import (
	"fmt"
	"notification/internal/producer"
	"notification/pkg/logger"
	"time"
)

// SendPush send push notifications
func SendPush(notificationInput <-chan producer.Notification) <-chan string {
	notificationOutput := make(chan string)

	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			logger.Info(fmt.Sprintf("Processing push notification ID: %d", n.ID))
			// simulate push sending
			time.Sleep(2 * time.Second)
			msg := fmt.Sprintf("[Push Sent] ID: %d, Message: %s\n", n.ID, n.Message)
			logger.Info(fmt.Sprintf("Push notification sent successfully: %s", msg))
			notificationOutput <- msg
		}
	}()

	return notificationOutput
}
