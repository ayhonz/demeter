package main

import (
	"log"
	"net/http"
)

func main() {
	server := &CookBookServer{NewInMemoryRecipeStore()}
	log.Fatal(http.ListenAndServe(":8081", server))
}
