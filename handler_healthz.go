package main

import (
	"errors"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct{}{})
}

func errorzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusBadRequest, errors.New("Something went wrong."))
}
