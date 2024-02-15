---
title: Invoice Events
description: Events for an invoice
layout: ../../../layouts/docs.astro
lang: en
---

Listen for `invoice` events on your CustomerOS instance so you can take action in your application.

## Invoice events

You can subscribe to any of the invoice events described below.  All events send an event identifier and the [`invoice` data object](objects/invoice-object) to your configured webhook endpoint.

You can use the same webhook URL for as many events as you like.

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

### invoice.status.voided
This event occurs when an invoice has been voided and is no longer able to be acted on.

## Example event

```json
{
  "data": {
    "amountDue": 0.0,
    "amountPaid": 0.0,
    "amountRemaining": 0.0,
    "contract": {
      "contractName": "My Contract",
      "contractStatus": "LIVE",
      "metadata": {
        "id": "96d612a8-b086-4dae-9f10-a12796f30c55"
      }
    },
    "currency": "usd",
    "due": "2024-02-01T19:42:00.656805391Z",
    "invoiceLineItems": [
      {
        "description": "My Amazing Product",
        "metadata": {
          "id": "96d612a8-b086-4dae-9f10-a12796f30c55"
        }
      }
    ],
    "invoiceNumber": "LAM-23423",
    "invoicePeriodEnd": "2024-01-26T00:00:00Z",
    "invoicePeriodStart": "2024-01-26T00:00:00Z",
    "invoiceUrl": "https://customeros.ai/invoices/96d612a8-b086-4dae-9f10-a12796f30c55",
    "metadata": {
      "created": "2024-02-01T19:42:00.656805391Z",
      "id": "96d612a8-b086-4dae-9f10-a12796f30c55"
    },
    "note": "",
    "organization": {
      "customerOsId": "C-XSC-SDF",
      "metadata": {
        "id": "96d612a8-b086-4dae-9f10-a12796f30c55"
      },
      "name": "Acme Corp"
    },
    "paid": false,
    "status": "DUE",
    "subtotal": 0.0,
    "taxDue": 0.0
  },
  "event": "invoice.finalized"
}
```