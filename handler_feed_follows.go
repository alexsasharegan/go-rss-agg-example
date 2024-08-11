package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexsasharegan/go-rss-agg-example/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiConf *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest,
			fmt.Errorf("malformed JSON payload: %s", err),
		)
		return
	}

	follow, err := apiConf.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("failed to create feed follow: %s", err),
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(follow))
}

func (apiConf *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := apiConf.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Errorf("error retrieving feed follows: %s", err),
		)
		return
	}

	payload := make([]FeedFollow, len(follows))
	for i, f := range follows {
		payload[i] = databaseFeedFollowToFeedFollow(f)
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (apiConf *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	ffIDStr := chi.URLParam(r, "feedFollowID")
	ffID, err := uuid.Parse(ffIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			fmt.Errorf("malformed feed follow id: %s", err),
		)
		return
	}

	err = apiConf.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     ffID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Errorf("failed to delete feed follow: %s", err),
		)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
