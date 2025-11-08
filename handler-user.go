package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/auth"
	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreaeteUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't able to create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

// Below is the Authenticated End Point.
func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter,r *http.Request){
	apiKey,err:=auth.GetAPIKey(r.Header)
	if(err!=nil){
		respondWithError(w,403,fmt.Sprintf("Auth Error: %v",err))
		return 
	}

	user,err:=apiConfig.DB.GetUserByAPIKEY(r.Context(),apiKey)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't get user: %v",err))
		return 
	}

	respondWithJSON(w,200,databaseUserToUser(user))
}