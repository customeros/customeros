---
title: Invoice line items object
description: The Invoice line items data object
layout: ../../../../layouts/docs.astro
lang: en
---

## The invoiceLineItems data object

```json
[
    {
        "metadata": {
            <metadata object>
        },
        "description": "My Amazing Product",
        "price": 0.0,
        "quantity": 0,
        "subtotal": 0.0,
        "tax": {
            <tax object>
        },
        "total": 0.0
    }
]
```

### id
Unique `string` identifying the object.  This property is always set.