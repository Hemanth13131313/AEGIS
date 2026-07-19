# AEGIS Vault Policy
# Grants AEGIS services read access to their respective secret paths

# Gateway: read JWKS signing key, Kafka credentials
path "secret/data/aegis/gateway/*" {
  capabilities = ["read"]
}

# Policy Engine: read DB credentials
path "secret/data/aegis/policy-engine/*" {
  capabilities = ["read"]
}

# Scanner: read model API keys
path "secret/data/aegis/scanner/*" {
  capabilities = ["read"]
}

# All services: read shared credentials
path "secret/data/aegis/shared/*" {
  capabilities = ["read"]
}

# Deny everything else
path "*" {
  capabilities = ["deny"]
}
