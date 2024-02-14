---
title: Billing details object
description: The billing details data object
layout: ../../../../layouts/docs.astro
lang: en
---

## The billing details data object

```json
"billingDetails": {
        "billingCycle": "MONTHLY_BILLING",
        "invoicingStarted": "2024-01-26T00:00:00Z",
        "addressLine1": "20 Main St",
        "addressLine2": "Apt 8",
        "locality": "San Francisco",
        "region": "CA",
        "country": "USA",
        "postalCode": "91010",
        "organizationLegalName": "Acme Corp",
        "billingEmail": "finance@acmecorp.com",
        "invoiceNote": "",
    }
```

### id
Unique `string` identifying the object.  This property is always set.

### created
ISO 8601 `timestamp` of when the object was created.

### lastUpdated
ISO 8601 `timestamp` of when the object was last updated.