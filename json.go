package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed JSON marshal: %v\n", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		log.Println("Failed response write:", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	if code >= http.StatusInternalServerError {
		log.Printf("Responding with status=%v for error: %v\n", code, err)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{
		Error: err.Error(),
	})
}
