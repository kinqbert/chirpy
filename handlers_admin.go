package main

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) handleGetFileserverHits() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
	})
}

func (cfg *ApiConfig) handleReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Swap(0)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if cfg.isDev {
			err := cfg.dbQueries.DeleteAllUsers(r.Context())

			if err != nil {
				writeError(fmt.Sprintf("Failed to remove users: %e", err), w)
			}
		}

		w.WriteHeader(http.StatusOK)
	})
}
