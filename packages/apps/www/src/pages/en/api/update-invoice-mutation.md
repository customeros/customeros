---
title: Invoice Mutation
description: Invoice mutation
layout: ../../../layouts/docs.astro
lang: en
---

## UpdateInvoice mutation request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "mutation UpdateInvoice { 
      invoice_Update(input: { 
        id: \"96d699a8-b986-4dae-9f10-a23196f30c90\", 
        paid: true, 
        status: \"PAID\" 
        patch: true
      }) { 
        id 
        paid 
        status 
      } 
    }"
  }' \
  https://cos.customeros.ai/query

```

The invoice mutation request requires that you pass the invoice `id` as a query parameter.  

In order to update only the fields specified in the request, you must pass `patch: true` as part of the query parameters.  If you do not, you must pass the full object in the request.

As this is a graphQL request, you are able to specify the exact payload you would like returned in the response.  In the example above, we've specified that `id`, `paid`, and `status`are returned, but you can choose from any of the response parameters defined in the [invoice object](objects/invoice)

## Invoice mutation response
```json
"data": {
    "invoice_Update": {
        "id": "96d699a8-b986-4dae-9f10-a23196f30c90",
        "paid": true,
        "status": "PAID"
    }
}
```