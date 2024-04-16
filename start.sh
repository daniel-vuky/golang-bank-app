#!/bin/sh

set -e

echo "Run the DB migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting the application"
exec "$@"