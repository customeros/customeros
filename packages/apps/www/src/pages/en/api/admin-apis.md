---
title: CustomerOS Admin APIs
description: CustomerOS Admin APIs
layout: ../../../layouts/docs.astro
lang: en
---

## Admin

### Get all users

With information stored against that user

```bash
curl \
-X POST \
-H "X-OPENLINE-TENANT-KEY: your-api-key" \
-H "Content-Type: application/json" \
-d '{"query": "query { users { totalPages totalElements content { id firstName lastName calendars { link } jobRoles{ jobTitle description } emails { email } createdAt }}} "}' https://api.customeros.ai/query