package app

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"token-service/internal/pkg/logger"
	"token-service/internal/pkg/store"
)

type tokenResponse struct {
	Token     string `json:"token"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

func StartHTTPServer(addr string, apiKey string, kv store.KeyValueStore, lgr *logger.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/token", apiKeyAuth(apiKey, tokenHandler(kv, lgr), lgr))

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lgr.PrintError("HTTP server error", err)
		}
	}()

	lgr.PrintSuccess("HTTP server listening on " + addr)
	return server
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func tokenHandler(kv store.KeyValueStore, lgr *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token, err := kv.GetKey("defenseToken")
		if err != nil {
			if errors.Is(err, store.ErrKeyNotFound) {
				http.Error(w, "token not available", http.StatusNotFound)
				return
			}
			lgr.PrintError("Error on get defense token from store", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		updatedAt, _ := kv.GetKey("defenseTokenUpdatedAt")

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tokenResponse{
			Token:     token,
			UpdatedAt: updatedAt,
		}); err != nil {
			lgr.PrintError("Error on write token response", err)
		}
	}
}

func apiKeyAuth(apiKey string, next http.Handler, lgr *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimSpace(r.Header.Get("X-API-Key"))
		if key == "" {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				key = strings.TrimSpace(authHeader[7:])
			}
		}

		if !matchAPIKey(key, apiKey) {
			lgr.PrintWarning("Unauthorized token request")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func matchAPIKey(got, want string) bool {
	if got == "" || want == "" {
		return false
	}
	if len(got) != len(want) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}
