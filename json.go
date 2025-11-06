package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	if(code>499){
		log.Printf("Responding with 5XX err Message: %s",message)
	}

	type errorResponse struct {
		Error string `json:"error"` // json:"desiredJSONKeyName" --> t's metadata that tells the JSON encoder/decoder how to handle this field when converting between Go structs and JSON.
	}
	respondWithJSON(w, code, errorResponse{Error: message}) // here the errorResponse{} is exactly the instance of that type with proper values
}

// Below one is like a template for the boilerplate for proper json response
func respondWithJSON(w http.ResponseWriter, code int , payload interface{}) { // here interface payload accepts any type as the payload argument 
	data, err := json.Marshal(payload)
	
	if err != nil {
		log.Printf("Failed to Marshal the JSON %v",payload)
		w.WriteHeader(http.StatusInternalServerError) // here we mention the status code
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}