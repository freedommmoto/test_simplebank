#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path db/migration -database "$DB_SOUECE" -verbose up

echo "start the app"
exec "$@"