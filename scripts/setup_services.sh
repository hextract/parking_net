#!/bin/bash

set -e

SCRIPT_DIR="/setup/scripts"

echo "=========================================="
echo "Starting automated service setup..."
echo "=========================================="

echo "Waiting for database to be ready..."
timeout=60
DB_CONTAINER="db"
while ! docker exec "$DB_CONTAINER" pg_isready -U "${POSTGRES_USER:-postgres}" >/dev/null 2>&1; do
    sleep 2
    timeout=$((timeout - 2))
    if [ $timeout -le 0 ]; then
        echo "ERROR: Database not ready, timeout"
        exit 1
    fi
done
echo "Database is ready"

sleep 3

echo "Verifying database tables..."
PARKING_DB="${PARKING_DB_NAME:-parking_db}"
BOOKING_DB="${BOOKING_DB_NAME:-booking_db}"

if ! docker exec "$DB_CONTAINER" psql -U "${POSTGRES_USER:-postgres}" -d "$PARKING_DB" -c "\dt" 2>/dev/null | grep -q "parking_places"; then
    echo "Creating parking_places table..."
    docker exec "$DB_CONTAINER" psql -U "${POSTGRES_USER:-postgres}" -d "$PARKING_DB" \
        -f /docker-entrypoint-initdb.d/init_sql/init_parking.sql >/dev/null 2>&1 || true
fi

if ! docker exec "$DB_CONTAINER" psql -U "${POSTGRES_USER:-postgres}" -d "$BOOKING_DB" -c "\dt" 2>/dev/null | grep -q "bookings"; then
    echo "Creating bookings table..."
    docker exec "$DB_CONTAINER" psql -U "${POSTGRES_USER:-postgres}" -d "$BOOKING_DB" \
        -f /docker-entrypoint-initdb.d/init_sql/init_booking.sql >/dev/null 2>&1 || true
fi

echo "Waiting for Keycloak to be ready..."
timeout=120
KC_PORT="${KEYCLOAK_PORT:-8080}"
while ! curl -s "http://keycloak:${KC_PORT}" >/dev/null 2>&1 && ! curl -s "http://localhost:${KC_PORT}" >/dev/null 2>&1; do
    sleep 2
    timeout=$((timeout - 2))
    if [ $timeout -le 0 ]; then
        echo "ERROR: Keycloak not ready, timeout"
        exit 1
    fi
done

echo "Setting up Keycloak..."
bash "$SCRIPT_DIR/setup_keycloak.sh"

echo "=========================================="
echo "Service setup complete!"
echo "=========================================="
