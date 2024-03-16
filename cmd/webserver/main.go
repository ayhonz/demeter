package main

import (
	"github.com/ayhonz/racook"

	"log"
	"net/http"
	"os"
)

const dbFileName = "recipes.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := cookbook.NewFileSystemRecipeStore(db)
	if err != nil {
		log.Fatalf("problem creating file system recipe store, %v ", err)
	}

	server := cookbook.NewCookBookServer(store)

	log.Println("Starting server on :6969")
	if err := http.ListenAndServe(":6969", server); err != nil {
		log.Fatalf("could not listen on port 6969 %v", err)
	}
}
