---
title: Organization Events
description: Events for an organization
layout: ../../../layouts/docs.astro
lang: en
---

Listen for `organization` events on your CustomerOS instance so you can take action in your application.

## Organization events

You can subscribe to any of the organization events described below.  All events send an event identifier and the [`organization` data object](organization-object) to your configured webhook endpoint.

You can use the same webhook URL for as many events as you like.

### organization.archived

### organization.created

### organization.merged

### organization.updated

### organization.onboarding.done

### organization.onboarding.late

### organization.onboarding.stuck

### organization.relationship.customer

### organization.relationship.prospect

## Example event

```json
{
  "event": "organization.onboarding.late",
  "data": {
    <organization object>
}
```