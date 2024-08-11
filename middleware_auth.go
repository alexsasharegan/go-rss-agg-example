package main

import (
	"fmt"
	"net/http"

	"github.com/alexsasharegan/go-rss-agg-example/internal/auth"
	"github.com/alexsasharegan/go-rss-agg-example/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConf *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.ExtractAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized,
				fmt.Errorf("failed authorization: %s", err),
			)
			return
		}

		user, err := apiConf.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound,
				fmt.Errorf("user not found: %s", err),
			)
			return
		}

		handler(w, r, user)
	}
}
