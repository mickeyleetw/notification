package test

import (
	"notification/core/history"
	"notification/core/processor"
	"notification/core/producer"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessor_ValidMessage(t *testing.T) {
	// Setup
	input := make(chan producer.Notification, 1)
	errs := make(chan error, 1)
	historyCh := make(chan history.Entry, 1)

	// Test data
	testNotification := producer.Notification{
		ID:      1,
		Type:    "email",
		Message: "hello world",
	}

	// Send test data
	input <- testNotification
	close(input)

	// Run processor
	output := processor.Processor(input, errs, historyCh)

	// Verify results
	select {
	case result := <-output:
		// Check message processing
		assert.Equal(t, "[Processed] HELLO WORLD", result.Message, "Message should be processed correctly")
		assert.Equal(t, testNotification.ID, result.ID, "ID should remain unchanged")
		assert.Equal(t, testNotification.Type, result.Type, "Type should remain unchanged")

		// Check history entry
		historyEntry := <-historyCh
		assert.Equal(t, "Processor", historyEntry.Stage)
		assert.Equal(t, "Processed successfully", historyEntry.Message)
		assert.Equal(t, result, historyEntry.Notification)

	case <-time.After(1 * time.Second):
		assert.Fail(t, "timeout waiting for processor output")
	}

	// Verify no errors
	select {
	case err := <-errs:
		assert.Fail(t, "unexpected error", err)
	default:
		// No error expected
	}
}

func TestProcessor_EmptyMessage(t *testing.T) {
	// Setup
	input := make(chan producer.Notification, 1)
	errs := make(chan error, 1)
	historyCh := make(chan history.Entry, 1)

	// Test data
	testNotification := producer.Notification{
		ID:      1,
		Type:    "email",
		Message: "",
	}

	// Send test data
	input <- testNotification
	close(input)

	// Run processor
	output := processor.Processor(input, errs, historyCh)

	// Verify error
	select {
	case err := <-errs:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty message") // check if error message contains "empty message"
		assert.Contains(t, err.Error(), "1")             // Check ID in error message
	case <-time.After(1 * time.Second):
		assert.Fail(t, "timeout waiting for error")
	}

	// Verify output (if any)
	select {
	case result := <-output:
		assert.Empty(t, result.Message, "Output message should be empty")
		assert.Equal(t, 0, result.ID, "ID should remain unchanged")
		assert.Equal(t, "", result.Type, "Type should remain unchanged")
	default:
		// No output is also acceptable
	}

	// Verify no history entry
	select {
	case entry := <-historyCh:
		assert.Fail(t, "unexpected history entry", entry)
	default:
		// No history entry expected
	}
}

func TestProcessor_MultipleMessages(t *testing.T) {
	// Setup
	input := make(chan producer.Notification, 3)
	errs := make(chan error, 3)
	historyCh := make(chan history.Entry, 3)

	// Test data
	messages := []producer.Notification{
		{ID: 1, Type: "email", Message: "first message"},
		{ID: 2, Type: "sms", Message: "second message"},
		{ID: 3, Type: "push", Message: "third message"},
	}

	// Send test data
	for _, msg := range messages {
		input <- msg
	}
	close(input)

	// Run processor
	output := processor.Processor(input, errs, historyCh)

	// Verify results
	for i := 0; i < len(messages); i++ {
		select {
		case result := <-output:
			expectedMsg := "[Processed] " + strings.ToUpper(messages[i].Message)
			assert.Equal(t, expectedMsg, result.Message, "Message should be processed correctly")
			assert.Equal(t, messages[i].ID, result.ID, "ID should remain unchanged")
			assert.Equal(t, messages[i].Type, result.Type, "Type should remain unchanged")

			// Verify history entry
			historyEntry := <-historyCh
			assert.Equal(t, "Processor", historyEntry.Stage)
			assert.Equal(t, "Processed successfully", historyEntry.Message)
			assert.Equal(t, result, historyEntry.Notification)

		case <-time.After(1 * time.Second):
			assert.Fail(t, "timeout waiting for processor output")
		}
	}

	// Verify no errors
	select {
	case err := <-errs:
		assert.Fail(t, "unexpected error", err)
	default:
		// No error expected
	}
}
