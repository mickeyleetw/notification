package processor

import (
	"notification/internal/producer"
	"strings"
)

// Processor is a function that processes the notification
func Processor(notificationInput <-chan producer.Notification) <-chan producer.Notification {
	notificationOutput := make(chan producer.Notification)

	go func() {
		defer close(notificationOutput)
		for n := range notificationInput {
			// mock processing
			n.Message = strings.TrimSpace(n.Message)
			n.Message = strings.ToUpper(n.Message)

			notificationOutput <- n
		}
	}()

	return notificationOutput
}
