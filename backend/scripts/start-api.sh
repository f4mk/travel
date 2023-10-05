#!/bin/sh

set -e

set -o allexport
source $CONFIG_PATH_ENV
set +o allexport

# TODO: move to app and wait there (think where to run migrations)
/wait-for-it.sh $RABBIT_HOST_NAME:$RABBIT_PORT
/wait-for-it.sh $PG_HOST_NAME:$PG_PORT
until PGPASSWORD=$PG_PASSWORD psql -h $PG_HOST_NAME -p $PG_PORT -U $PG_USER -lqt | cut -d \| -f 1 | grep -qw travel; do
  echo "Waiting for the 'travel' database to be available..."
  sleep 5
done

echo "'travel' database is available."
/wait-for-it.sh $RD_HOST_NAME:$RD_PORT

export PGPASSWORD=$PG_PASSWORD
psql -h $PG_HOST_NAME -p $PG_PORT -U "$PG_USER" -d $PG_DB_NAME -c "CREATE EXTENSION IF NOT EXISTS postgis;"

./migrate -migrate

./app
