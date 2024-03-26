export const SITE = {
  title: "CustomerOS",
  description: "CustomerOS provides you everything you need to be customer driven.",
  defaultLanguage: "en_US",
};

export const OPEN_GRAPH = {
  image: {
    src: "images/og-image.png",
    alt: "CustomerOS provides you everything you need to be customer driven.",
  },
  twitter: "OpenlineAI",
};

// This is the type of the frontmatter you put in the docs markdown files.
export type Frontmatter = {
  title: string;
  description: string;
  layout: string;
  image?: { src: string; alt: string };
  dir?: "ltr" | "rtl";
  ogLocale?: string;
  lang?: KnownLanguageCode;
  isMdx?: boolean;
};

export const KNOWN_LANGUAGES = {
  // Add more languages here
  // sv: "Svenska",
  // ar: "العربية",
  en: "English",
  // es: "Español",
  // fr: "Français",
  // ja: "日本語",
  // pt: "Português",
  // ru: "Русский",
  // no: "Norsk",
  // pl: "Polski",
  // "zh-hans": "简体中文",
} as const;
export type KnownLanguageCode = keyof typeof KNOWN_LANGUAGES;

export const GITHUB_EDIT_URL = `https://github.com/openline-ai/openline-customer-os/tree/otter/packages/apps/www`;

export const COMMUNITY_INVITE_URL = `https://join.slack.com/t/openline-ai/shared_invite/zt-1i6umaw6c-aaap4VwvGHeoJ1zz~ngCKQ`;

// See "Algolia" section of the README for more information.
export const ALGOLIA = {
  indexName: "openline",
  appId: "ZP9RKP00LB",
  apiKey: "6e9726a151f778e76c05baa495c27895",
};

export type OuterHeaders = "CustomerOS" | "APIs" | "Integrations" | "CLI";

export type SidebarItem<TCode extends KnownLanguageCode = KnownLanguageCode> = {
  text: string;
  link: `/${TCode}/${string}`;
};

export type SidebarItemLink = SidebarItem["link"];

export type Sidebar = {
  [TCode in KnownLanguageCode]: {
    [THeader in OuterHeaders]?: SidebarItem<TCode>[];
  };
};
export const SIDEBAR: Sidebar = {
  // For Translations:
  // Keep the "outer headers" in English so we can match them.
  // Translate the "inner headers" to the language you're translating to.
  // Omit any files you haven't translated, they'll fallback to English.
  // Example:
  // sv: {
  //   "CustomerOS API": [
  //     { text: "Introduktion", link: "sv/introduction" },
  //     { text: "Installation", link: "sv/installation" },
  //   ],
  //   Usage: [{ text: "Miljövariabler", link: "sv/usage/env-variables" }],
  // },
  en: {
    "CustomerOS": [
      { text: "Introduction", link: "en/introduction" },
      { text: "Why CustomerOS?", link: "en/why" },
      { text: "Where to Start?", link: "en/where-to-start" },
      { text: "FAQ", link: "en/faq" },
    ],
    Integrations: [
      { text: "Integrations", link: "en/integrations" },
    ],
    "API Overview": [
      { text: "Getting Started", link: "en/api/getting-started" },
      { text: "Webhooks", link: "en/api/webhooks"},
    ],
    "Contracts API": [
      { text: "Query: Contracts", link: "en/api/contract-query"},
      { text: "Mutation: Update Contracts", link: "en/api/update-contract-mutation"},
      { text: "Mutation: Create ContractLineItem", link: "en/api/create-contract-line-item-mutation"},
      { text: "Mutation: Update ContractLineItem", link: "en/api/update-contract-line-item-mutation"},
      { text: "Contract Events", link: "en/api/contract-events"},
    ],
    "Invoices API": [
      { text: "Query: Invoice", link: "en/api/invoice-query"},
      { text: "Mutation: Update Invoice", link: "en/api/update-invoice-mutation"},
      { text: "Invoice Events", link: "en/api/invoice-events"},
    ],
    "Organizations API": [
      { text: "Query: Organization", link: "en/api/organization-query"},
      { text: "Mutation: Create Organization", link: "en/api/create-organization-mutation"},
      { text: "Mutation: Update Organization", link: "en/api/update-organization-mutation"},
      { text: "Organization Events", link: "en/api/organization-events"},
    ],
    "API Data Objects": [
      { text: "billingDetails", link: "en/api/objects/billing-details"},
      { text: "contract", link: "en/api/objects/contract"},
      { text: "contractLineItems", link: "en/api/objects/contract-line-items"},
      { text: "customFieldTemplate", link: "en/api/objects/custom-field-template"},
      { text: "customFields", link: "en/api/objects/custom-fields"},
      { text: "email", link: "en/api/objects/email"},
      { text: "externalLinks", link: "en/api/objects/external-links"},
      { text: "invoice", link: "en/api/objects/invoice"},
      { text: "invoicelineItems", link: "en/api/objects/invoice-line-items"},
      { text: "jobRole", link: "en/api/objects/job-role"},
      { text: "metadata", link: "en/api/objects/metadata"},
      { text: "opportunities", link: "en/api/objects/opportunities"},
      { text: "organization", link: "en/api/objects/organization"},
      { text: "tax", link: "en/api/objects/invoice-tax"},
    ],
    CLI: [
      { text: "Getting Started", link: "en/cli/getting-started" },
      { text: "Commands", link: "en/cli/commands" },
      { text: "Services", link: "en/cli/services" },
    ],
  }
};

export const SIDEBAR_HEADER_MAP: Record<
  Exclude<KnownLanguageCode, "en">,
  Record<OuterHeaders, string>
> = {
  // Translate the sidebar's "outer headers" here
  // sv: {
  //   "CustomerOS API": "CustomerOS API",
  //   APIs: "APIs",
  // },
};
