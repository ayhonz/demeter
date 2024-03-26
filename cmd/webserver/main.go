package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/ayhonz/racook"
	"github.com/ayhonz/racook/internal/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	addr := flag.String("addr", ":6969", "HTTP network address")
	dbURL := flag.String("dbUrl", "postgres://example:example@localhost:5432/racook", "database url")

	flag.Parse()

	postgres, err := sql.Open("postgres", *dbURL)
	if err != nil {
		log.Fatalf("unable to connect to database %v", err)
	}
	defer postgres.Close()

	db := sqlx.NewDb(postgres, "postgres")

	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to ping database %v", err)
	}

	storage := database.NewStorage(db)
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(postgres)
	sessionManager.Lifetime = 12 * time.Hour

	server := cookbook.NewCookBookServer(storage, sessionManager)

	log.Println("Starting server on", *addr)
	if err := http.ListenAndServe(*addr, server); err != nil {
		log.Fatalf("could not listen on port %s %v", *addr, err)
	}
}
