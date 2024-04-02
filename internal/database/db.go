package database

import (
	"database/sql"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	Db             *sqlx.DB
	SessionManager *scs.SessionManager
}

func NewDatabaseConnection(dbURL string) (*Database, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(db, "postgres")

	err = dbx.Ping()
	if err != nil {
		return nil, err
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour

	return &Database{
		Db:             dbx,
		SessionManager: sessionManager,
	}, nil
}
