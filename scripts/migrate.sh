#!/bin/bash

DB_URL="postgres://compose-postgres:compose-postgres@localhost:5432/urlshortener?sslmode=disable"
MIGRATIONS_PATH="db/migrations"

case "$1" in
  up)
    echo "Applying migrations..."
    migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" up
    ;;
  down)
    echo "Rolling back last migration..."
    migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" down 1
    ;;
  version)
    echo "Current migration version:"
    migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" version
    ;;
  force)
    if [ -z "$2" ]; then
      echo "Usage: $0 force <version>"
      exit 1
    fi
    echo "Forcing version to $2..."
    migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" force "$2"
    ;;
  create)
    if [ -z "$2" ]; then
      echo "Usage: $0 create <migration_name>"
      exit 1
    fi
    migrate create -ext sql -dir "$MIGRATIONS_PATH" -seq "$2"
    ;;
  *)
    echo "Usage: $0 {up|down|version|force <version>|create <name>}"
    exit 1
esac
