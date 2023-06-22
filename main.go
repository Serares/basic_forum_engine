package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Print("Can't load .env file")
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	db, err := CreateDbConnection()
	if err != nil {
		log.Fatal("Error on db connection")
	}
	server := NewApiServer(":"+PORT, db)
	server.Run()
}
