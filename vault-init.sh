#!/bin/sh
set -e

echo "=== Starting Vault initialization ==="

export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"

echo "Waiting for Vault to be ready..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if vault status >/dev/null 2>&1; then
        echo "Vault is ready!"
        break
    fi
    attempt=$((attempt + 1))
    echo "Attempt $attempt/$max_attempts - Vault not ready yet, waiting..."
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "ERROR: Vault failed to start after $max_attempts attempts"
    exit 1
fi

echo "Checking KV v2 secrets engine..."
if ! vault secrets list | grep -q "secret/"; then
    echo "Enabling KV v2 secrets engine..."
    vault secrets enable -path=secret kv-v2
    echo "✓ KV v2 secrets engine enabled"
else
    echo "✓ KV v2 secrets engine already exists"
fi

# Check and load RSA keys
echo "Checking RSA keys..."
if [ -f "/vault/keys/private.key" ] && [ -f "/vault/keys/public.key" ]; then
    echo "✓ RSA keys found"

    # Read the keys
    echo "Reading RSA keys..."
    PRIVATE_KEY=$(cat /vault/keys/private.key)
    PUBLIC_KEY=$(cat /vault/keys/public.key)

    # Store JWT keys in Vault
    echo "Storing JWT keys in Vault..."
    vault kv put secret/jwt \
        private_key="$PRIVATE_KEY" \
        public_key="$PUBLIC_KEY"

    echo "✓ JWT keys successfully stored in Vault!"

    # Verify storage
    echo "Verifying stored keys..."
    if vault kv get secret/jwt >/dev/null 2>&1; then
        echo "✓ Keys verification successful"
        vault kv get secret/jwt
    else
        echo "✗ Keys verification failed"
        exit 1
    fi
else
    echo "✗ ERROR: RSA keys not found at /vault/keys/"
    echo "Expected files:"
    echo "  - /vault/keys/private.key"
    echo "  - /vault/keys/public.key"
    echo "Current contents of /vault/keys/:"
    ls -la /vault/keys/ || echo "Directory not accessible"
    exit 1
fi

echo "=== Vault initialization completed successfully! ==="