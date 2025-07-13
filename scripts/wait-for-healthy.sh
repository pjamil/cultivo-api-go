#!/bin/bash
set -e

CONTAINER_NAME=$1
TIMEOUT=60

echo "Waiting for container '$CONTAINER_NAME' to be healthy..."

for i in $(seq $TIMEOUT); do
  STATUS=$(docker inspect -f '{{.State.Health.Status}}' "$CONTAINER_NAME" 2>/dev/null || echo "starting")
  if [ "$STATUS" = "healthy" ]; then
    echo "Container '$CONTAINER_NAME' is healthy!"
    exit 0
  fi
  echo "Current status: $STATUS. Waiting... ($i/$TIMEOUT)"
  sleep 1
done

echo "Error: Timeout waiting for container '$CONTAINER_NAME' to be healthy."
exit 1
