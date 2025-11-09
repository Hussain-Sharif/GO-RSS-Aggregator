package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct{
		FeedID uuid.UUID `json:"feed_id"` 
	}

	decoder := json.NewDecoder(r.Body)
	params:= parameters{}

	err:=decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	createdFeedFollow,err:=cfg.DB.CreateFeed_Follows(r.Context(),database.CreateFeed_FollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	respondWithJSON(w,200,databaseFeedFollowsToFeedFollows(createdFeedFollow))

}

