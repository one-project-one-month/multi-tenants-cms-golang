#!/bin/sh

sleep 3

export VAULT_ADDR='http://localhost:8200'
vault login $VAULT_DEV_ROOT_TOKEN_ID


vault secrets enable -path=secret kv-v2

vault kv put secret/jwt \
  private_key="$(cat /vault/keys/private.key)" \
  public_key="$(cat /vault/keys/public.key)"
