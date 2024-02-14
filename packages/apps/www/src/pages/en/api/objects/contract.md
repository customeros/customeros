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
  "committedPeriods": 1,
  "contractEnded": "2024-02-01T20:53:13.381944294Z",
  "contractLineItems": [
    {
      <contractLineItems object>
    }
  ],
  "contractRenewalCycle": "MONTHLY_RENEWAL",
  "contractSigned": "2024-02-01T20:53:13.381944294Z",
  "contractUrl": "https://acme.salesforce.com/contracts/asdtasd8wetsfset3yasdt34tges",
  "currency": "eur",
  "invoicingEnabled": true,  
  "metadata": {
    <metadata object>
  },
  "contractName": "My Contract",
  "opportunities": [
    {
      <opportunities object>
    }
  ],
  "serviceStarted": "2024-01-26T00:00:00Z",
  "contractStatus": "LIVE"
}
  
```

### amountDue
`Float` representing the total amount due for this invoice, including any applicable taxes.

### amountPaid
`Float` representing the amount that has been paid against the invoice.

### amountRemaining
`Float` representing the difference between the `amountDue` and `amountPaid`.

### contract
The [`contract` object](contract-object)

### currency
`Enum` representing the three-letter ISO currency code, in lower case.

### due
ISO 8601 `timestamp` representing the date payment for this invoice is due.

### invoiceNumber
A unique, identifying `string` that appears on the invoice.

### invoicePeriodStart
ISO 8601 `timestamp` of the first day in the invoice period.

### invoicePeriodEnd
ISO 8601 `timestamp` of the last day in the invoice period.

### invoiceUrl
A nullable `string` representing the URL for the hosted invoice page which allows customers to view and pay an invoice.

### lineItems
The invoice [`line items` object](invoice-line-items-object)

### metadata
The [`metadata` object](metadata-object)

### note
A `string` representing any notes on the invoice.  Can be an empty string.

### organization
The [`organization` object](organization-object)

### paid
A `boolean` indicating if the invoice has been paid in full.

### subtotal
A `float` representing the sum of all line items minus tax/VAT.

### status
An `enum` representing the current status of the invoice.  Valid values are:
- DRAFT
- DUE
- PAID
- OVERDUE
- UNCOLLECTABLE
- VOIDED

### tax
The invoice [`tax` object](invoice-tax-object)

