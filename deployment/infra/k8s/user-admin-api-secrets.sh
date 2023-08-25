#!/bin/bash

# Secret name 
SECRET_NAME=user-admin-api-secret

# Get secret data from env vars
GOOGLE_OAUTH_CLIENT_ID_BASE64_ENCODED=$(echo "$GOOGLE_OAUTH_CLIENT_ID" | base64)
GOOGLE_OAUTH_CLIENT_SECRET_BASE64_ENCODED=$(echo "$GOOGLE_OAUTH_CLIENT_SECRET" | base64)

# Generate secret YAML 
cat <<EOF > ${SECRET_NAME}.yaml
apiVersion: v1
kind: Secret
metadata:
  name: ${SECRET_NAME}
type: Opaque
data:
  GOOGLE_OAUTH_CLIENT_ID: ${GOOGLE_OAUTH_CLIENT_ID_BASE64_ENCODED}
  GOOGLE_OAUTH_CLIENT_SECRET: ${GOOGLE_OAUTH_CLIENT_SECRET_BASE64_ENCODED}
EOF