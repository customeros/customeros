interface Field {
  name: string;
  label: string;
  textarea?: boolean;
}

export interface IntegrationItem {
  key: string;
  name: string;
  icon: string;
  fields: Field[];
  identifier: string;
  state: 'INACTIVE' | 'ACTIVE';
  isFromIntegrationApp?: boolean;
}

export const integrationsData: IntegrationItem[] = [
  {
    key: 'gsuite',
    state: 'INACTIVE',
    icon: '/integrationOptionLogos/google-icon.svg',

    identifier: 'gsuite',
    name: 'G Suite',
    fields: [
      {
        name: 'privateKey',
        label: 'Private key',
        textarea: true,
      },
      {
        name: 'clientEmail',
        label: 'Service account email',
      },
    ],
  },
  {
    key: 'hubspot',
    state: 'INACTIVE',
    icon: '/integrationOptionLogos/hubspot.svg',
    identifier: 'hubspot',
    name: 'Hubspot',
    fields: [
      {
        name: 'privateAppKey',
        label: 'API key',
      },
    ],
    isFromIntegrationApp: true,
  },
  {
    key: 'smartsheet',
    state: 'INACTIVE',
    icon: '/integrationOptionLogos/smartsheet.svg',
    identifier: 'smartsheet',
    name: 'Smartsheet',
    fields: [
      {
        name: 'id',
        label: 'ID',
      },
      {
        name: 'accessToken',
        label: 'API key',
      },
    ],
  },
  {
    key: 'jira',
    state: 'INACTIVE',
    icon: '/integrationOptionLogos/jira.svg',
    identifier: 'jira',
    name: 'Jira',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
      {
        name: 'domain',
        label: 'Domain',
      },
      {
        name: 'email',
        label: 'Email',
      },
    ],
  },
  {
    key: 'trello',
    state: 'INACTIVE',
    identifier: 'trello',
    name: 'Trello',
    icon: '/integrationOptionLogos/trello.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
      {
        name: 'apiKey',
        label: 'API key',
      },
    ],
  },
  {
    key: 'aha',
    state: 'INACTIVE',
    identifier: 'aha',
    name: 'Aha',
    icon: 'integrationOptionLogos/aha.svg',

    fields: [
      {
        name: 'apiUrl',
        label: 'API Url',
      },
      {
        name: 'apiKey',
        label: 'API key',
      },
    ],
  },
  {
    key: 'airtable',
    state: 'INACTIVE',
    identifier: 'airtable',
    name: 'Airtable',
    icon: 'integrationOptionLogos/airtable.svg',
    fields: [
      {
        name: 'personalAccessToken',
        label: 'Personal access token',
      },
    ],
  },
  {
    key: 'amplitude',
    state: 'INACTIVE',
    identifier: 'amplitude',
    name: 'Amplitude',
    icon: 'integrationOptionLogos/amplitude.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API key',
      },
      {
        name: 'secretKey',
        label: 'Secret key',
      },
    ],
  },
  {
    key: 'asana',
    state: 'INACTIVE',
    identifier: 'asana',
    name: 'Asana',
    icon: 'integrationOptionLogos/asana.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'baton',
    state: 'INACTIVE',
    identifier: 'baton',
    name: 'Baton',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiKey',
        label: 'API key',
      },
    ],
  },
  {
    key: 'babelforce',
    state: 'INACTIVE',
    identifier: 'babelforce',
    name: 'Babelforce',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'regionEnvironment',
        label: 'Region / Environment',
      },
      {
        name: 'accessKeyId',
        label: 'Access Key Id',
      },
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'bigquery',
    state: 'INACTIVE',
    identifier: 'bigquery',
    name: 'BigQuery',
    icon: 'integrationOptionLogos/bigquery.svg',
    fields: [
      {
        name: 'serviceAccountKey',
        label: 'Service account key',
      },
    ],
  },

  {
    key: 'braintree',
    state: 'INACTIVE',
    identifier: 'braintree',
    name: 'Braintree',
    icon: 'integrationOptionLogos/braintree.svg',
    fields: [
      {
        name: 'publicKey',
        label: 'Public Key',
      },
      {
        name: 'privateKey',
        label: 'Private Key',
      },
      {
        name: 'environment',
        label: 'Environment',
      },
      {
        name: 'merchantId',
        label: 'Merchant Id',
      },
    ],
  },
  {
    key: 'callrail',
    state: 'INACTIVE',
    identifier: 'callrail',
    name: 'CallRail',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'account',
        label: 'Account',
      },
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'chargebee',
    state: 'INACTIVE',
    identifier: 'chargebee',
    name: 'Chargebee',
    icon: 'integrationOptionLogos/chargebee.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'productCatalog',
        label: 'Product Catalog',
      },
    ],
  },
  {
    key: 'chargify',
    state: 'INACTIVE',
    identifier: 'chargify',
    name: 'Chargify',
    icon: 'integrationOptionLogos/chargify.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'domain',
        label: 'Domain',
      },
    ],
  },
  {
    key: 'clickup',
    state: 'INACTIVE',
    identifier: 'clickup',
    name: 'ClickUp',
    icon: 'integrationOptionLogos/clickup.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'closecom',
    state: 'INACTIVE',
    identifier: 'closecom',
    name: 'Close.com',
    icon: 'integrationOptionLogos/close.com.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },

  {
    key: 'coda',
    state: 'INACTIVE',
    identifier: 'coda',
    name: 'Coda',
    icon: 'integrationOptionLogos/coda.svg',
    fields: [
      {
        name: 'authToken',
        label: 'Auth Token',
      },
      {
        name: 'documentId',
        label: 'Document Id',
      },
    ],
  },
  {
    key: 'confluence',
    state: 'INACTIVE',
    identifier: 'confluence',
    name: 'Confluence',
    icon: 'integrationOptionLogos/confluence.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
      {
        name: 'domain',
        label: 'Domain',
      },
      {
        name: 'loginEmail',
        label: 'Login Email',
      },
    ],
  },
  {
    key: 'courier',
    state: 'INACTIVE',
    identifier: 'courier',
    name: 'Courier',
    icon: 'integrationOptionLogos/courier.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'customerio',
    state: 'INACTIVE',
    identifier: 'customerio',
    name: 'Customer.io',
    icon: 'integrationOptionLogos/customer-io.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'datadog',
    state: 'INACTIVE',
    identifier: 'datadog',
    name: 'Datadog',
    icon: 'integrationOptionLogos/datadog.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'applicationKey',
        label: 'Application Key',
      },
    ],
  },
  {
    key: 'delighted',
    state: 'INACTIVE',
    identifier: 'delighted',
    name: 'Delighted',
    icon: 'integrationOptionLogos/delighted.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'dixa',
    state: 'INACTIVE',
    identifier: 'dixa',
    name: 'Dixa',
    icon: 'integrationOptionLogos/dixa.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'drift',
    state: 'INACTIVE',
    identifier: 'drift',
    name: 'Drift',
    icon: 'integrationOptionLogos/drift.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'emailoctopus',
    state: 'INACTIVE',
    identifier: 'emailoctopus',
    name: 'EmailOctopus',
    icon: 'integrationOptionLogos/emailoctopus.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'facebookMarketing',
    state: 'INACTIVE',
    identifier: 'facebookMarketing',
    name: 'Facebook Marketing',
    icon: 'integrationOptionLogos/facebook.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'fastbill',
    state: 'INACTIVE',
    identifier: 'fastbill',
    name: 'Fastbill',
    icon: 'integrationOptionLogos/fastbill.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'projectId',
        label: 'Project Id',
      },
    ],
  },
  {
    key: 'flexport',
    state: 'INACTIVE',
    identifier: 'flexport',
    name: 'Flexport',
    icon: 'integrationOptionLogos/flexport.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'freshcaller',
    state: 'INACTIVE',
    identifier: 'freshcaller',
    name: 'Freshcaller',
    icon: 'integrationOptionLogos/freshcaller.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'freshdesk',
    state: 'INACTIVE',
    identifier: 'freshdesk',
    name: 'Freshdesk',
    icon: 'integrationOptionLogos/freshdesk.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'domain',
        label: 'Domain',
      },
    ],
  },
  {
    key: 'freshsales',
    state: 'INACTIVE',
    identifier: 'freshsales',
    name: 'Freshsales',
    icon: 'integrationOptionLogos/freshsales.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'domain',
        label: 'Domain',
      },
    ],
  },
  {
    key: 'freshservice',
    state: 'INACTIVE',
    identifier: 'freshservice',
    name: 'Freshservice',
    icon: 'integrationOptionLogos/freshservice.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'domain',
        label: 'Domain',
      },
    ],
  },
  {
    key: 'genesys',
    state: 'INACTIVE',
    identifier: 'genesys',
    name: 'Genesys',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'region',
        label: 'Region',
      },
      {
        name: 'clientId',
        label: 'Client Id',
      },
      {
        name: 'clientSecret',
        label: 'Client Secret',
      },
    ],
  },
  {
    key: 'github',
    state: 'INACTIVE',
    identifier: 'github',
    name: 'Github',
    icon: 'integrationOptionLogos/github.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'gitlab',
    state: 'INACTIVE',
    identifier: 'gitlab',
    name: 'GitLab',
    icon: 'integrationOptionLogos/gitlab.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'gocardless',
    state: 'INACTIVE',
    identifier: 'gocardless',
    name: 'GoCardless',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
      {
        name: 'environment',
        label: 'Environment',
      },
      {
        name: 'version',
        label: 'Version',
      },
    ],
  },
  {
    key: 'gong',
    state: 'INACTIVE',
    identifier: 'gong',
    name: 'Gong',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'harvest',
    state: 'INACTIVE',
    identifier: 'harvest',
    name: 'Harvest',
    icon: 'integrationOptionLogos/harvest.svg',
    fields: [
      {
        name: 'accountId',
        label: 'Account Id',
      },
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'instagram',
    state: 'INACTIVE',
    identifier: 'instagram',
    name: 'Instagram',
    icon: 'integrationOptionLogos/instagram.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'instatus',
    state: 'INACTIVE',
    identifier: 'instatus',
    name: 'Instatus',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'intercom',
    state: 'INACTIVE',
    identifier: 'intercom',
    name: 'Intercom',
    icon: 'integrationOptionLogos/intercom.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },

  {
    key: 'klaviyo',
    state: 'INACTIVE',
    identifier: 'klaviyo',
    name: 'Klaviyo',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'kustomer',
    state: 'INACTIVE',
    identifier: 'kustomer',
    name: 'Kustomer',
    icon: 'integrationOptionLogos/kustomer.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'looker',
    state: 'INACTIVE',
    identifier: 'looker',
    name: 'Looker',
    icon: 'integrationOptionLogos/looker.svg',
    fields: [
      {
        name: 'clientId',
        label: 'Client Id',
      },
      {
        name: 'clientSecret',
        label: 'Client Secret',
      },
      {
        name: 'domain',
        label: 'Domain',
      },
    ],
  },
  {
    key: 'mailchimp',
    state: 'INACTIVE',
    identifier: 'mailchimp',
    name: 'Mailchimp',
    icon: 'integrationOptionLogos/mailchimp.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'mailjetemail',
    state: 'INACTIVE',
    identifier: 'mailjetemail',
    name: 'Mailjet Email',
    icon: 'integrationOptionLogos/mailjetemail.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'apiSecret',
        label: 'API Secret',
      },
    ],
  },
  {
    key: 'marketo',
    state: 'INACTIVE',
    identifier: 'marketo',
    name: 'Marketo',
    icon: 'integrationOptionLogos/marketo.svg',
    fields: [
      {
        name: 'clientId',
        label: 'Client Id',
      },
      {
        name: 'clientSecret',
        label: 'Client Secret',
      },
      {
        name: 'domainUrl',
        label: 'Domain Url',
      },
    ],
  },
  {
    key: 'microsoftteams',
    state: 'INACTIVE',
    identifier: 'microsoftteams',
    name: 'Microsoft Teams',
    icon: 'integrationOptionLogos/microsoftteams.svg',
    fields: [
      {
        name: 'tenantId',
        label: 'Tenant Id',
      },
      {
        name: 'clientId',
        label: 'Client Id',
      },
      {
        name: 'clientSecret',
        label: 'Client Secret',
      },
    ],
  },

  {
    key: 'monday',
    state: 'INACTIVE',
    identifier: 'monday',
    name: 'Monday',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'paypaltransaction',
    state: 'INACTIVE',
    identifier: 'paypaltransaction',
    name: 'Paypal Transaction',
    icon: 'integrationOptionLogos/paypal.svg',
    fields: [
      {
        name: 'clientId',
        label: 'Client Id',
      },
      {
        name: 'secret',
        label: 'Secret',
      },
    ],
  },
  {
    key: 'notion',
    state: 'INACTIVE',
    identifier: 'notion',
    name: 'Notion',
    icon: 'integrationOptionLogos/notion.svg',
    fields: [
      {
        name: 'internalAccessToken',
        label: 'Internal Access Token',
      },
      {
        name: 'publicClientId',
        label: 'Public Client Id',
      },
      {
        name: 'publicClientSecret',
        label: 'Public Client Secret',
      },
      {
        name: 'publicAccessToken',
        label: 'Public Access Token',
      },
    ],
  },
  {
    key: 'oraclenetsuite',
    state: 'INACTIVE',
    identifier: 'oraclenetsuite',
    name: 'Oracle Netsuite',
    icon: 'integrationOptionLogos/oraclenetsuite.svg',
    fields: [
      {
        name: 'accountId',
        label: 'Account Id',
      },
      {
        name: 'consumerKey',
        label: 'Consumer Key',
      },
      {
        name: 'consumerSecret',
        label: 'Consumer Secret',
      },
      {
        name: 'tokenId',
        label: 'Token Id',
      },
      {
        name: 'tokenSecret',
        label: 'Token Secret',
      },
    ],
  },
  {
    key: 'orb',
    state: 'INACTIVE',
    identifier: 'orb',
    name: 'Orb',
    icon: 'integrationOptionLogos/orb.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'orbit',
    state: 'INACTIVE',
    identifier: 'orbit',
    name: 'Orbit',
    icon: 'integrationOptionLogos/orbit.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'pagerduty',
    state: 'INACTIVE',
    identifier: 'pagerduty',
    name: 'PagerDuty',
    icon: 'integrationOptionLogos/pagerduty.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'paystack',
    state: 'INACTIVE',
    identifier: 'paystack',
    name: 'Paystack',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'secretKey',
        label: 'Secret Key (mandatory)',
      },
      {
        name: 'lookbackWindow',
        label: 'Lookback Window (in days)',
      },
    ],
  },
  {
    key: 'pendo',
    state: 'INACTIVE',
    identifier: 'pendo',
    name: 'Pendo',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'pipedrive',
    state: 'INACTIVE',
    identifier: 'pipedrive',
    name: 'Pipedrive',
    icon: 'integrationOptionLogos/pipedrive.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'plaid',
    state: 'INACTIVE',
    identifier: 'plaid',
    name: 'Plaid',
    icon: 'integrationOptionLogos/plaid.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },

  {
    key: 'plausible',
    state: 'INACTIVE',
    identifier: 'plausible',
    name: 'Plausible',
    icon: 'integrationOptionLogos/plausible.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'siteId',
        label: 'Site Id',
      },
    ],
  },
  {
    key: 'posthog',
    state: 'INACTIVE',
    identifier: 'posthog',
    name: 'PostHog',
    icon: 'integrationOptionLogos/posthog.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key (mandatory)',
      },
      {
        name: 'baseUrl',
        label: 'Base URL (optional)',
      },
    ],
  },
  {
    key: 'qualaroo',
    state: 'INACTIVE',
    identifier: 'qualaroo',
    name: 'Qualaroo',
    icon: 'integrationOptionLogos/qualaroo.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'quickbooks',
    state: 'INACTIVE',
    identifier: 'quickbooks',
    name: 'QuickBooks',
    icon: 'integrationOptionLogos/quickbooks.svg',
    fields: [
      {
        name: 'clientId',
        label: 'Client Id',
      },
      {
        name: 'clientSecret',
        label: 'Client Secret',
      },
      {
        name: 'realmId',
        label: 'Realm Id',
      },
      {
        name: 'refreshToken',
        label: 'Refresh Token',
      },
    ],
  },
  {
    key: 'recharge',
    state: 'INACTIVE',
    identifier: 'recharge',
    name: 'Recharge',
    icon: 'integrationOptionLogos/recharge.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'recruitee',
    state: 'INACTIVE',
    identifier: 'recruitee',
    name: 'Recruitee',
    icon: 'integrationOptionLogos/recruitee.svg',
    fields: [
      {
        name: 'companyId',
        label: 'Company Id',
      },
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },

  {
    key: 'recurly',
    state: 'INACTIVE',
    identifier: 'recurly',
    name: 'Recurly',
    icon: 'integrationOptionLogos/recurly.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'retently',
    state: 'INACTIVE',
    identifier: 'retently',
    name: 'Retently',
    icon: 'integrationOptionLogos/retently.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'salesforce',
    state: 'INACTIVE',
    identifier: 'salesforce',
    name: 'Salesforce',
    icon: 'integrationOptionLogos/salesforce.svg',
    fields: [
      {
        name: 'clientId',
        label: 'Client Id',
      },
      {
        name: 'clientSecret',
        label: 'Client Secret',
      },
      {
        name: 'refreshToken',
        label: 'Refresh Token',
      },
    ],
  },
  {
    key: 'salesloft',
    state: 'INACTIVE',
    identifier: 'salesloft',
    name: 'SalesLoft',
    icon: 'integrationOptionLogos/salesloft.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'sendgrid',
    state: 'INACTIVE',
    identifier: 'sendgrid',
    name: 'SendGrid',
    icon: 'integrationOptionLogos/sendgrid.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'sentry',
    state: 'INACTIVE',
    identifier: 'sentry',
    name: 'Sentry',
    icon: 'integrationOptionLogos/sentry.svg',
    fields: [
      {
        name: 'project',
        label: 'Project',
      },
      {
        name: 'authenticationToken',
        label: 'Authentication Token',
      },
      {
        name: 'organization',
        label: 'Organization',
      },
      {
        name: 'host',
        label: 'Host (optional)',
      },
    ],
  },
  {
    key: 'slack',
    state: 'INACTIVE',
    identifier: 'slack',
    name: 'Slack',
    icon: 'integrationOptionLogos/slack.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
      {
        name: 'channelFilter',
        label: 'Channel Filter',
      },
      {
        name: 'lookbackWindow',
        label: 'Lookback Window (in days, optional)',
      },
    ],
  },
  {
    key: 'stripe',
    state: 'INACTIVE',
    identifier: 'stripe',
    name: 'Stripe',
    icon: 'integrationOptionLogos/stripe.svg',
    fields: [
      {
        name: 'accountId',
        label: 'Account Id',
      },
      {
        name: 'secretKey',
        label: 'Secret Key',
      },
    ],
  },
  {
    key: 'surveymonkey',
    state: 'INACTIVE',
    identifier: 'surveymonkey',
    name: 'SurveyMonkey',
    icon: 'integrationOptionLogos/surveymonkey.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'talkdesk',
    state: 'INACTIVE',
    identifier: 'talkdesk',
    name: 'TalkDesk',
    icon: 'integrationOptionLogos/talkdesk.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'tiktok',
    state: 'INACTIVE',
    identifier: 'tiktok',
    name: 'TikTok',
    icon: 'integrationOptionLogos/tiktok.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
    ],
  },
  {
    key: 'todoist',
    state: 'INACTIVE',
    identifier: 'todoist',
    name: 'Todoist',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'typeform',
    state: 'INACTIVE',
    identifier: 'typeform',
    name: 'Typeform',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },
  {
    key: 'vittally',
    state: 'INACTIVE',
    identifier: 'vittally',
    name: 'Vittally',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'wrike',
    state: 'INACTIVE',
    identifier: 'wrike',
    name: 'Wrike',
    icon: 'integrationOptionLogos/wrike.svg',
    fields: [
      {
        name: 'accessToken',
        label: 'Access Token',
      },
      {
        name: 'hostUrl',
        label: 'Host URL',
      },
    ],
  },
  {
    key: 'xero',
    state: 'INACTIVE',
    identifier: 'xero',
    name: 'Xero',
    icon: '/integrationOptionLogos/customer-os.png',
    fields: [
      {
        name: 'clientId',
        label: 'Client ID',
      },
      {
        name: 'clientSecret',
        label: 'Client Secret',
      },
      {
        name: 'tenantId',
        label: 'Tenant ID',
      },
      {
        name: 'scopes',
        label: 'Scopes',
      },
    ],
  },
  {
    key: 'zendesk',
    state: 'INACTIVE',
    identifier: 'zendesksupport',
    name: 'Zendesk Support',
    icon: '/integrationOptionLogos/zendesksupport.svg',
    fields: [
      {
        name: 'apiKey',
        label: 'API key',
      },
      {
        name: 'subdomain',
        label: 'Subdomain',
      },
      {
        name: 'adminEmail',
        label: 'Admin email',
      },
    ],
    isFromIntegrationApp: true,
  },
  {
    key: 'zendeskchat',
    state: 'INACTIVE',
    identifier: 'zendeskchat',
    name: 'Zendesk Chat',
    icon: 'integrationOptionLogos/zendeskchat.svg',
    fields: [
      {
        name: 'subdomain',
        label: 'Subdomain',
      },
      {
        name: 'accessKey',
        label: 'Access Key',
      },
    ],
  },

  {
    key: 'zendesktalk',
    state: 'INACTIVE',
    identifier: 'zendesktalk',
    name: 'Zendesk Talk',
    icon: 'integrationOptionLogos/zendesktalk.svg',
    fields: [
      {
        name: 'subdomain',
        label: 'Subdomain',
      },
      {
        name: 'accessKey',
        label: 'Access Key',
      },
    ],
  },

  {
    key: 'zendesksell',
    state: 'INACTIVE',
    identifier: 'zendesksell',
    name: 'Zendesk Sell',
    icon: 'integrationOptionLogos/zendesksell.svg',
    fields: [
      {
        name: 'apiToken',
        label: 'API Token',
      },
    ],
  },

  {
    key: 'zendesksunshine',
    state: 'INACTIVE',
    identifier: 'zendesksunshine',
    name: 'Zendesk Sunshine',
    icon: 'integrationOptionLogos/zendesksunshine.svg',
    fields: [
      {
        name: 'subdomain',
        label: 'Subdomain',
      },
      {
        name: 'apiToken',
        label: 'API Token',
      },
      {
        name: 'email',
        label: 'Email',
      },
    ],
  },

  {
    key: 'zenefits',
    state: 'INACTIVE',
    identifier: 'zenefits',
    name: 'Zenefits',
    icon: 'integrationOptionLogos/zenefits.svg',
    fields: [
      {
        name: 'token',
        label: 'Token',
      },
    ],
  },
];
