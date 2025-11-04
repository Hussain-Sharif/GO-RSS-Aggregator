package main

import "net/http"

func (apiCfg *apiConfig)handlerCreaeteUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Name string `name`
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}