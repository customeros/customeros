---
title: Contracts Events
description: Events for a contract
layout: ../../../layouts/docs.astro
lang: en
---

Listen for `contract` events on your CustomerOS instance so you can take action in your application.

## Contract events

You can subscribe to any of the contract events described below.  All events send an event identifier and the [`contract` data object](contract-object) to your configured webhook endpoint.

You can use the same webhook URL for as many events as you like.

### contract.created
This event occurs when a new contract is created on an organization.

### contract.renewed
This event occurs when a contract is renewed.

### contract.ended
This event occurs when a contract is ended.

## Example event

```json
{
  "event": "contract.renewed",
  "data": {
    <contract object>
}
```