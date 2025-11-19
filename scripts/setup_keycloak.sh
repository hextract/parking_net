#!/bin/bash

set -e

SCRIPT_DIR="/setup/scripts"

echo "Setting up Keycloak client secret..."

KC_ADMIN_USER="${KEYCLOAK_ADMIN:-admin}"
KC_ADMIN_PASS="${KEYCLOAK_ADMIN_PASSWORD:-admin}"
KC_SERVER="http://keycloak:${KEYCLOAK_INNER_PORT:-8080}"
KC_REALM="${KEYCLOAK_REALM:-parking-users}"
KC_CLIENT_ID="${KEYCLOAK_CLIENT:-parking-auth}"
KC_CLIENT_SECRET="${KEYCLOAK_CLIENT_SECRET}"

if [ -z "$KC_CLIENT_SECRET" ]; then
    echo "ERROR: KEYCLOAK_CLIENT_SECRET environment variable is not set"
    exit 1
fi

echo "Waiting for Keycloak to be ready..."
timeout=120
KC_CONTAINER="keycloak"
while ! docker exec "$KC_CONTAINER" /opt/keycloak/bin/kcadm.sh config credentials \
    --server "$KC_SERVER" \
    --realm master \
    --user "$KC_ADMIN_USER" \
    --password "$KC_ADMIN_PASS" >/dev/null 2>&1; do
    sleep 2
    timeout=$((timeout - 2))
    if [ $timeout -le 0 ]; then
        echo "ERROR: Keycloak not ready, timeout"
        exit 1
    fi
done

echo "Keycloak is ready"

CLIENT_ID=$(docker exec "$KC_CONTAINER" bash -c \
    "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh config credentials \
     --server $KC_SERVER \
     --realm master \
     --user '$KC_ADMIN_USER' \
     --password '$KC_ADMIN_PASS' >/dev/null 2>&1 && \
     /opt/keycloak/bin/kcadm.sh get clients \
     -r $KC_REALM \
     -q clientId=$KC_CLIENT_ID \
     --fields id \
     --format csv \
     --noquotes 2>/dev/null | tail -1")

if [ -z "$CLIENT_ID" ]; then
    echo "ERROR: Could not find client '$KC_CLIENT_ID' in realm '$KC_REALM'"
    exit 1
fi

echo "Found client ID: $CLIENT_ID"

docker exec "$KC_CONTAINER" bash -c \
    "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh config credentials \
     --server $KC_SERVER \
     --realm master \
     --user '$KC_ADMIN_USER' \
     --password '$KC_ADMIN_PASS' >/dev/null 2>&1 && \
     /opt/keycloak/bin/kcadm.sh update clients/$CLIENT_ID \
     -r $KC_REALM \
     -s secret='$KC_CLIENT_SECRET'" >/dev/null 2>&1

echo "Client secret configured successfully"
echo "Setup complete!"
