#!/bin/bash

set -e
set -u

function create_user_and_database() {
    local database=$1
    echo "  Creating user and database '$database'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        DO \$\$
        BEGIN
            IF NOT EXISTS (SELECT FROM pg_user WHERE usename = '$database') THEN
                CREATE USER $database;
            END IF;
        END
        \$\$;
        SELECT 'CREATE DATABASE $database'
        WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$database')\gexec
        GRANT ALL PRIVILEGES ON DATABASE $database TO $database;
EOSQL
    echo "  Creating tables in '$database'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$database" -f "/docker-entrypoint-initdb.d/init_sql/init_${database}.sql" 2>/dev/null || true
}
  echo

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
    echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
    for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
        create_user_and_database $db
    done
    echo "Multiple databases created"
fi
