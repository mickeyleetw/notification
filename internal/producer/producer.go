package producer

import "errors"

// Notification is a struct that contains the notification data
type Notification struct {
	ID      int
	Type    string
	Message string
}

// Producer generate notifications and send them through a channel, and return an error channel
func Producer() (<-chan Notification, <-chan error) {
	out := make(chan Notification)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		notifications := []Notification{
			{ID: 1, Type: "Email", Message: "Welcome email"},
			{ID: 2, Type: "SMS", Message: "Your OTP code"},
			{ID: 3, Type: "Push", Message: "New push notification"},
			{ID: 4, Type: "Webhook", Message: "Webhook triggered"},
		}

		if len(notifications) == 0 {
			errs <- errors.New("no notifications to generate")
			return
		}

		for _, n := range notifications {
			out <- n
		}
	}()

	return out, errs
}
