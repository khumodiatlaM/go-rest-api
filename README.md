# Go REST API Starter Project

This project is a boilerplate for building a production-ready REST API in Go. It follows the **Hexagonal Architecture** pattern to ensure a clear separation of concerns, making the application scalable, maintainable, and highly testable.

***

## Key Technologies

* **Go:** The core language for the API.
* **Hexagonal Architecture:** For a clean, decoupled design (Ports and Adapters).
* **PostgreSQL:** The primary database.
* **Docker & Docker Compose:** To containerize the database for a consistent development and testing environment.
* **Viper:** For robust configuration management using `.env` files.
* **Zap:** A high-performance structured logging library.
* **Testify:** A powerful toolkit for writing clear and expressive tests (mocks and assertions).
* **Testcontainers-go:** To spin up a real database in a Docker container for isolated integration tests.
* **Golang Migrate:** A CLI tool for managing database schema migrations.
* **Juliensmit/httprouter:** A fast and simple HTTP router.
* **Go Modules:** For dependency management.
* **golang-jwt/jwt/v5:** For secure JWT authentication.
* **bcrypt:** For secure password hashing.
* **Prometheus:** For collecting and monitoring API metrics.
* **Makefile:** For automating common tasks like building, testing, and running migrations.

***

## Project Structure

The project is organized into the following directories:

```
├── cmd/                # Application entry points
├── config/             # Configuration management
├── internal/           # Application core (domain, services, repositories)
│   ├── core/           # Business logic, Domain models and interfaces
│   ├── handlers/       # HTTP handlers
│   └── db/             # Data access layer
|   └── metrics/        # Prometheus metrics setup
├── pkg/                # Shared utilities and packages  
├── docs/               # API documentation
├── migrations/         # Database migration files
├── tests/              # Integration and unit tests
├── .env.example        # Example environment variables file
├── .gitignore          # Git ignore file
├── go.mod              # Go module file
├── go.sum              # Go module checksum file
├── Makefile            # Makefile for common tasks
├── README.md           # Project documentation
├── docker-compose.yml  # Docker Compose file for setting up the database
└── Dockerfile          # Dockerfile for containerizing the application
```

***

## Getting Started

### Prerequisites

* Go (1.18 or newer)
* Docker & Docker Compose

### 1. Clone the repository

```bash
git clone <repository_url>
cd <project_directory>

```
### 2. Set up environment variables
Copy the `.env.example` file to `.env` and modify it as needed:

```bash
cp .env.example .env
```
### 3. Run the PostgreSQL database, Prometheus, and Grafana using Docker Compose

```bash
make docker-up
``` 
### 4. Run database migrations
Make sure you have `golang-migrate` installed. Then run:
```bash
make migrate-up
```
### 5. Build and run the application

```bash
make build && make run 
```
The API will be accessible at `http://localhost:8080`.

### Testing
Run unit and integration tests using:

```bash
make test-unit
```

```bash
make test-integration
```

### API Endpoints
The API provides the following endpoints:
* `GET /health`:    Check the health status of the API.
* `POST /users`:    Create a new user.
* `POST /users/login`: Authenticate a user and return a JWT token.
* `GER /users/:id`: Retrieve a user by id. **(Protected, requires JWT token)**


### Monitoring
The API is instrumented with Prometheus metrics for monitoring. You can access the metrics at http://localhost:8080/metrics.


### Documentation
For detailed API documentation, refer to the `docs/` directory.