---
title: Invoice Events
description: Events for an invoice
layout: ../../../layouts/docs.astro
lang: en
---

Listen for events on your CustomerOS account so you can take action in your application.

## Invoice events

You can subscribe to any of the invoice events described below.  All events send the [`invoice` data object](invoice-object) to your configured webhook endpoint.

### invoice.finalized
This event occurs when an open invoice is generated for a customer. 

### invoice.sent
This event occurs when an invoice email is sent to the customer.

### invoice.status.paid
This event occurs when an invoice has been paid in full.

### invoice.status.uncollectible
This event occurs when an invoice has been marked as uncollectible.

### invoice.status.overdue
This event occurs when an invoice is 30 days overdue.

### invoice.voided
This event occurs when an invoice has been voided and is no longer able to be acted on.

## The invoice object

```json
{
  "triggerEvent": "invoice.finalized",
  "data": {
    "id": "96d612a8-b086-4dae-9f10-a12796f30c55",
    "object": "invoice",
    "amountDue": 0.0,
    "amountPaid": 0.0,
    "amountRemaining": 0.0,
    "contract": {
      "id": "96d612a8-b086-4dae-9f10-a12796f30c55",
      "invoicingStarted": "2024-01-26T00:00:00Z",
      "name": "My Contract"
    },
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
    "note": "",
    "number": "LAM-23423",
    "organization": {
      "id": "96d612a8-b086-4dae-9f10-a12796f30c55",
      "customerOsId": "C-SDF-WER",  
      "name": "Acme Corp"
    },
    "paid": false,
    "periodStart": "2024-01-26T00:00:00Z",
    "periodEnd": "2024-01-26T00:00:00Z", 
    "subtotal": 0.0,
    "status": "OPEN", 
    "tax": {
      "salesTax": 0.0,
      "vat": 200.0
    }
  }
}
```
