build: swagger-gen
	go build -o bin/cultivo-api ./cmd/cultivo-api-go

test-up:
	docker-compose -f docker-compose.test.yml up -d

test-down:
	docker-compose -f docker-compose.test.yml down

test: test-up
	@echo "Waiting for test database to be healthy..."
	@while [ "$(docker inspect -f '{{.State.Health.Status}}' cultivo-api-go-test-db)" != "healthy" ]; do sleep 1; done
	@echo "Test database is healthy!"
	go test ./...
	$(MAKE) test-down

swagger-gen:
	GO111MODULE=on swag init -g cmd/cultivo-api-go/main.go -parseDependency

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/infrastructure/database/migrations -seq $name

migrate-up:
	migrate -path internal/infrastructure/database/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migrate-up-test:
	migrate -path internal/infrastructure/database/migrations -database "postgres://$(DB_USER_TEST):$(DB_PASSWORD_TEST)@$(DB_HOST_TEST):$(DB_PORT_TEST)/$(DB_NAME_TEST)?sslmode=disable" -verbose up

migrate-down:
	migrate -path internal/infrastructure/database/migrations -database "postgres://postgres:753951465827PJamil@localhost:5432/cultivo-api-go?sslmode=disable" -verbose down

gitea-sessao:
	go run gitea_sessao.go
