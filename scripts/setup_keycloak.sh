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

BACKEND_DOMAIN="${BACKEND_DOMAIN:-backend.parking-net.space}"
AUTH_SERVICE_URL="${AUTH_SERVICE_URL:-http://auth:8800}"
GOOGLE_REDIRECT_URI="${GOOGLE_OAUTH_REDIRECT_URI:-}"

if [ -z "$GOOGLE_REDIRECT_URI" ]; then
    if [ -n "$BACKEND_DOMAIN" ] && [[ "$BACKEND_DOMAIN" != "localhost"* ]] && [[ "$BACKEND_DOMAIN" != "127.0.0.1"* ]]; then
        GOOGLE_REDIRECT_URI="https://${BACKEND_DOMAIN}/auth/google/callback"
    else
        GOOGLE_REDIRECT_URI="${AUTH_SERVICE_URL}/auth/google/callback"
    fi
fi

FRONTEND_URL="${FRONTEND_URL:-http://localhost:3000}"
FRONTEND_CALLBACK="${FRONTEND_URL}/auth/callback"

echo "Configuring client secret and redirect URIs..."
echo "Google redirect URI: $GOOGLE_REDIRECT_URI"
echo "Frontend callback: $FRONTEND_CALLBACK"

echo "Updating client configuration..."
docker exec "$KC_CONTAINER" bash -c \
    "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh config credentials \
     --server $KC_SERVER \
     --realm master \
     --user '$KC_ADMIN_USER' \
     --password '$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh update clients/$CLIENT_ID \
     -r $KC_REALM \
     -s secret='$KC_CLIENT_SECRET' \
     -s 'redirectUris=[\"$GOOGLE_REDIRECT_URI\",\"http://localhost:8800/auth/google/callback\",\"$FRONTEND_CALLBACK\"]' \
     -s 'webOrigins=[\"+\"]'"

UPDATE_RESULT=$?

if [ $UPDATE_RESULT -ne 0 ]; then
    echo "WARNING: First attempt failed (exit code: $UPDATE_RESULT), trying alternative method..."
    docker exec "$KC_CONTAINER" bash -c \
        "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
         /opt/keycloak/bin/kcadm.sh config credentials \
         --server $KC_SERVER \
         --realm master \
         --user '$KC_ADMIN_USER' \
         --password '$KC_ADMIN_PASS' && \
         /opt/keycloak/bin/kcadm.sh update clients/$CLIENT_ID \
         -r $KC_REALM \
         -s secret='$KC_CLIENT_SECRET' \
         -s 'redirectUris[0]'='$GOOGLE_REDIRECT_URI' \
         -s 'redirectUris[1]'='http://localhost:8800/auth/google/callback' \
         -s 'redirectUris[2]'='$FRONTEND_CALLBACK' \
         -s 'webOrigins[0]'='+'"
    
    if [ $? -ne 0 ]; then
        echo "ERROR: Both methods failed. Please run scripts/fix_keycloak_redirect.sh manually"
        exit 1
    fi
fi

echo "Client secret and redirect URIs configured successfully"
echo "Setup complete!"
