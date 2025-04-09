package sender

import (
	"fmt"
	"notification/core/history"
	"notification/core/producer"
	"notification/pkg/logger"
	"time"
)

// SendPush send push notifications
func SendPush(notificationInput <-chan producer.Notification, errs chan<- error, historyCh chan<- history.Entry) <-chan string {
	notificationOutput := make(chan string)

	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			logger.Info(fmt.Sprintf("Processing push notification ID: %d", n.ID))

			if n.Message == "" {
				errs <- fmt.Errorf("[Sender Error] Push message is empty. ID: %d", n.ID)
				continue
			}

			// simulate push sending
			time.Sleep(2 * time.Second)
			history.StoreNotification(n, "sender")
			msg := fmt.Sprintf("[Push Sent] ID: %d, Message: %s\n", n.ID, n.Message)
			logger.Info(fmt.Sprintf("Push notification sent successfully: %s", msg))

			// Add to history
			historyCh <- history.Entry{
				Stage:        "Push Sender",
				Notification: n,
				Message:      "Push sent successfully",
			}

			notificationOutput <- msg
		}
	}()

	return notificationOutput
}
