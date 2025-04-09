package sender

import (
	"fmt"
	"notification/core/history"
	"notification/core/producer"
	"notification/pkg/logger"
	"time"
)

// SendEmail send email notifications
func SendEmail(notificationInput <-chan producer.Notification, errs chan<- error, historyCh chan<- history.Entry) <-chan string {
	notificationOutput := make(chan string)

	go func() {
		defer close(notificationOutput)

		for n := range notificationInput {
			logger.Info(fmt.Sprintf("Processing email notification ID: %d", n.ID))

			if n.Message == "" {
				errs <- fmt.Errorf("[Sender Error] Email message is empty. ID: %d", n.ID)
				continue
			}

			// mock email sending delay
			time.Sleep(1 * time.Second)
			history.StoreNotification(n, "email")
			msg := fmt.Sprintf("[Email Sent] ID: %d, Message: %s", n.ID, n.Message)
			logger.Info(fmt.Sprintf("Email sent successfully: %s", msg))

			// Add to history
			historyCh <- history.Entry{
				Stage:        "Email Sender",
				Notification: n,
				Message:      "Email sent successfully",
			}

			notificationOutput <- msg
		}
	}()

	return notificationOutput
}
