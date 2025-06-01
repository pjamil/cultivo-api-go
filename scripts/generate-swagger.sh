#!/bin/bash

echo "Generating Swagger docs..."
swag init -g cmd/cultivo-api-go-swagger/main.go -o ./docs