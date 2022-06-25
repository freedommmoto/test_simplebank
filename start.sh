#!/bin/sh

set -e

source /app/app.env
echo "run db migration"
/app/migrate -path db/migration -database "$DB_SOUECE" -verbose up

echo "start the app"
exec "$@"