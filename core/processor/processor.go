package processor

import (
	"fmt"
	"notification/core/history"
	"notification/core/producer"
	"strings"
)

// Processor is a function that processes the notification
func Processor(notificationInput <-chan producer.Notification, errs chan<- error, historyCh chan<- history.Entry) <-chan producer.Notification {
	notificationOutput := make(chan producer.Notification)
	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			if n.Message == "" {
				errs <- fmt.Errorf("empty message for notification ID: %d", n.ID)
				continue
			}
			// mock processing
			n.Message = strings.TrimSpace(n.Message)
			n.Message = strings.ToUpper(n.Message)
			n.Message = "[Processed] " + n.Message
			// Add to history
			historyCh <- history.Entry{
				Stage:        "Processor",
				Notification: n,
				Message:      "Processed successfully",
			}

			notificationOutput <- n
		}
	}()

	return notificationOutput
}
