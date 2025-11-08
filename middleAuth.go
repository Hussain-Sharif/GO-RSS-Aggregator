package main

import (
	"fmt"
	"net/http"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/auth"
	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request,database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) (http.HandlerFunc){
	return func(w http.ResponseWriter, r *http.Request){
		apiKey,err:=auth.GetAPIKey(r.Header)
		if(err!=nil){
			respondWithError(w,403,fmt.Sprintf("Auth Error: %v",err))
			return 
		}

		user,err:=cfg.DB.GetUserByAPIKEY(r.Context(),apiKey)
		if err!=nil{
			respondWithError(w,400,fmt.Sprintf("Couldn't get user: %v",err))
			return 
		}
		handler(w,r,user)
	}
}