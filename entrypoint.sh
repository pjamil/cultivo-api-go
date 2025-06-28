#!/bin/sh

set -e

# Executa as migrações
/go/bin/migrate -path /app/internal/infrastructure/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

# Inicia a aplicação
exec "$@"
