#!/bin/bash

set -e

echo "Fixing Keycloak redirect URIs..."

KC_ADMIN_USER="${KEYCLOAK_ADMIN:-admin}"
KC_ADMIN_PASS="${KEYCLOAK_ADMIN_PASSWORD:-admin}"
KC_SERVER="http://keycloak:${KEYCLOAK_INNER_PORT:-8080}"
KC_REALM="${KEYCLOAK_REALM:-parking-users}"
KC_CLIENT_ID="${KEYCLOAK_CLIENT:-parking-auth}"
KC_CONTAINER="keycloak"

BACKEND_DOMAIN="${BACKEND_DOMAIN:-backend.parking-net.space}"
GOOGLE_REDIRECT_URI="https://${BACKEND_DOMAIN}/auth/google/callback"
FRONTEND_URL="${FRONTEND_URL:-https://parking-net.space}"
FRONTEND_CALLBACK="${FRONTEND_URL}/auth/callback"

echo "Getting client ID..."
CLIENT_ID=$(docker exec "$KC_CONTAINER" bash -c \
    "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh config credentials \
     --server $KC_SERVER \
     --realm master \
     --user '$KC_ADMIN_USER' \
     --password '$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh get clients \
     -r $KC_REALM \
     -q clientId=$KC_CLIENT_ID \
     --fields id \
     --format csv \
     --noquotes | tail -1")

if [ -z "$CLIENT_ID" ]; then
    echo "ERROR: Could not find client '$KC_CLIENT_ID' in realm '$KC_REALM'"
    exit 1
fi

echo "Found client ID: $CLIENT_ID"
echo "Setting redirect URIs..."
echo "  - $GOOGLE_REDIRECT_URI"
echo "  - http://localhost:8800/auth/google/callback"
echo "  - $FRONTEND_CALLBACK"

docker exec "$KC_CONTAINER" bash -c \
    "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh config credentials \
     --server $KC_SERVER \
     --realm master \
     --user '$KC_ADMIN_USER' \
     --password '$KC_ADMIN_PASS' && \
     /opt/keycloak/bin/kcadm.sh update clients/$CLIENT_ID \
     -r $KC_REALM \
     -s 'redirectUris=[\"$GOOGLE_REDIRECT_URI\",\"http://localhost:8800/auth/google/callback\",\"$FRONTEND_CALLBACK\"]' \
     -s 'webOrigins=[\"+\"]'"

if [ $? -eq 0 ]; then
    echo "✓ Redirect URIs updated successfully"
    
    echo "Verifying redirect URIs..."
    docker exec "$KC_CONTAINER" bash -c \
        "export KC_ADMIN='$KC_ADMIN_USER' KC_ADMIN_PASSWORD='$KC_ADMIN_PASS' && \
         /opt/keycloak/bin/kcadm.sh config credentials \
         --server $KC_SERVER \
         --realm master \
         --user '$KC_ADMIN_USER' \
         --password '$KC_ADMIN_PASS' && \
         /opt/keycloak/bin/kcadm.sh get clients/$CLIENT_ID \
         -r $KC_REALM \
         --fields redirectUris \
         --format json" | grep -o '"redirectUris":\[[^]]*\]' || echo "Could not verify"
else
    echo "✗ Failed to update redirect URIs"
    exit 1
fi

echo "Done!"

