---
title: Create Organization Mutation
description: Create Organization Mutation query
layout: ../../../layouts/docs.astro
lang: en
---

## CreateOrganization mutation request

```curl
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-CUSTOMER-OS-API-KEY: <MY_API_KEY_HERE>" \
  -d '{
    "query": "mutation CreateOrganization { 
      organization_Create(input: { 
        customId: "myCustomId", 
        name: "Acme Corp",
        description: "A description of the organization",
        notes: "Notes go here!",
        domains: [
            "acmecorp.com",
            "acmecorp.io"
        ]
        website: "https://acmecorp.io",
        industry: "software",
        subIndustry: "subIndustry",
        industryGroup: "industryGroup",
        public: true,
        isCustomer: true,
        market: "B2B",
        logo: "https://acmecorp.io/logo.png",
        employeeGrowthRate: "0 percent",
        headquarters: "San Francisco",
        yearFounded: 1964,
        employees: 1001,
        slackChannelId: "",
        appSource: "salesforce",
      }) {
        metadata {
          id
        } 
      } 
    }"
  }' 
  https://cos.customeros.ai/query

```

As this is a graphQL request, you are able to specify the exact payload you would like returned in the response. In the example above, we’ve specified that only id is returned, but you can choose from any of the response parameters defined in the [organization object](objects/organizations).

## CreateContractLineItem mutation response
```json
"data": {
    "organization_Create": {
        "metadata": {
            "id": "96d699a8-b986-4dae-9f10-a23196f30c90",
        }
    }
}
```