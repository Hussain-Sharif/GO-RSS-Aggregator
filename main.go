package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	// "github.com/go-chi/cors"
)


func main(){
	fmt.Println("Hello World")

	godotenv.Load(".env")

	portString:=os.Getenv("PORT")
	if(portString==""){
		log.Fatal("Port is not found") // This will stop the program to execute
	}

	router:= chi.NewRouter()

	srv:= &http.Server{
		Handler:router,
		Addr:":"+portString,
	}

	err=srv.ListenAndServe()
	if(err!=nil){
		log.Fatal(err)
	}

	fmt.Println(portString)
}