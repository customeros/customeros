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
  "customerOsId": "C-XSC-SDF",
  "customFields": [
    <customFields object>
  ],
  "customId": "myCustomId",
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
  "lastFundingRound": SERIES_B,
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
  "owner": [
    <user object>
  ],
  "parentCompanies": [
    <parentCompany object>
  ],
  "people": [
    <people object>
  ],
  "public": true,
  "socialMedia": [
    <socialMedia object>
  ],
  "subIndustry": "subIndustry",
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

### lastFundingRound
An enum representing the last funding round for organization. Valid values are:
- `PRE_SEED`
- `SEED`
- `SERIES_A`
- `SERIES_B`
- `SERIES_C`
- `SERIES_D`
- `SERIES_E`
- `SERIES_F`
- `IPO`
- `FRIENDS_AND_FAMILY`
- `ANGEL`
- `BRIDGE`

### market
An enum. Valid values are:
- `B2B`
- `B2C`
- `MARKETPLACE`