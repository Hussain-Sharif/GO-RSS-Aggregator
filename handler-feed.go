package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct{
		Name string `json:"name"`
		Url string `json:"url"`
	}

	params:= parameters{}

	err:=respondWithBody(w,r,&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// VALIDATE IT'S A REAL FEED!
    if !isValidFeed(params.Url) {
        // Try auto-discovery
        result, err := discoverFeeds(params.Url)
        if err != nil || len(result.FeedURLs) == 0 {
            respondWithError(w, 400, fmt.Sprintf(
                "Invalid RSS feed URL. Try /v1/discover-feed to find valid feeds from this URL: %s", 
                params.Url,
            ))
            return
        }
        
        // Suggest the first discovered feed
        respondWithError(w, 400, fmt.Sprintf(
            "Not a valid feed URL. Did you mean: %s? Use /v1/discover-feed first.",
            result.FeedURLs[0].URL,
        ))
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



// POST /v1/discover-feed - Discover RSS feeds from URL
func (cfg *apiConfig) handlerDiscoverFeed(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        URL string `json:"url"`
    }
    
    var params parameters
    if err := respondWithBody(w, r, &params); err != nil {
        respondWithError(w, 400, "Invalid JSON")
        return
    }
    
    if params.URL == "" {
        respondWithError(w, 400, "URL is required")
        return
    }
    
    // Validate/fixing URL format
    if !strings.HasPrefix(params.URL, "http") {
        params.URL = "https://" + params.URL
    }
    
    result, err := discoverFeeds(params.URL)
    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Discovery failed: %v", err))
        return
    }
    
    respondWithJSON(w, 200, result)
}
