package main

import (
	"context"
	"fmt"
	"notification/core/dispatcher"
	"notification/core/history"
	"notification/core/limiter"
	"notification/core/processor"
	"notification/core/producer"
	"notification/core/sender"
	"notification/pkg/logger"
	"sync"
	"time"
)

func main() {
	logger.Info("🚀 Start notification system...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Step 0: Create shared channels
	errsChannel := make(chan error)
	historyCh := make(chan history.Entry)
	historyStore := history.NewStore()

	// Step 0.1: Start history saving goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer logger.Info("History goroutine stopped")
		for entry := range historyCh {
			historyStore.Add(entry)
		}
	}()

	// Step 1: Producer
	notificationsChannel := producer.Producer(errsChannel)

	// Step 2: Error handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer logger.Info("Error handler goroutine stopped")
		printErrorResults("Error", errsChannel)
	}()

	// Step 3~7: Pipeline
	processedCh := processor.Processor(notificationsChannel, errsChannel, historyCh)
	limitedCh := limiter.RateLimiter(processedCh, 500*time.Millisecond, errsChannel, historyCh)
	dispatchMap := dispatcher.Route(limitedCh, errsChannel, historyCh)

	emailResultCh := sender.SendEmail(dispatchMap.EmailCh, errsChannel, historyCh)
	smsResultCh := sender.SendSMS(dispatchMap.SMSCh, errsChannel, historyCh)
	webhookResultCh := sender.SendWebhook(dispatchMap.WebhookCh, errsChannel, historyCh)
	pushResultCh := sender.SendPush(dispatchMap.PushCh, errsChannel, historyCh)

	// Step 8: Print results
	wg.Add(4)
	go func() {
		defer wg.Done()
		printResults("Email", emailResultCh)
	}()
	go func() {
		defer wg.Done()
		printResults("SMS", smsResultCh)
	}()
	go func() {
		defer wg.Done()
		printResults("Webhook", webhookResultCh)
	}()
	go func() {
		defer wg.Done()
		printResults("Push", pushResultCh)
	}()

	// Step 9: Wait
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Info("Context cancelled, shutting down...")
	case <-done:
		logger.Info("All notifications processed, shutting down...")
	case <-time.After(10 * time.Second):
		logger.Info("Timeout reached, shutting down...")
		cancel()
	}

	// Step 10: Close shared channels
	close(errsChannel)
	close(historyCh)

	// Export history to a file
	err := history.ExportToFile("history.json")
	if err != nil {
		logger.Error(fmt.Sprintf("Error exporting history to file: %v", err))
	} else {
		logger.Info("History successfully exported to history.json")
	}

	wg.Wait()
	logger.Info("Notification system stopped.")
}

// Utility function to print result or error messages
func printResults(tag string, ch <-chan string) {
	for msg := range ch {
		fmt.Printf("[%s Result] %s\n", tag, msg)
	}
}

func printErrorResults(tag string, ch <-chan error) {
	for err := range ch {
		fmt.Printf("[%s Error] %s\n", tag, err)
	}
}
