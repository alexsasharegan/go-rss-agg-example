package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexsasharegan/go-rss-agg-example/internal/database"
	"github.com/google/uuid"
)

func (apiConf *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest,
			fmt.Errorf("malformed JSON payload: %s", err),
		)
		return
	}

	user, err := apiConf.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("failed to create user: %s", err),
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiConf *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiConf *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiConf.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Errorf("failed to get posts: %s", err),
		)
		return
	}

	payload := make([]Post, len(posts))
	for i, p := range posts {
		payload[i] = databasePostToPost(p)
	}

	respondWithJSON(w, http.StatusOK, payload)
}
