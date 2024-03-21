include .env

run:
	go run cmd/webserver/main.go -addr=:6969 -dbUrl=${DB_URL}
migration-down:
	 goose -dir internal/database/schema postgres postgres://example:example@localhost:5432/racook down 
migration-up:
	 goose -dir internal/database/schema postgres postgres://example:example@localhost:5432/racook up

