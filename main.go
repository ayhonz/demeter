package main

import (
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

	store := NewFileSystemRecipeStore(db)
	server := NewCookBookServer(store)
	if err := http.ListenAndServe(":6969", server); err != nil {
		log.Fatalf("could not listen on port 6969 %v", err)

	}
}
