APP_NAME=todo-api-go
DB_URL=postgres://sonha:sonha@localhost:5432/todo?sslmode=disable

# Run the Go server
run:
	go run ./cmd/api

# Build binary
build:
	go build -o bin/$(APP_NAME) ./cmd/api

# Start Docker containers
up:
	docker compose up -d

# Stop containers
down:
	docker compose down

# Reset DB (dangerous: drops data!)
reset-db:
	docker compose down -v
	docker compose up -d

# Run migrations
migrate:
	goose -dir ./migrations postgres "$(DB_URL)" up

# Rollback last migration
rollback:
	goose -dir ./migrations postgres "$(DB_URL)" down

# Create a new migration file: make new name=create_users_table
new:
	goose -dir ./migrations create $(name) sql

# Lint
lint:
	golangci-lint run ./...

# Generate sqlc code
sqlc:
	sqlc generate

# Swagger docs
swagger:
	swag init -g cmd/api/main.go -o docs
