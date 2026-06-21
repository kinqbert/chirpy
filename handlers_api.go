package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CreateUserBody struct {
	Email string `json:"email"`
}

func (cfg *ApiConfig) handleCreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body CreateUserBody
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)

		if err != nil {
			writeError("Something went wrong", w)
			return
		}

		dbUser, err := cfg.dbQueries.CreateUser(r.Context(), body.Email)
		if err != nil {
			fmt.Printf("%v", err)
			writeError("Failed to create user", w)
			return
		}

		user := User{
			ID:        dbUser.ID.String(),
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt.Format(time.RFC3339),
			UpdatedAt: dbUser.UpdatedAt.Format(time.RFC3339),
		}

		writeResponse(w, user, http.StatusCreated)
	})
}

type ValidateChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
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

type ValidateChirpBody struct {
	Body string `json:"body"`
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body ValidateChirpBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		writeError("Something went wrong", w)
		return
	}

	if len(body.Body) > 140 {
		writeError("Chirp is too long", w)
		return
	}

	cleanedChirp := filterProfanes(body.Body)
	response := ValidateChirpResponse{CleanedBody: cleanedChirp}

	err = writeResponse(w, response)

	if err != nil {
		writeError("Something went wrong", w)
		return
	}
}
