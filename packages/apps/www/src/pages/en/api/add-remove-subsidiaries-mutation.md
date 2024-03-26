---
title: Add & Remove Organization Subsidiaries Mutation
description: Add & Remove Organization Subsidiaries Mutation query
layout: ../../../layouts/docs.astro
lang: en
---

## AddSubsidiaryToOrganization mutation request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "mutation AddSubsidiaryToOrganization { 
      organization_AddSubsidiary(input: { 
        organizationId: "ed3b1fde-6905-47e3-80fe-8f5327672bb1",
        subsidiaryId: "c7452931-8e2e-4796-a97e-ee75a6f908aa",
        type: "Branch"
      }) {
          <organization object>
        } 
      } 
    }"
  }' 
  https://cos.customeros.ai/query

```

This request will allow you to link 2 organizations a parent/child relationship (typically one organization will be a subsidiary of another). This can be used to handle organizations that have relationships such as head office and branches or a main office with store locations. You can use the freetext ```type``` field to denote the relationship, such as ```store``` or ```branch```.

As this is a graphQL request, you are able to specify the exact payload you would like returned in the response. In the example above, we’ve specified that only id is returned, but you can choose from any of the response parameters defined in the [organization object](objects/organization).

## AddSubsidiaryToOrganization mutation response
```json
"data": {
    "organization_AddSubsidiary": {
        "metadata": {
            "id": "ed3b1fde-6905-47e3-80fe-8f5327672bb1",
        }
    }
}
```

The response will contain the parent organization object or parts of it if specified in the request.

## RemoveSubsidiaryFromOrganization mutation request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "mutation RemoveSubsidiaryFromOrganization { 
      organization_RemoveSubsidiary(input: { 
        organizationId: "ed3b1fde-6905-47e3-80fe-8f5327672bb1",
        subsidiaryId: "c7452931-8e2e-4796-a97e-ee75a6f908aa",
      }) {
          <organization object>
        } 
      } 
    }"
  }' 
  https://cos.customeros.ai/query

```

This request will allow you to unlink a subsidiary Organization from its parent Organization.

## RemoveSubsidiaryFromOrganization mutation response
```json
"data": {
    "organization_RemoveSubsidiary": {
        <organization object>
    }
}
```

The response will contain the parent organization object or parts of it if specified in the request.