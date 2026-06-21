package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kinqbert/chirpy/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL, platform := os.Getenv("DB_URL"), os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to db: %v", err)
		return
	}

	dbQueries := database.New(db)

	cfg := ApiConfig{
		dbQueries: *dbQueries,
		isDev:     platform == "dev",
	}
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("GET /admin/metrics", cfg.handleGetFileserverHits())
	mux.Handle("POST /admin/reset", cfg.handleReset())

	mux.Handle("POST /api/users", cfg.handleCreateUser())

	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	log.Println("Starting server...")
	log.Println(server.ListenAndServe())
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
