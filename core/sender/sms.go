package sender

import (
	"fmt"
	"notification/core/history"
	"notification/core/producer"
	"notification/pkg/logger"
	"time"
)

// SendSMS send sms notifications
func SendSMS(notificationInput <-chan producer.Notification, errs chan<- error, historyCh chan<- history.Entry) <-chan string {
	notificationOutput := make(chan string)

	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			logger.Info(fmt.Sprintf("Processing SMS notification ID: %d", n.ID))

			if n.Message == "" {
				errs <- fmt.Errorf("[Sender Error] SMS message is empty. ID: %d", n.ID)
				continue
			}

			// simulate sms sending
			time.Sleep(500 * time.Millisecond)
			history.StoreNotification(n, "sms")
			msg := fmt.Sprintf("[SMS Sent] ID: %d, Message: %s\n", n.ID, n.Message)
			logger.Info(fmt.Sprintf("SMS sent successfully: %s", msg))

			// Add to history
			historyCh <- history.Entry{
				Stage:        "SMS Sender",
				Notification: n,
				Message:      "SMS sent successfully",
			}
			notificationOutput <- msg
		}
	}()

	return notificationOutput
}
