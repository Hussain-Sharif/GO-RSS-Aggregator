package main

import (
	"fmt"
	"log"
	"os"
)


func main(){
	fmt.Println("Hello World")

	portString:=os.Getenv("PORT")
	if(portString==""){
		log.Fatal("Port is not found") // This will stop the program to execute
	}

	fmt.Println(portString)
}