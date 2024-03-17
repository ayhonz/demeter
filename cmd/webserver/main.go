package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/ayhonz/racook"
	"github.com/ayhonz/racook/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	addr := flag.String("addr", ":6969", "HTTP network address")
	dbURL := flag.String("dbUrl", "postgres://example:example@localhost:5432/racook", "database url")

	flag.Parse()

	conn, err := sql.Open("postgres", *dbURL)
	if err != nil {
		log.Fatalf("unable to connect to database %v", err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatalf("unable to ping database %v", err)
	}

	queries := database.New(conn)

	server := cookbook.NewCookBookServer(queries)

	log.Println("Starting server on", *addr)
	if err := http.ListenAndServe(*addr, server); err != nil {
		log.Fatalf("could not listen on port %s %v", *addr, err)
	}
}
