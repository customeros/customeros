---
title: Update Contract Line Item Mutation
description: Update contract mutation
layout: ../../../layouts/docs.astro
lang: en
---

## UpdateContractLineItem mutation request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "mutation UpdateContractLineItem { 
      contractLineItem_Update(input: { 
        id: \"96d699a8-b986-4dae-9f10-a23196f30c90\", 
        description: "My Fantastic Product",
        quantity: 10,
        price: 100.00,
        patch: true 
      }) { 
        id
      } 
    }"
  }' 
  https://cos.customeros.ai/query

```

The contractLineItem mutation request requires that you pass the contract `id` as a query parameter.  

In order to update only the fields specified in the request, you must pass `patch: true` as part of the query parameters.  If you do not, you must pass the full object in the request.

As this is a graphQL request, you are able to specify the exact payload you would like returned in the response.  In the example above, we've specified that only `id` is returned, but you can choose from any of the response parameters defined in the [contractLineItem object](objects/contract-line-items)

## UpdateContract mutation response
```json
"data": {
    "contractLineItem_Update": {
        "id": "96d699a8-b986-4dae-9f10-a23196f30c90",
    }
}
```