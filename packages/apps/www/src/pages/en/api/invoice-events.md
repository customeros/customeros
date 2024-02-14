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
  "event": "invoice.finalized",
  "data": {
    <invoice object>
}
```