package main

import (
	"encoding/json"
	"net/http"
)

func writeResponse(w http.ResponseWriter, response any, status ...int) error {
	if len(status) == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(status[0])
	}

	encodedResponse, err := json.Marshal(response)

	if err != nil {
		return err
	}

	w.Write(encodedResponse)
	return nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(errorMessage string, w http.ResponseWriter) {
	errorResponse := ErrorResponse{Error: errorMessage}

	w.WriteHeader(http.StatusBadRequest)
	errorResponseJson, err := json.Marshal(errorResponse)

	if err != nil {
		w.Write([]byte(""))
	}

	w.Write(errorResponseJson)
}
