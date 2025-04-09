package producer

import (
	"errors"
	"fmt"
)

// Notification is a struct that contains the notification data
type Notification struct {
	ID      int
	Type    string
	Message string
}

// Producer generate notifications and send them through a channel, and return an error channel
func Producer(errs chan<- error) <-chan Notification {
	out := make(chan Notification)

	go func() {
		defer close(out)

		notifications := []Notification{
			{ID: 1, Type: "Email", Message: "Welcome email"},
			{ID: 2, Type: "SMS", Message: "Your OTP code"},
			{ID: 3, Type: "Push", Message: "Push hello"},
			{ID: 4, Type: "Webhook", Message: "Webhook triggered"},
			{ID: 5, Type: "Unknown", Message: "Oops!"}, // unknown type
		}

		if len(notifications) == 0 {
			errs <- fmt.Errorf("Producer error: %w", errors.New("no notifications to generate"))
			return
		}

		for _, n := range notifications {
			out <- n
		}
	}()

	return out
}
