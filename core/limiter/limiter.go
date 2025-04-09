package limiter

import (
	"fmt"
	"notification/core/history"
	"notification/core/producer"
	"time"
)

// RateLimiter limit the rate of notifications
func RateLimiter(notificationInput <-chan producer.Notification, interval time.Duration, errs chan<- error, historyCh chan<- history.Entry) <-chan producer.Notification {
	notificationOutput := make(chan producer.Notification)

	go func() {
		defer close(notificationOutput)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for n := range notificationInput {
			select {
			case <-ticker.C:
				notificationOutput <- n
			default:
				errs <- fmt.Errorf("rate limited notification ID: %d", n.ID)
				time.Sleep(interval)
				// Add to history
				historyCh <- history.Entry{
					Stage:        "Rate Limiter",
					Notification: n,
					Message:      fmt.Sprintf("Notification ID: %d rate limited", n.ID),
				}

				notificationOutput <- n
			}
		}
	}()

	return notificationOutput
}
