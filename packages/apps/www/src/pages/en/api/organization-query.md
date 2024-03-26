---
title: Organization Query
description: Organization query
layout: ../../../layouts/docs.astro
lang: en
---

## Organization query request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "query {
      organization(id: \"96d699a8-b986-4dae-9f10-a23196f30c90\") {
        <organization-object>
        name
        description
      }
    }"
  }' \
  https://cos.customeros.ai/query

```

## Organization query by custom ID request

You are also able to query organizations by specific fields. For instance, when querying by your own identifier that you have updated in the ```customId``` field, you can use the following query:

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "query {
      organization_ByCustomId(customId: \"96d699a8-b986-4dae-9f10-a23196f30c90\") {
        <organization-object>
      }
    }"
  }' \
  https://cos.customeros.ai/query

```

The organization query by Custom ID request requires that you pass the organization `customId` as a query parameter.  

As this is a graphQL request, you are able to specify the exact payload you would like returned in the response.  In the first example above, we've specified that `name`, and `description` are returned, but you can choose from any of the response parameters defined in the [organization object](objects/organization)

## Invoice query response
```json
"data": {
    "query": {
        <organization object>
    }
}
```



