package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func (h *Handler) RedisMiddlewareCache(next http.Handler, ctx context.Context, db *redis.Client) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req CreateShortLinkRequest

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Incorrect JSON format", http.StatusBadRequest)
				return
			}

			if req.OriginLink == "" {
				http.Error(w, "Link is empty", http.StatusBadRequest)
			}

			if shortLink, err := db.Get(ctx, req.OriginLink).Result(); err == nil && shortLink != "" {
				w.Header().Set("Content-Type", "application-json")
				w.WriteHeader(http.StatusCreated)

				if err := json.NewEncoder(w).Encode(shortLink); err != nil {
					http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
				}
				return

			}
			next.ServeHTTP(w, r)
		})
}
