---
title: Getting Started
description: CustomerOS API Getting Started guide
layout: ../../../layouts/docs.astro
lang: en
---

## Introduction

The CustomerOS API is a GraphQL API that allows you to interact with your CustomerOS tenant programmatically.

## Authentication

We use a simple API key authentication system. You can request an API key for your tenant from the CustomerOS team.

When making API calls, just include the following headers:

```
-H "X-OPENLINE-TENANT-KEY: "your-api-key" \
-H "Content-Type: application/json" \
```

## Rate Limiting

We rate limit API calls to 1000 per minute. If you exceed this limit, you will receive a 429 response. If you need an increased limit, please contact the CustomerOS team.