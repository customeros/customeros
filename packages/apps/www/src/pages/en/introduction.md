---
title: Introduction
description: Introduction to CustomerOS
layout: ../../layouts/docs.astro
lang: en
---

## CustomerOS, your customer data mesh

CustomerOS was borne from the need to have a single source of truth for all customer data in order to drive customer-centric workflows. It is essentially a data mesh that has a number of different data sources that are all mapped into the mesh, and then made accessible to any application via the CustomerOS API suite.

### Wait, what is a data mesh exactly?

A data mesh is a data architecture that is designed to allow for the easy integration of data from multiple sources. It is a way to organize data in a way that allows for the easy integration of data from multiple sources.

<div class="embed">
<iframe width="560" height="315" src="https://www.youtube.com/embed/zfFyE3xmJ7I" title="Data Mesh 101: What is Data Mesh?" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>

Some great reading about this can be found by Zhamak Dehghani [here](https://martinfowler.com/articles/data-mesh-principles.html) and from the team and INNOQ [here](https://www.datamesh-architecture.com/).

### But aren't you just describing a CDP?

Well, yes. But every CDP that we have worked with is less of a Customer Data **Platform**, and more of a Customer Data **Pipeline**. They are designed to take data from a number of different sources, and then push that data to a number of different destinations. They are not designed to be a single source of truth for all customer data, available in real-time and accessible via an API.

In short, we are similar to a CDP, but more in the truest sense of the word compared to the current available options on the market.

## What kind of data sources do you support?

CustomerOS supports a number of different data sources, including, but not limited to:

- [Hubspot](https://hubspot.com)
- [Salesforce](https://salesforce.com)
- [Zendesk](https://zendesk.com)
- [Intercom](https://intercom.com)
- [Stripe](https://stripe.com)
- [Sendgrid](https://sendgrid.com)
- [Google Drive](https://drive.google.com)
- [Google Calendar](https://calendar.google.com)
- [Google Meet](https://meet.google.com)

All of our integrations can be found [here](/en/integrations)

## What kind of data can you pull from these sources?

We can pull any data that is available via the API of the data source. For example, we can pull the following data from Hubspot:

- Contacts
- Companies
- Deals
- Tickets

## Can you push data to these sources as destinations and keep it all in sync for me?

Right now we don't support pushing data to any of our integrations, but we look at adding this functionality in the future.

## So... sum it up for me, what is CustomerOS? A workflow tool?

CustomerOS is a data mesh first and foremost. However, it's design allows for workflow tooling to be built on top of it. We are currently working on building out a number of different workflow tools, starting with Customer Success tooling. However we have exposed both the source code of the CustomerOS data mesh, and our first workflow tool, so that you can enhance the data mesh, integrate your own data to the mesh, or even build your own tools on top of it.

If you want to build on top of our data mesh using our APIs, start [here](/en/api).

If you want to add data sources or fork CustomerOS to suit your own needs, start [here](/en/cli).