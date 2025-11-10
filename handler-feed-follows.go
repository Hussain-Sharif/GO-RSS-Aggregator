package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
	"github.com/go-chi/chi"
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
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't able to create the Feed: %v", err))
		return
	}

	respondWithJSON(w,200,databaseFeedFollowsToFeedFollows(createdFeedFollow))

}


func (cfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	allFeedFollows,err:=cfg.DB.GetAllFeed_Follows(r.Context(),user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't able to get all feedfollows: %v", err))
		return
	}

	respondWithJSON(w,200,databaseGetFeedFollowsToGetFeedFollows(allFeedFollows))
}


func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	
	feedFollowIDStr:=chi.URLParam(r,"feedFollowID")

	feedFollowID,err:=uuid.Parse(feedFollowIDStr)
	if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't able to parse feedFollowID: %v", err))
		return
	}



	err=cfg.DB.DeleteFeed_Follows(r.Context(),database.DeleteFeed_FollowsParams{
		ID:feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't able to Delete feedfollows: %v", err))
		return
	}

	respondWithJSON(w,200, struct{}{})
}