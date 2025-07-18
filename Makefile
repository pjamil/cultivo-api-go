build: swagger-gen
	go build -o bin/cultivo-api ./cmd/cultivo-api-go

test-up:
	docker-compose -f docker-compose.test.yml down -v || true
	docker-compose -f docker-compose.test.yml up -d

test-down:
	docker-compose -f docker-compose.test.yml down

test:
	docker-compose -f docker-compose.test.yml down -v --remove-orphans
	docker volume rm cultivo-api-go_test_db_data || true
	docker-compose -f docker-compose.test.yml up -d
	./scripts/wait-for-healthy.sh cultivo-api-go-test-db
	go test ./...



swagger-gen:
	GO111MODULE=on swag init -g cmd/cultivo-api-go/main.go -parseDependency

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/infrastructure/database/migrations -seq $name

migrate-up:
	migrate -path internal/infrastructure/database/migrations -database "postgres://postgres:$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migrate-up-test:
	migrate -path internal/infrastructure/database/migrations -database "postgres://$(DB_USER_TEST):$(DB_PASSWORD_TEST)@$(DB_HOST_TEST):$(DB_PORT_TEST)/$(DB_NAME_TEST)?sslmode=disable" -verbose up

migrate-down:
	migrate -path internal/infrastructure/database/migrations -database "postgres://postgres:753951465827PJamil@localhost:5432/cultivo-api-go?sslmode=disable" -verbose down

gitea-sessao:
	go run gitea_sessao.go
