package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)


func main(){
	fmt.Println("Hello World")

	// env setup
	godotenv.Load(".env") // This will load the .env file
	portString:=os.Getenv("PORT") // Getting the PORT from the .env file
	if(portString==""){ // if don't find the PORT
		log.Fatal("Port is not found") // This will stop the program to execute
	}

	// router creation
	router:= chi.NewRouter() // Creating a new router
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*","http://*"},
		AllowedMethods: []string{"GET","POST","PUT","DELETE","PATCH"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router:=chi.NewRouter() // It's like a subrouter
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.HandleFunc("/err",handlerError)


	router.Mount("/v1",v1Router)

	srv:= &http.Server{    // Creating a new server with that created router
		Handler:router,
		Addr:":"+portString,
	}

	log.Printf("Server Starting at port %v", portString) 
	err:=srv.ListenAndServe()
	if(err!=nil){
		log.Fatal(err)
	}

	fmt.Println(portString)
}