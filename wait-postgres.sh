#!/bin/sh
# wait-postgres.sh

set -e

host="127.0.0.1:5432"
shift

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
cd /cmd/mfc-api/ && ./mfc-api