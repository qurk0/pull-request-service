#!/bin/sh
set -e

echo "Waiting for Postgres at ${DB_HOST}:${DB_PORT}..."

until nc -z "$DB_HOST" "$DB_PORT"; do
  sleep 1
done

echo "Postgres is up, running goose migrations..."

goose -dir /migrations \
  postgres "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}?sslmode=disable" \
  up

echo "Migrations finished successfully."
