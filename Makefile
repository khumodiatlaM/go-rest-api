DB_DSN := "postgres://user:password@localhost:5438/rest_api_db?sslmode=disable"

run-db-container:
	docker-compose up -d db

stop-db-container:
	docker-compose down

migrate-up: run-db-container
	migrate -path ./migrations -database $(DB_DSN) up

migrate-down: run-db-container
	migrate -path ./migrations -database $(DB_DSN) down