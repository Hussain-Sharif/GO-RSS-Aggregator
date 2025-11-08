package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct{
		Name string `json:"name"`
		Url string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params:= parameters{}

	err:=decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed,err:=cfg.DB.CreateFeed(r.Context(),database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url: params.Url,
		UserID: user.ID,
	})

	respondWithJSON(w,200,databaseFeedToFeed(feed))

}



func (cfg *apiConfig) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request){

	allFeeds,err:=cfg.DB.GetFeeds(r.Context())
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't get All feeds: %v",err))
		return 
	}

	respondWithJSON(w,200,databaseFeedsToFeeds(allFeeds))

}