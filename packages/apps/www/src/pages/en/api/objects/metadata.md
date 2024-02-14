---
title: Metadata Object
description: The metadata data object
layout: ../../../../layouts/docs.astro
lang: en
---

## The metadata data object

```json
{
    "id": "96d612a8-b086-4dae-9f10-a12796f30c55",
    "created": "2024-02-01T19:42:00.656805391Z",
    "customerOsId": "C-XSC-SDF",
    "lastUpdated": "2024-02-01T20:53:13.381944294Z",
    "createdBy": "bob@acme.com",
    "source": "CustomerOS",
    "sourceOfTruth": "Hubspot",
    "appSource": "xyz",
    "externalLinks": [
        {
            <externalLinks object>
        }
    ],
    "referenceId": "myreferenceId",
    "public": true
}
```

### id
Unique `string` identifying the object.  This property is always set.

### created
ISO 8601 `timestamp` of when the object was created.

### lastUpdated
ISO 8601 `timestamp` of when the object was last updated.