---
title: Organization object
description: The organization data object
layout: ../../../../layouts/docs.astro
lang: en
---

## The organization data object

```json
{
  "accountDetails": {
    <accountDetails object>
  },
  "contracts": [
    <contract object>
  ],
  "customFields": [
    <customFields object>
  ],
  "description": "The best company in the world",
  "domains": [
    "acmecorp.com",
    "acmecorp.io" 
  ],
  "employeeGrowthRate": "0 percent",
  "employees": 0,
  "headquarters": "San Francisco",
  "industry": "software",
  "industryGroup": "industryGroup",
  "isCustomer": true,
  "lastFundingAmount": "100000000",
  "lastFundingRound": "SERIES B",
  "lastTouchpoint": {
    "<lastTouchpoint object>"
  },
  "locations": [
    <location object>
  ],
  "logo": "https://acmecorp.io/logo.png",
  "market": "B2B",  
  "metadata": {
    <metadata object>
  },
  "name": "Acme Corp",
  "notes": "",
  "owner": "sara@acmecorp.io",
  "parentCompany": [
    <parentCompany object>
  ],
  "people": [
    <people object>
  ],
  "socialMedia": [
    <socialMedia object>
  ],
  "subindustry": "subindustry",
  "subsidiaries": [
    <subsidiaries object>
  ],
  "tags": [
    <tag object>
  ],
  "targetAudience": "Developers and Product Managers",
  "timelineEvents": [
    <timelineEvents object>
  ],
  "valueProposition": "To enable you to win",
  "website": "https://acmecorp.io",
  "yearFounded": 2000
}
```

### id
Unique `string` identifying the object.  This property is always set.