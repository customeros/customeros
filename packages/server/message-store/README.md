<div align="center">
  <a href="https://openline.ai">
    <img
      src="https://www.openline.ai/TeamHero.svg"
      alt="Openline Logo"
      height="64"
    />
  </a>
  <br />
  <p>
    <h3>
      <b>
        Openline Oasis app
      </b>
    </h3>
  </p>
  <p>
    Openline customerOS is the easiest way to consolidate, warehouse, and build applications with your customer data.
  </p>
  <p>

[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen?logo=github)][oasis-repo] 
[![license](https://img.shields.io/badge/license-Apache%202-blue)][apache2] 
[![stars](https://img.shields.io/github/stars/openline-ai/openline-customer-os?style=social)][oasis-repo] 
[![twitter](https://img.shields.io/twitter/follow/openlineAI?style=social)][twitter] 
[![slack](https://img.shields.io/badge/slack-community-blueviolet.svg?logo=slack)][slack]

  </p>
  <p>
    <sub>
      Built with ❤︎ by the
      <a href="https://openline.ai">
        Openline
      </a>
      community!
    </sub>
  </p>
</div>


## 👋 Overview

Message store

## 🚀 Installation

### generate
https://grpc.io/docs/protoc-installation/
https://grpc.io/docs/languages/go/quickstart/

To build the project run

```
make
```
## Development

To run this service to run on your laptop you need the following environemnt vars to be set

| Variable                  | Meaning                                                        |
|---------------------------|----------------------------------------------------------------|
| DB_HOST                   | hostname of postgres db                                        |
| DB_PORT                   | port of postgres db                                            |
| DB_NAME                   | database name                                                  |
| DB_USER                   | user to log into db as                                         |
| DB_PASS                   | the database password                                          |
| MESSAGE_STORE_SERVER_PORT | port to bind for GRPC communications should be set to "9009"   |
| CUSTOMER_OS_API           | URL to access the Customer OS API                              |


## 🙌 Features

TBD

## 🤝 Resources

- For help, feature requests, or chat with fellow Openline enthusiasts, check out our [slack community][slack]!
- Our [docs site][docs] has references for developer functionality, including the Graph API

## 💪 Contributions

- We love contributions big or small!  Please check out our [guide on how to get started][contributions].
- Not sure where to start?  [Book a free, no-pressure, no-commitment call][call] with the team to discuss the best way to get involved.

## ✨ Contributors

A massive thank you goes out to all these wonderful people ([emoji key][emoji]):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->


<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## 🪪 License

- This repo is licensed under [Apache 2.0][apache2], with the exception of the ee directory (if applicable).
- Premium features (contained in the ee directory) require an Openline Enterprise license.  See our [pricing page][pricing] for more details.


[apache2]: https://www.apache.org/licenses/LICENSE-2.0
[call]: https://meetings-eu1.hubspot.com/matt2/customer-demos
[careers]: https://openline.ai
[contributions]: https://github.com/openline-ai/community/blob/main/README.md
[docs]: https://openline.ai
[emoji]: https://allcontributors.org/docs/en/emoji-key
[oasis-repo]: https://github.com/openline-ai/openline-customer-os/
[pricing]: https://openline.ai/pricing
[slack]: https://join.slack.com/t/openline-ai/shared_invite/zt-1i6umaw6c-aaap4VwvGHeoJ1zz~ngCKQ
[twitter]: https://twitter.com/OpenlineAI
