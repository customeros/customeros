---
title: Location object
description: The location data object
layout: ../../../../layouts/docs.astro
lang: en
---

## The location data object

```json
{
  "addressLine1": "20 Main St",
  "addressLine2": "Apt 8",
  "addressType": "Commercial",
  "commercial": true,
  "country": "USA",
  "district": "",
  "latitude": "37.7749",
  "locality": "San Francisco", 
  "locationName": "my location",
  "longitude": "-122.4194",
  "metadata": {
    "<metadata object>"
  },
  "plusFour": "0000",
  "postalCode": "91010",
  "predirection": "20",
  "rawAddress": "20 Main St. Apt. 8 San Francisco, CA, 91010",
  "region": "CA",
  "street": "Main St.",  
  "timeZone": "US-Pacific",
  "utcOffset": 8  
}
```

### id
Unique `string` identifying the object.  This property is always set.