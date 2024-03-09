package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewCookBookServer(NewInMemoryRecipeStore())
	log.Fatal(http.ListenAndServe(":8081", server))
}
