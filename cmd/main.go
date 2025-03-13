package main

import (
	"context"
	"fmt"
	"notification/internal/dispatcher"
	"notification/internal/limiter"
	"notification/internal/processor"
	"notification/internal/producer"
	"notification/internal/sender"
	"notification/pkg/logger"
	"sync"
	"time"
)

func main() {
	logger.Info("🚀 Start notification system...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Step 1: Create notifications and start a notification producer channel
	notificationsChannel, errsChannel := producer.Producer()

	// Step 2: Handle errors using different goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer logger.Info("Error handler goroutine stopped")
		for err := range errsChannel {
			if err != nil {
				fmt.Println("[ERROR] Producer:", err)
			}
		}
	}()

	// Step 3: Processor
	processedCh := processor.Processor(notificationsChannel)

	// Step 4: Rate Limiter
	limitedCh := limiter.RateLimiter(processedCh, 500*time.Millisecond) // 500ms one notification

	// Step 5: Dispatcher
	dispatchMap := dispatcher.Route(limitedCh)

	// Step 6: Sender (returns result channels)
	emailResultCh := sender.SendEmail(dispatchMap.EmailCh)
	smsResultCh := sender.SendSMS(dispatchMap.SMSCh)
	webhookResultCh := sender.SendWebhook(dispatchMap.WebhookCh)
	pushResultCh := sender.SendPush(dispatchMap.PushCh)

	// Step 7: Combine and print results
	wg.Add(4)
	go func() {
		defer wg.Done()
		defer logger.Info("Email result printer stopped")
		printResults("Email", emailResultCh)
	}()
	go func() {
		defer wg.Done()
		defer logger.Info("SMS result printer stopped")
		printResults("SMS", smsResultCh)
	}()
	go func() {
		defer wg.Done()
		defer logger.Info("Webhook result printer stopped")
		printResults("Webhook", webhookResultCh)
	}()
	go func() {
		defer wg.Done()
		defer logger.Info("Push result printer stopped")
		printResults("Push", pushResultCh)
	}()

	// Create a channel to signal when all notifications are processed
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for all goroutines to complete or context cancellation
	select {
	case <-ctx.Done():
		logger.Info("Context cancelled, shutting down...")
	case <-done:
		logger.Info("All notifications processed, shutting down...")
	case <-time.After(10 * time.Second): // 增加等待時間到 10 秒
		logger.Info("Timeout reached, shutting down...")
		cancel()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	logger.Info("Notification system stopped.")
}

// Utility function to print result messages
func printResults(tag string, ch <-chan string) {
	for msg := range ch {
		fmt.Printf("[%s Result] %s\n", tag, msg)
	}
}
