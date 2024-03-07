---
title: Contract line items object
description: The contract line items data object
layout: ../../../../layouts/docs.astro
lang: en
---

## The contractLineItems data object

```json
{
   "billingCycle": "MONTHLY",
   "comments": "",
   "description": "My Amazing Product",
   "metadata": {
      <metadata object>
   },
   "parentId": "96d612a8-b086-4dae-9f10-a12796f30c55",
   "price": 0.0,
   "quantity": 0,
   "serviceEnded": "2024-01-26T00:00:00Z",
   "serviceStarted": "2024-01-26T00:00:00Z",
   "tax": {
      <tax object>
   }
}
```

### id
Unique `string` identifying the object.  This property is always set.

### created
ISO 8601 `timestamp` of when the object was created.

### lastUpdated
ISO 8601 `timestamp` of when the object was last updated.