---
title: Contract Query
description: Contract query
layout: ../../../layouts/docs.astro
lang: en
---

## Contract query request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "query { 
      contract(id: \"96d699a8-b986-4dae-9f10-a23196f30c90\") { 
        metadata {
          created
          id
        }
        name
        currency
        status
        owner
      }
    }"
  }' \
  https://cos.customeros.ai/query

```

The contract query request requires that you pass the invoice `id` as a query parameter.  

As this is a graphQL request, you are able to specify the exact payload you would like returned in the response.  In the example above, we've specified that `metadata`, `name`, `currency`, `status`, and `owner` are returned, but you can choose from any of the response parameters defined in the [contract object](objects/contract)

## Invoice query response
```json
"data": {
    "query": {
        <contract object>
    }
}
```