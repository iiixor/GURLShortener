package service

import (
	"URLShortener/internal/config"
	"URLShortener/internal/repository"
	"URLShortener/pkg/random"
	"context"
	"fmt"
)

const maxRetries = 5 // Number of times to try generating a unique alias.

// URLShortenerService provides the business logic for shortening URLs.
type URLShortenerService struct {
	storage repository.Storage
	config  *config.URLShortener
}

// NewURLShortener creates a new URLShortenerService.
func NewURLShortener(storage repository.Storage, cfg *config.URLShortener) *URLShortenerService {
	return &URLShortenerService{
		storage: storage,
		config:  cfg,
	}
}

// Shorten generates a unique short alias for the given URL and saves it.
func (s *URLShortenerService) Shorten(ctx context.Context, urlToShorten string) (string, error) {
	var alias string
	var err error

	for i := 0; i < maxRetries; i++ {
		alias, err = random.NewRandomString(s.config.AliasLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate alias: %w", err)
		}

		exists, err := s.storage.AliasExists(ctx, alias)
		if err != nil {
			return "", fmt.Errorf("failed to check alias existence: %w", err)
		}

		if !exists {
			break
		}
	}

	if err := s.storage.SaveURL(ctx, urlToShorten, alias); err != nil {
		return "", fmt.Errorf("failed to save url: %w", err)
	}

	return alias, nil
}
