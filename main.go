package main

import (
	"log"

	"github.com/IamStubborN/filmtracker/gsrv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("File .env not found, reading configuration from ENV")
	}
}

func main() {
	server := gsrv.CreateServer()
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
