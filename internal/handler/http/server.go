package http

import (
	"URLShortener/internal/repository"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// RedirectHandler handles the redirection from a short alias to the original URL.
func RedirectHandler(log *zap.Logger, storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the alias from the URL path, removing the leading "/"
		alias := strings.TrimPrefix(r.URL.Path, "/")
		if alias == "" {
			http.NotFound(w, r)
			return
		}

		log.Info("received redirect request", zap.String("alias", alias))

		// Retrieve the original URL from storage
		originalURL, err := storage.GetURL(r.Context(), alias)
		if err != nil {
			// If the alias is not found, it's a 404
			if errors.Is(err, repository.ErrAliasNotFound) {
				log.Warn("alias not found", zap.String("alias", alias))
				http.NotFound(w, r)
				return
			}
			// For any other error, it's an internal server error
			log.Error("failed to get URL from storage", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Perform the redirect
		http.Redirect(w, r, originalURL, http.StatusFound) // 302 Found
	}
}

// NewServer creates a new http.Server with the redirect handler.
func NewServer(addr string, log *zap.Logger, storage repository.Storage) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", RedirectHandler(log, storage))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
