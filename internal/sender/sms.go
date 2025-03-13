package sender

import (
	"fmt"
	"notification/internal/producer"
	"notification/pkg/logger"
	"time"
)

// SendSMS send sms notifications
func SendSMS(notificationInput <-chan producer.Notification) <-chan string {
	notificationOutput := make(chan string)

	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			logger.Info(fmt.Sprintf("Processing SMS notification ID: %d", n.ID))
			// simulate sms sending
			time.Sleep(500 * time.Millisecond)
			msg := fmt.Sprintf("[SMS Sent] ID: %d, Message: %s\n", n.ID, n.Message)
			logger.Info(fmt.Sprintf("SMS sent successfully: %s", msg))
			notificationOutput <- msg
		}
	}()

	return notificationOutput
}
