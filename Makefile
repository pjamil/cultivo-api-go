build: swagger-gen
	go build -o bin/cultivo-api ./cmd/cultivo-api-go

swagger-gen:
	GO111MODULE=on swag init -g cmd/cultivo-api-go/main.go -parseDependency

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/infrastructure/database/migrations -seq $name

migrate-up:
	migrate -path internal/infrastructure/database/migrations -database "postgres://postgres:753951465827PJamil@localhost:5432/cultivo-api-go?sslmode=disable" -verbose up

migrate-down:
	migrate -path internal/infrastructure/database/migrations -database "postgres://postgres:753951465827PJamil@localhost:5432/cultivo-api-go?sslmode=disable" -verbose down
