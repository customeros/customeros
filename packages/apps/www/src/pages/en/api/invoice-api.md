---
title: Invoice API
description: Invoice data object
layout: ../../../layouts/docs.astro
lang: en
---

## Invoice query request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "query { 
      invoice(id: \"96d699a8-b986-4dae-9f10-a23196f30c90\") { 
        metadata 
        amountDue 
        currency
        due
      }
    }"
  }' \
  https://cos.customeros.ai/query

```

The invoice query request requires that you pass the invoice `id` as a query parameter.  

As this is a graphQL request, you are able to specify the exact payload you would like returned in the response.  In the example above, we've specified that `metadata`, `amountDue`, `currency`, and `due` are returned, but you can choose from any of the response parameters defined in the [invoice object](objects/invoice-object)

## Invoice query response
```json
"data": {
    "query": {
        <invoice object>
    }
}
```