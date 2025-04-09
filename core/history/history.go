package history

import (
	"encoding/json"
	"fmt"
	"notification/core/producer"
	"os"
	"sync"
)

// Entry is a struct that contains the history entry
type Entry struct {
	Stage        string
	Notification producer.Notification
	Message      string // optional: detailed info or result message
}

// Store is a struct that contains the history store
type Store struct {
	mu    sync.Mutex
	items []Entry
}

var globalStore = NewStore()

// NewStore is a function that creates a new history store
func NewStore() *Store {
	return &Store{
		items: make([]Entry, 0),
	}
}

// Add is a function that adds a history entry to the store
func (s *Store) Add(entry Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, entry)
}

// All is a function that returns all history entries
func (s *Store) All() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Entry(nil), s.items...) // return a copy
}

// StoreNotification allows shared access to add a notification to the global store
func StoreNotification(n producer.Notification, stage string) {
	globalStore.Add(Entry{
		Stage:        stage,
		Notification: n,
	})
}

// GetAllHistory returns all recorded history entries
func GetAllHistory() []Entry {
	return globalStore.All()
}

// ExportToFile is a function that exports the history to a file
func ExportToFile(filePath string) error {
	globalStore.mu.Lock()
	defer globalStore.mu.Unlock()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(globalStore.items); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
