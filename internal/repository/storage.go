package repository

import (
	"context"
	"errors"
)

// ErrAliasNotFound is returned when an alias is not found in the storage.
var ErrAliasNotFound = errors.New("alias not found")

// Storage defines the interface for URL storage.
// This allows us to swap the implementation (e.g., from in-memory to a real database)
// without changing the business logic that uses it.
type Storage interface {
	// SaveURL saves the original URL and its alias to the storage.
	SaveURL(ctx context.Context, urlToSave, alias string) error

	// GetURL retrieves the original URL by its alias.
	// It should return ErrAliasNotFound if the alias does not exist.
	GetURL(ctx context.Context, alias string) (string, error)

	// AliasExists checks if an alias already exists in the storage.
	AliasExists(ctx context.Context, alias string) (bool, error)
}
