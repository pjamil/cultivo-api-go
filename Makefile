build:
    go build -o bin/cultivo-api ./cmd/api

test:
    go test -v ./...

docker-build:
    docker build -t cultivo-api .