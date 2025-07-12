#!/bin/sh

set -e

# Executa as migrações
# Adicionado um loop de espera para o banco de dados
echo "Waiting for postgres..."
while ! nc -z db 5432; do
  sleep 0.1
done
echo "PostgreSQL started"

/go/bin/migrate -path /app/internal/infrastructure/database/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

# Inicia a aplicação
exec /app/cultivo-api-go