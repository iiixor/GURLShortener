package memory

import (
	"URLShortener/internal/repository"
	"context"
	"sync"
)

// MemStorage is an in-memory implementation of the Storage interface.
// It uses a map to store URLs and a mutex to handle concurrent access.
type MemStorage struct {
	mu   sync.RWMutex
	data map[string]string
}

// New creates and returns a new MemStorage instance.
func New() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

// SaveURL saves a URL and its alias.
func (s *MemStorage) SaveURL(_ context.Context, urlToSave, alias string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[alias] = urlToSave
	return nil
}

// GetURL retrieves a URL by its alias.
func (s *MemStorage) GetURL(_ context.Context, alias string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.data[alias]
	if !ok {
		return "", repository.ErrAliasNotFound
	}

	return url, nil
}

// AliasExists checks if an alias already exists.
func (s *MemStorage) AliasExists(_ context.Context, alias string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[alias]
	return ok, nil
}
