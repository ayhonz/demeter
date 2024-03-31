include .env

all: build

build:
	@echo "Building..."
	@templ generate
	@go build -o main cmd/webserver/main.go

# Run the application
run:
	@go run cmd/webserver/main.go -addr=:${PORT} -dbUrl=${DB_URL}

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
		air --build.bin "./main -addr=:${PORT} -dbUrl=${DB_URL}"; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi


# Run the migrations
migration-up:
	@if command -v goose > /dev/null; then \
	    goose -dir internal/database/schema postgres postgres://example:example@localhost:5432/racook up; \
	else \
	    read -p "Go's 'goose' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/pressly/goose/v3/cmd/goose@latest; \
	    	goose -dir internal/database/schema postgres postgres://example:example@localhost:5432/racook up; \
	    else \
	        echo "You chose not to install goose. Exiting..."; \
	        exit 1; \
	    fi; \
	fi


# RUN Rollback of migrations
migration-down:
	@if command -v goose > /dev/null; then \
	    goose -dir internal/database/schema postgres postgres://example:example@localhost:5432/racook down; \
	else \
	    read -p "Go's 'goose' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/pressly/goose/v3/cmd/goose@latest; \
	    	goose -dir internal/database/schema postgres postgres://example:example@localhost:5432/racook down; \
	    else \
	        echo "You chose not to install goose. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
