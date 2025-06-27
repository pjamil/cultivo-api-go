build:
	go build -o bin/cultivo-api ./cmd/cultivo-api-go

test:
	go test -v ./...

docker-build:
	docker build -t cultivo-api .