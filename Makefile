build:
	go build -o bin/cultivo-api ./cmd/cultivo-api-go

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/infrastructure/database/migrations -seq $name

migrate-up:
	migrate -path internal/infrastructure/database/migrations -database "postgres:postgre://localhost:5432/cultivo-api-go" -verbose up

migrate-down:
	migrate -path internal/infrastructure/database/migrations -database "postgres://postgres:postgres@localhost:5432/cultivo-api-go?sslmode=disable" -verbose down
