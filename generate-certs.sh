#!/bin/bash

set -e

CERT_DIR="./certs"
mkdir -p "$CERT_DIR"

echo "Generating TLS certificate with SANs..."

# Create a temporary OpenSSL configuration file
cat > "$CERT_DIR/openssl.cnf" << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = RU
ST = Dev
L = Dev
O = LocalDev
CN = localhost

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF

# Generate the certificate with SANs
openssl req -newkey rsa:2048 -nodes -keyout "$CERT_DIR/server.key" \
  -x509 -days 365 -out "$CERT_DIR/server.crt" \
  -config "$CERT_DIR/openssl.cnf" \
  -extensions v3_req

# Clean up the temporary config file
rm "$CERT_DIR/openssl.cnf"

echo "Certificates created:"
echo " - $CERT_DIR/server.crt"
echo " - $CERT_DIR/server.key"
