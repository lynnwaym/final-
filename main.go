package main

import (
	"final_project/pkg/db"
	"final_project/pkg/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

//Если проверяет Антон Саттаров, то спасибо за все ревью. Даже в зелёных проектах указывал на ошибки

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env NOT FOUND")
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)

	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
