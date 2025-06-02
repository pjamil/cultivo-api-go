#!/bin/bash

echo "Generating Swagger docs..."
swag init -g cmd/cultivo-api-go/main.go -o ./docs