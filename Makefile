include .env

run:
	go run cmd/webserver/main.go -addr=:6969 -dbUrl=${DB_URL}
migration:
	 goose postgres postgres://example:example@localhost:5432/racook up

