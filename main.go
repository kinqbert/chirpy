package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getFileserverHits() http.Handler {
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

func (cfg *apiConfig) resetFileserverHits() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Swap(0)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	})
}

func main() {
	var cfg apiConfig
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("GET /admin/metrics", cfg.getFileserverHits())
	mux.Handle("POST /admin/reset", cfg.resetFileserverHits())

	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	log.Println("Starting server...")
	log.Println(server.ListenAndServe())
}

type ValidateChirpBody struct {
	Body string `json:"body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func returnError(errorMessage string, w http.ResponseWriter) {
	errorResponse := ErrorResponse{Error: errorMessage}

	w.WriteHeader(http.StatusBadRequest)
	errorResponseJson, err := json.Marshal(errorResponse)

	if err != nil {
		w.Write([]byte(""))
	}

	w.Write(errorResponseJson)
}

func filterProfanes(text string) string {
	profanes := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	fmt.Printf("%s", text)

	words := strings.Fields(text)

	for i, word := range words {
		normalized := strings.ToLower(word)

		if _, isProfane := profanes[normalized]; isProfane {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body ValidateChirpBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		returnError("Something went wrong", w)
		return
	}

	if len(body.Body) > 140 {
		returnError("Chirp is too long", w)
		return
	}

	cleanedChirp := filterProfanes(body.Body)

	successResponse := struct {
		CleanedBody string `json:"cleaned_body"`
	}{CleanedBody: cleanedChirp}
	encodedResponse, err := json.Marshal(successResponse)

	if err != nil {
		returnError("Something went wrong", w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(encodedResponse)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
