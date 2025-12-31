#!/bin/bash

set -e

SCRIPT_DIR="/setup/scripts"

echo "Setting up Google Identity Provider in Keycloak..."

KC_ADMIN_USER="${KEYCLOAK_ADMIN:-admin}"
KC_ADMIN_PASS="${KEYCLOAK_ADMIN_PASSWORD:-admin}"
KC_SERVER="http://keycloak:${KEYCLOAK_INNER_PORT:-8080}"
KC_REALM="${KEYCLOAK_REALM:-parking-users}"

GOOGLE_CLIENT_ID="${GOOGLE_OAUTH_CLIENT_ID}"
GOOGLE_CLIENT_SECRET="${GOOGLE_OAUTH_CLIENT_SECRET}"

if [ -z "$GOOGLE_CLIENT_ID" ] || [ -z "$GOOGLE_CLIENT_SECRET" ]; then
    echo "WARNING: GOOGLE_OAUTH_CLIENT_ID or GOOGLE_OAUTH_CLIENT_SECRET not set. Skipping Google IdP setup."
    exit 0
fi

BACKEND_DOMAIN="${BACKEND_DOMAIN:-backend.parking-net.space}"
AUTH_SERVICE_URL="${AUTH_SERVICE_URL:-http://auth:8800}"
KEYCLOAK_FRONTEND_URL="${KEYCLOAK_FRONTEND_URL:-https://keycloak.backend.parking-net.space}"
KEYCLOAK_BROKER_REDIRECT_URI="${KEYCLOAK_FRONTEND_URL}/realms/${KC_REALM}/broker/google/endpoint"

echo "Keycloak broker redirect URI: $KEYCLOAK_BROKER_REDIRECT_URI"
echo "IMPORTANT: This URI must be registered in Google OAuth application as Authorized redirect URI!"

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

echo "Checking if Google IdP already exists..."
EXISTING_IDP=$(docker exec "$KC_CONTAINER" bash -c \
    "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh config credentials \
     --server $KC_SERVER \
     --realm master \
     --user '$KC_ADMIN_USER' \
     --password '$KC_ADMIN_PASS' >/dev/null 2>&1 && \
     /opt/keycloak/bin/kcadm.sh get identity-provider/instances \
     -r $KC_REALM \
     --fields alias \
     --format csv \
     --noquotes 2>/dev/null | grep -i google || echo ''")

if [ -n "$EXISTING_IDP" ]; then
    echo "Google IdP already exists, updating..."
    docker exec "$KC_CONTAINER" bash -c \
        "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
         /opt/keycloak/bin/kcadm.sh config credentials \
         --server $KC_SERVER \
         --realm master \
         --user '$KC_ADMIN_USER' \
         --password '$KC_ADMIN_PASS' >/dev/null 2>&1 && \
         /opt/keycloak/bin/kcadm.sh update identity-provider/instances/google \
         -r $KC_REALM \
         -s clientId='$GOOGLE_CLIENT_ID' \
         -s clientSecret='$GOOGLE_CLIENT_SECRET' \
         -s 'config.redirectUri'='$KEYCLOAK_BROKER_REDIRECT_URI' \
         -s 'config.acceptsPromptNoneForwardFromClient'=false" >/dev/null 2>&1
    echo "Google IdP updated successfully"
else
    echo "Creating Google IdP..."
    docker exec "$KC_CONTAINER" bash -c \
        "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
         /opt/keycloak/bin/kcadm.sh config credentials \
         --server $KC_SERVER \
         --realm master \
         --user '$KC_ADMIN_USER' \
         --password '$KC_ADMIN_PASS' >/dev/null 2>&1 && \
         /opt/keycloak/bin/kcadm.sh create identity-provider/instances \
         -r $KC_REALM \
         -s alias=google \
         -s providerId=google \
         -s enabled=true \
         -s 'config.clientId'='$GOOGLE_CLIENT_ID' \
         -s 'config.clientSecret'='$GOOGLE_CLIENT_SECRET' \
         -s 'config.redirectUri'='$KEYCLOAK_BROKER_REDIRECT_URI' \
         -s 'config.acceptsPromptNoneForwardFromClient'=false" >/dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        echo "Google IdP created successfully"
    else
        echo "ERROR: Failed to create Google IdP"
        exit 1
    fi
fi

echo "Google Identity Provider setup complete!"

