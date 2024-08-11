package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexsasharegan/go-rss-agg-example/internal/database"
	"github.com/google/uuid"
)

func (apiConf *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest,
			fmt.Errorf("malformed JSON payload: %s", err),
		)
		return
	}

	feed, err := apiConf.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("failed to create feed: %s", err),
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func (apiConf *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConf.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Errorf("error retrieving feeds: %s", err),
		)
		return
	}

	payload := make([]Feed, len(feeds))
	for i, f := range feeds {
		payload[i] = databaseFeedToFeed(f)
	}

	respondWithJSON(w, http.StatusOK, payload)
}
