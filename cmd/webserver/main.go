package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"racook/internal/database"
	"racook/internal/models"
	"racook/internal/server"

	_ "github.com/lib/pq"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dbURL := flag.String("dbUrl", "postgres://example:example@localhost:5432/racook?sslmode=disable", "database url")

	flag.Parse()

	connection, err := database.NewDatabaseConnection(*dbURL)
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
	defer connection.Db.Close()

	app := &server.Application{
		Recipes:        &models.RecipeModel{DB: connection.Db},
		SessionManager: connection.SessionManager,
	}

	log.Printf("Server is running on %s...", *addr)
	err = http.ListenAndServe(*addr, app.Routes())
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
