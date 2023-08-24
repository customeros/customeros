#!/bin/bash

# Secret name 
SECRET_NAME=user-admin-api

# Get secret data from env vars
GOOGLE_OAUTH_CLIENT_ID_BASE64_ENCODED=$(echo "$GOOGLE_OAUTH_CLIENT_ID" | tr -d '\n' | base64)
GOOGLE_OAUTH_CLIENT_SECRET_BASE64_ENCODED=$(echo "$GOOGLE_OAUTH_CLIENT_SECRET" | tr -d '\n' | base64)

# Generate secret YAML 
cat <<EOF > ${SECRET_NAME}-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: ${SECRET_NAME}
type: Opaque
data:
  GOOGLE_OAUTH_CLIENT_ID: ${GOOGLE_OAUTH_CLIENT_ID_BASE64_ENCODED}
  GOOGLE_OAUTH_CLIENT_SECRET: ${GOOGLE_OAUTH_CLIENT_SECRET_BASE64_ENCODED}
---
EOF
cat user-admin-api.yaml >> user-admin-api-secret.yaml
mv user-admin-api-secret.yaml user-admin-api.yaml
