<div align="center">
  <a href="https://customeros.ai">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/.github/CustomerOS-logo-darkmode-small.svg">
      <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/.github/CustomerOS-logo-small.svg">
      <img alt="CustomerOS Logo" src="https://raw.githubusercontent.com/openline-ai/openline-customer-os/otter/.github/CustomerOS-logo-small.svg">
      </picture>
  </a>
  <br />
  <p>
    <h3>
      <b>
        CustomerOS
      </b>
    </h3>
  </p>
  <p>
    CustomerOS is the easiest way to consolidate, warehouse, and build applications with your customer data.
  </p>
  <p>

[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen?logo=github)][customerOS-repo] 
[![license](https://img.shields.io/badge/license-Apache%202-blue)][apache2] 
[![stars](https://img.shields.io/github/stars/openline-ai/openline-customer-os?style=social)][customerOS-repo] 
[![twitter](https://img.shields.io/twitter/follow/openlineAI?style=social)][twitter] 
[![slack](https://img.shields.io/badge/slack-community-blueviolet.svg?logo=slack)][slack]

  </p>
  <p>
    <sub>
      Built with ❤︎ by the
      <a href="https://customeros.ai">
        CustomerOS
      </a>
      community!
    </sub>
  </p>
</div>


## 🚀 Installation

1. Download and install the [CustomerOS CLI][cli]
2. Run the following command

```sh-session
openline dev start customer-os
```

## 🤝 Resources

- Our [docs site][docs] has numerous guides and reference material for to make building on customerOS easy.
- For help, feature requests, or chat with fellow Openline enthusiasts, check out our [slack community][slack]!

## Codebase

### Technologies

Here's a list of the big technologies that we use:

- **PostgreSQL** & **Neo4j** - Data storage
- **Go** - Back end & API
- **TypeScript** - Web components
- **React** - Front end apps and UI components

### Folder structure

```sh
openline-customer-os/
├── architecture            # Architectural documentation
├── deployment              
│   ├── infra               # Infrastructure-as-code
│   └── scripts             # Deployment scripts
└── packages
    ├── apps                # Front end web applications
    │   ├── launcher        # customerOS app launcher & home screen
    │   └── settings        # customerOS system settings & app configuration
    ├── auth                # Authentication
    ├── components
    │   ├── react           # React component library
    │   └── web             # Web & UI component library
    ├── core                # Shared core libraries
    └── server              # Back end database & API server
```

## 💪 Contributions

- We love contributions big or small!  Please check out our [guide on how to get started][contributions].
- Not sure where to start?  [Book a free, no-pressure, no-commitment call][call] with the team to discuss the best way to get involved.

## ✨ Contributors

A massive thank you goes out to all these wonderful people ([emoji key][emoji]):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/mattbr0wn"><img src="https://avatars.githubusercontent.com/u/113338429?v=4?s=100" width="100px;" alt="Matt Brown"/><br /><sub><b>Matt Brown</b></sub></a><br /><a href="https://github.com/openline-ai/openline-customer-os/commits?author=mattbr0wn" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://openline.ai"><img src="https://avatars.githubusercontent.com/u/88987042?v=4?s=100" width="100px;" alt="Vasi Coscotin"/><br /><sub><b>Vasi Coscotin</b></sub></a><br /><a href="https://github.com/openline-ai/openline-customer-os/commits?author=xvasi" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/alexopenline"><img src="https://avatars.githubusercontent.com/u/95470380?v=4?s=100" width="100px;" alt="alexopenline"/><br /><sub><b>alexopenline</b></sub></a><br /><a href="https://github.com/openline-ai/openline-customer-os/commits?author=alexopenline" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/edifirut"><img src="https://avatars.githubusercontent.com/u/108661145?v=4?s=100" width="100px;" alt="edifirut"/><br /><sub><b>edifirut</b></sub></a><br /><a href="#infra-edifirut" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/openline-ai/openline-customer-os/pulls?q=is%3Apr+reviewed-by%3Aedifirut" title="Reviewed Pull Requests">👀</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/jontyk"><img src="https://avatars.githubusercontent.com/u/81759836?v=4?s=100" width="100px;" alt="Jonty Knox"/><br /><sub><b>Jonty Knox</b></sub></a><br /><a href="https://github.com/openline-ai/openline-customer-os/pulls?q=is%3Apr+reviewed-by%3Ajontyk" title="Reviewed Pull Requests">👀</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/tsearle"><img src="https://avatars.githubusercontent.com/u/4540323?v=4?s=100" width="100px;" alt="tsearle"/><br /><sub><b>tsearle</b></sub></a><br /><a href="https://github.com/openline-ai/openline-customer-os/commits?author=tsearle" title="Code">💻</a> <a href="https://github.com/openline-ai/openline-customer-os/commits?author=tsearle" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://www.openline.ai/"><img src="https://avatars.githubusercontent.com/u/135133139?v=4?s=100" width="100px;" alt="Alex Calinica"/><br /><sub><b>Alex Calinica</b></sub></a><br /><a href="https://github.com/openline-ai/openline-customer-os/pulls?q=is%3Apr+reviewed-by%3Aacalinica" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/openline-ai/openline-customer-os/commits?author=acalinica" title="Tests">⚠️</a> <a href="https://github.com/openline-ai/openline-customer-os/commits?author=acalinica" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/SilviuOpenline"><img src="https://avatars.githubusercontent.com/u/125381623?v=4?s=100" width="100px;" alt="SilviuOpenline"/><br /><sub><b>SilviuOpenline</b></sub></a><br /><a href="https://github.com/openline-ai/openline-customer-os/pulls?q=is%3Apr+reviewed-by%3ASilviuOpenline" title="Reviewed Pull Requests">👀</a> <a href="https://github.com/openline-ai/openline-customer-os/commits?author=SilviuOpenline" title="Tests">⚠️</a> <a href="https://github.com/openline-ai/openline-customer-os/commits?author=SilviuOpenline" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## 🪪 License

- This repo is licensed under [Apache 2.0][apache2], with the exception of the ee directory (if applicable).
- Premium features (contained in the ee directory) require an Openline Enterprise license.  See our [pricing page][pricing] for more details.
- Copyright &copy; Openline Technologies Inc. 2022 - 2024

<!--- References --->

[apache2]: https://www.apache.org/licenses/LICENSE-2.0
[call]: https://meetings-eu1.hubspot.com/matt2/customer-demos
[cli]: https://docs.customeros.ai/en/cli/getting-started
[contributions]: https://github.com/openline-ai/community/blob/main/README.md
[customerOS-repo]: https://github.com/openline-ai/openline-customer-os/
[docs]: https://docs.customeros.ai/
[emoji]: https://allcontributors.org/docs/en/emoji-key
[oasis]: https://github.com/openline-ai/openline-oasis
[pricing]: https://www.customeros.ai/pricing
[slack]: https://join.slack.com/t/openline-ai/shared_invite/zt-1i6umaw6c-aaap4VwvGHeoJ1zz~ngCKQ
[twitter]: https://twitter.com/OpenlineAI
