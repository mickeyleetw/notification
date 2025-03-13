package sender

import (
	"fmt"
	"notification/internal/producer"
	"notification/pkg/logger"
	"time"
)

// SendEmail send email notifications
func SendEmail(notificationInput <-chan producer.Notification) <-chan string {
	notificationOutput := make(chan string)

	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			logger.Info(fmt.Sprintf("Processing email notification ID: %d", n.ID))
			// mock email sending delay
			time.Sleep(1 * time.Second)
			msg := fmt.Sprintf("[Email Sent] ID: %d, Message: %s", n.ID, n.Message)
			logger.Info(fmt.Sprintf("Email sent successfully: %s", msg))
			notificationOutput <- msg
		}
	}()

	return notificationOutput
}
