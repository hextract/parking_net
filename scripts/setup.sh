#!/bin/bash

set -e

SCRIPT_DIR="/setup/scripts"
cd /setup

echo "=========================================="
echo "Parking Net - Service Setup"
echo "=========================================="

if [ -z "${KEYCLOAK_CLIENT_SECRET:-}" ] || [ "${KEYCLOAK_CLIENT_SECRET}" = "your-secret-key-here" ]; then
    echo "ERROR: KEYCLOAK_CLIENT_SECRET is not set or using default value"
    echo "Please set a secure value in your .env file"
    exit 1
fi

echo ""
echo "Running service setup..."
bash "$SCRIPT_DIR/setup_services.sh"

echo ""
echo "=========================================="
echo "Setup complete!"
echo "=========================================="
