package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Here the _ helps to saying to go that include this code even though we are not using it!
)

type apiConfig struct{
	DB *database.Queries
}


func main(){
	fmt.Println("Hello World")

	// env setup
	godotenv.Load(".env") // This will load the .env file
	portString:=os.Getenv("PORT") // Getting the PORT from the .env file
	if(portString==""){ // if don't find the PORT
		log.Fatal("Port is not found") // This will stop the program to execute
	}
	dbUrl:=os.Getenv("DB_URL") 
	if(dbUrl==""){ 
		log.Fatal("DB Url is missing from env")
	}

	conn,err:=sql.Open("postgres",dbUrl)

	if(err!=nil){
		log.Fatal("error to connect with DB",err)
	}

	

	apiCfg:=apiConfig{ // this helps to have handlers to use the DB
		DB:database.New(conn),
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
	v1Router.Post("/users",apiCfg.handlerCreaeteUser)

	router.Mount("/v1",v1Router)

	srv:= &http.Server{    // Creating a new server with that created router
		Handler:router,
		Addr:":"+portString,
	}

	log.Printf("Server Starting at port %v", portString) 
	srvErr:=srv.ListenAndServe()
	if(srvErr!=nil){
		log.Fatal(srvErr)
	}

	fmt.Println(portString)
}