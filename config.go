package main

import (
	"net/http"
	"sync/atomic"

	"github.com/kinqbert/chirpy/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      database.Queries
	isDev          bool
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
