package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)


func main(){
	fmt.Println("Hello World")

	godotenv.Load(".env")

	portString:=os.Getenv("PORT")
	if(portString==""){
		log.Fatal("Port is not found") // This will stop the program to execute
	}

	fmt.Println(portString)
}