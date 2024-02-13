---
title: Invoice Object
description: The invoice data object
layout: ../../../layouts/docs.astro
lang: en
---


## The invoice data object

```json
{
  "id": "96d612a8-b086-4dae-9f10-a12796f30c55",
  "object": "invoice",
  "amountDue": 0.0,
  "amountPaid": 0.0,
  "amountRemaining": 0.0,
  "created": "2024-02-01T19:42:00.656805391Z",
  "currency": "usd",
  "due": "2024-02-01T19:42:00.656805391Z",
  "invoiceUrl": "https://customeros.ai/invoices/96d612a8-b086-4dae-9f10-a12796f30c55",
  "lastUpdated": "2024-02-01T20:53:13.381944294Z",
  "lineItems": [
    {
      "id": "6d235a8-b086-4dae-9f10-a12796f30c55",
      "amountDue": 0.0,
      "created": "2024-02-01T19:42:00.656805391Z",
      "description": "My Subscription Plan",
      "price": 0.0,
      "quantity": 0,
      "tax": {
        "salesTax": 0.0,
        "vat": 0.0  
      }
    }
  ],
  "note": null,
  "number": "LAM-23423",
  "organization": {
    "id": "96d612a8-b086-4dae-9f10-a12796f30c55",
    "customerOsId": "C-SDF-WER",  
    "name": "Acme Corp"
  },
  "paid": false,
  "subtotal": 0.0,
  "status": "OPEN",
  "tax": {
    "salesTax": 0.0,
    "vat": 200.0
  }
}
```


### id
Unique `string` idenfying the invoice object.  This property is always set.

### object
`String` representing the object's type.  Objects of the same type share the same value. 

### amountDue
`Float` representing the total amount due for this invoice, including any applicable taxes.

### amountPaid
`Float` representing the amount that has been paid against the invoice.

### amountRemaining
`Float` representing the difference between the `amountDue` and `amountPaid`.

### created
ISO 8601 `timestamp` of when the invoice object was created.

### currency
`Enum` representing the three-letter ISO currency code, in lower case.

### due
ISO 8601 `timestamp` representing the date payment for this invoice is due.

### invoiceUrl
A nullable `string` representing the URL for the hosted invoice page which allows customers to view and pay an invoice.

### lastUpdated
ISO 8601 `timestamp` of when the invoice was last updated.

### lineItems


### note

### number

### organization

### paid

### subtotal

### status

### tax