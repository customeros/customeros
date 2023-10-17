---
title: CustomerOS Customer APIs
description: CustomerOS Customer APIs
layout: ../../../layouts/docs.astro
lang: en
---

## Contacts

### Get all contacts

With email validation

```bash
curl \
-X POST \
-H "X-OPENLINE-TENANT-KEY: your-api-key" \
-H "Content-Type: application/json" \
-d '{"query":"query { contact(id: \"CONTACT-ID-HERE\") { id emails{id emailValidationDetails {validated isReachable isValidSyntax canConnectSmtp acceptsMail hasFullInbox isCatchAll isDeliverable isDisabled}}} "}' \
https://api.customeros.ai/query
```

### Create contact with email validation

```bash
curl \
-X POST \
-H "X-OPENLINE-TENANT-KEY: your-api-key" \
-H "Content-Type: application/json" \
-d '{"query": "mutation { customer_contact_Create(input: {prefix: \"Ms.\", firstName: \"X\", lastName: \"Y\", appSource:\"YOUR-APP\", email: {primary:true, email:\"someone@somedomain.com\", label: WORK, appSource:\"YOUR-APP\"}}) { id, email {id}}}"}' \
https://api.customeros.ai/query
```