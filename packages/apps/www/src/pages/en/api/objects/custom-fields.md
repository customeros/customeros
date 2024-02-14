---
title: Custom fields object
description: The custom fields data object
layout: ../../../../layouts/docs.astro
lang: en
---

## The customFields data object

```json
{
  "dataType": "TEXT",
  "fieldName": "my custom field",
  "metadata": {
    <metadata object>
  },
  "template": {
    <customFieldTemplate object>
  },
  "value": "my custom value"
}
```

### id
Unique `string` identifying the object.  This property is always set.