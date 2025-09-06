
DB_DSN := "postgres://user:password@localhost:5438/rest_api_db?sslmode=disable"

## Start the Postgres container
run-db-container:
	docker-compose up -d db

## Stop the Postgres container
stop-db-container:
	docker-compose down

## Apply the migration
migrate-up: run-db-container
	migrate -path ./migrations -database $(DB_DSN) up

## Rollback the migration
migrate-down: run-db-container
	migrate -path ./migrations -database $(DB_DSN) down

## Build the Api
build:
	go build -C cmd/api -o ../../bin/go-rest-api

## Run the Api
run:
	cd cmd/api && go run .

## Run unit tests
test-unit:
	go test -tags 'unit' -v ./internal/core/

## Run integration tests
test-integration:
	go test -tags 'integration' -v ./internal/handlers/... && \
	go test -tags 'integration' -v ./internal/db/...