---
title: Update Organization Mutation
description: Update Organization Mutation query
layout: ../../../layouts/docs.astro
lang: en
---

## UpdateOrganization mutation request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "mutation UpdateOrganization { 
      organization_Update(input: { 
        id: "96d699a8-b986-4dae-9f10-a23196f30c90", 
        customId: "MyNewCustomId",
        patch: true
      }) { 
        id 
      } 
    }"
  }' \
  https://cos.customeros.ai/query

```

The organization mutation request requires that you pass the invoice `id` as a query parameter.  

To update specific fields, include `patch: true` in the query parameters; otherwise, provide the full object.

As this is a graphQL request, you are able to specify the exact payload you would like returned in the response.  In the example above, we've specified that `id` is returned, but you can choose from any of the response parameters defined in the [organization object](objects/organization).

## Organization mutation response
```json
"data": {
    "organization_Update": {
        "id": "96d699a8-b986-4dae-9f10-a23196f30c90",
    }
}
```