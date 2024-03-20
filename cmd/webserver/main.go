package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ayhonz/racook"
	"github.com/ayhonz/racook/internal/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	addr := flag.String("addr", ":6969", "HTTP network address")
	dbURL := flag.String("dbUrl", "postgres://example:example@localhost:5432/racook", "database url")

	flag.Parse()

	db, err := sqlx.Open("postgres", *dbURL)
	if err != nil {
		log.Fatalf("unable to connect to database %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to ping database %v", err)
	}

	storage := database.NewStorage(db)

	server := cookbook.NewCookBookServer(storage)

	log.Println("Starting server on", *addr)
	if err := http.ListenAndServe(*addr, server); err != nil {
		log.Fatalf("could not listen on port %s %v", *addr, err)
	}
}
