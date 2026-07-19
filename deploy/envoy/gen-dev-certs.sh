#!/bin/bash
# DEVELOPMENT ONLY - DO NOT USE IN PRODUCTION
# Generates self-signed certificates for local Envoy testing

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CERT_DIR="$DIR/certs"

mkdir -p "$CERT_DIR"

echo "Generating self-signed dev certificates in $CERT_DIR..."

openssl req -x509 -newkey rsa:4096 -keyout "$CERT_DIR/server.key" -out "$CERT_DIR/server.crt" \
  -days 365 -nodes -subj "/C=US/ST=State/L=City/O=AEGIS/CN=localhost"

echo "Dev certificates generated successfully."
