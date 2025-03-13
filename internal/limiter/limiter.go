package limiter

import (
	"notification/internal/producer"
	"time"
)

// RateLimiter limit the rate of notifications
func RateLimiter(notificationInput <-chan producer.Notification, interval time.Duration) <-chan producer.Notification {
	notificationOutput := make(chan producer.Notification)

	go func() {
		defer close(notificationOutput)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for n := range notificationInput {
			<-ticker.C
			notificationOutput <- n
		}
	}()

	return notificationOutput
}
