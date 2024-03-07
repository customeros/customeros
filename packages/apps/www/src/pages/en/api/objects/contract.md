---
title: Contract Object
description: The contract data object
layout: ../../../../layouts/docs.astro
lang: en
---


## The contract data object

```json
{
  "billingDetails": {
    <billingDetails object>  
  },
  "billingEnabled": true,  
  "committedPeriods": 1,
  "contractEnded": "2024-02-01T20:53:13.381944294Z",
  "contractLineItems": [
    {
      <contractLineItems object>
    }
  ],
  "contractName": "My Contract",
  "contractRenewalCycle": "MONTHLY_RENEWAL",
  "contractSigned": "2024-02-01T20:53:13.381944294Z",
  "contractStatus": "LIVE",
  "contractUrl": "https://acme.salesforce.com/contracts/asdtasd8wetsfset3yasdt34tges",
  "currency": "eur",
  "metadata": {
    <metadata object>
  },
  "opportunities": [
    {
      <opportunities object>
    }
  ],
  "serviceStarted": "2024-01-26T00:00:00Z"
}
  
```

### billingDetails
`Float` representing the total amount due for this invoice, including any applicable taxes.


