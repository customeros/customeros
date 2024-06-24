import orb from '@assets/integrationOptionLogos/orb.svg';
import aha from '@assets/integrationOptionLogos/aha.svg';
import jira from '@assets/integrationOptionLogos/jira.svg';
import coda from '@assets/integrationOptionLogos/coda.svg';
import dixa from '@assets/integrationOptionLogos/dixa.svg';
import wrike from '@assets/integrationOptionLogos/wrike.svg';
import slack from '@assets/integrationOptionLogos/slack.svg';
import plaid from '@assets/integrationOptionLogos/plaid.svg';
import orbit from '@assets/integrationOptionLogos/orbit.svg';
import asana from '@assets/integrationOptionLogos/asana.svg';
import drift from '@assets/integrationOptionLogos/drift.svg';
import attio from '@assets/integrationOptionLogos/attio.svg';
import tiktok from '@assets/integrationOptionLogos/tiktok.svg';
import stripe from '@assets/integrationOptionLogos/stripe.svg';
import sentry from '@assets/integrationOptionLogos/sentry.svg';
import notion from '@assets/integrationOptionLogos/notion.svg';
import looker from '@assets/integrationOptionLogos/looker.svg';
import gitlab from '@assets/integrationOptionLogos/gitlab.svg';
import github from '@assets/integrationOptionLogos/github.svg';
import trello from '@assets/integrationOptionLogos/trello.svg';
import recurly from '@assets/integrationOptionLogos/recurly.svg';
import posthog from '@assets/integrationOptionLogos/posthog.svg';
import marketo from '@assets/integrationOptionLogos/marketo.svg';
import harvest from '@assets/integrationOptionLogos/harvest.svg';
import hubspot from '@assets/integrationOptionLogos/hubspot.svg';
import clickup from '@assets/integrationOptionLogos/clickup.svg';
import courier from '@assets/integrationOptionLogos/courier.svg';
import datadog from '@assets/integrationOptionLogos/datadog.svg';
import zenefits from '@assets/integrationOptionLogos/zenefits.svg';
import talkdesk from '@assets/integrationOptionLogos/talkdesk.svg';
import sendgrid from '@assets/integrationOptionLogos/sendgrid.svg';
import retently from '@assets/integrationOptionLogos/retently.svg';
import recharge from '@assets/integrationOptionLogos/recharge.svg';
import qualaroo from '@assets/integrationOptionLogos/qualaroo.svg';
import kustomer from '@assets/integrationOptionLogos/kustomer.svg';
import intercom from '@assets/integrationOptionLogos/intercom.svg';
import unthread from '@assets/integrationOptionLogos/unthread.png';
import airtable from '@assets/integrationOptionLogos/airtable.svg';
import bigquery from '@assets/integrationOptionLogos/bigquery.svg';
import chargify from '@assets/integrationOptionLogos/chargify.svg';
import facebook from '@assets/integrationOptionLogos/facebook.svg';
import fastbill from '@assets/integrationOptionLogos/fastbill.svg';
import flexport from '@assets/integrationOptionLogos/flexport.svg';
import mixpanel from '@assets//integrationOptionLogos/mixpanel.svg';
import gsuite from '@assets/integrationOptionLogos/google-icon.svg';
import salesloft from '@assets/integrationOptionLogos/salesloft.svg';
import recruitee from '@assets/integrationOptionLogos/recruitee.svg';
import plausible from '@assets/integrationOptionLogos/plausible.svg';
import pipedrive from '@assets/integrationOptionLogos/pipedrive.svg';
import pagerduty from '@assets/integrationOptionLogos/pagerduty.svg';
import mailchimp from '@assets/integrationOptionLogos/mailchimp.svg';
import instagram from '@assets/integrationOptionLogos/instagram.svg';
import freshdesk from '@assets/integrationOptionLogos/freshdesk.svg';
import amplitude from '@assets/integrationOptionLogos/amplitude.svg';
import braintree from '@assets/integrationOptionLogos/braintree.svg';
import chargebee from '@assets/integrationOptionLogos/chargebee.svg';
import delighted from '@assets/integrationOptionLogos/delighted.svg';
import salesforce from '@assets/integrationOptionLogos/salesforce.svg';
import quickbooks from '@assets/integrationOptionLogos/quickbooks.svg';
import freshsales from '@assets/integrationOptionLogos/freshsales.svg';
import smartsheet from '@assets/integrationOptionLogos/smartsheet.svg';
import closedotcom from '@assets/integrationOptionLogos/close.com.svg';
import confluence from '@assets/integrationOptionLogos/confluence.svg';
import customeros from '@assets/integrationOptionLogos/customer-os.png';
import customerio from '@assets/integrationOptionLogos/customer-io.svg';
import zendesksell from '@assets/integrationOptionLogos/zendesksell.svg';
import zendesktalk from '@assets/integrationOptionLogos/zendesktalk.svg';
import zendeskchat from '@assets/integrationOptionLogos/zendeskchat.svg';
import zendesk from '@assets//integrationOptionLogos/zendesksupport.svg';
import freshcaller from '@assets/integrationOptionLogos/freshcaller.svg';
import paypaltransaction from '@assets/integrationOptionLogos/paypal.svg';
import surveymonkey from '@assets/integrationOptionLogos/surveymonkey.svg';
import mailjetemail from '@assets/integrationOptionLogos/mailjetemail.svg';
import freshservice from '@assets/integrationOptionLogos/freshservice.svg';
import emailoctopus from '@assets/integrationOptionLogos/emailoctopus.svg';
import oraclenetsuite from '@assets/integrationOptionLogos/oraclenetsuite.svg';
import microsoftteams from '@assets/integrationOptionLogos/microsoftteams.svg';
import zendesksunshine from '@assets/integrationOptionLogos/zendesksunshine.svg';

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
    key: 'attio',
    state: 'INACTIVE',
    icon: attio,
    identifier: 'attio',
    name: 'Attio',
    fields: [],
    isFromIntegrationApp: true,
  },
  {
    key: 'unthread',
    state: 'INACTIVE',
    icon: unthread,
    identifier: 'unthread',
    name: 'Unthread',
    fields: [],
    isFromIntegrationApp: true,
  },
  {
    key: 'customeros-custom-payment-provider',
    state: 'INACTIVE',
    icon: customeros,
    identifier: 'customeros-custom-payment-provider',
    name: 'Custom payment provider',
    fields: [],
    isFromIntegrationApp: true,
  },
  {
    key: 'gsuite',
    state: 'INACTIVE',
    icon: gsuite,
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
    icon: hubspot,
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
    icon: smartsheet,
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
    icon: jira,
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
    icon: trello,
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
    icon: aha,

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
    icon: airtable,
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
    icon: amplitude,
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
    icon: asana,
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
    icon: customeros,
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
    icon: customeros,
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
    icon: bigquery,
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
    icon: braintree,
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
    icon: customeros,
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
    icon: chargebee,
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
    icon: chargify,
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
    icon: clickup,
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
  },
  {
    key: 'close-oauth',
    state: 'INACTIVE',
    identifier: 'close-oauth',
    name: 'Close.com',
    icon: closedotcom,
    fields: [
      {
        name: 'apiKey',
        label: 'API Key',
      },
    ],
    isFromIntegrationApp: true,
  },
  {
    key: 'coda',
    state: 'INACTIVE',
    identifier: 'coda',
    name: 'Coda',
    icon: coda,
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
    icon: confluence,
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
    icon: courier,
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
    icon: customerio,
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
    icon: datadog,
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
    icon: delighted,
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
    icon: dixa,
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
    icon: drift,
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
    icon: emailoctopus,
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
    icon: facebook,
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
    icon: fastbill,
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
    icon: flexport,
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
    icon: freshcaller,
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
    icon: freshdesk,
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
    icon: freshsales,
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
    icon: freshservice,
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
    icon: customeros,
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
    icon: github,
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
    icon: gitlab,
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
    icon: customeros,
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
    icon: customeros,
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
    icon: harvest,
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
    icon: instagram,
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
    icon: customeros,
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
    icon: intercom,
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
    icon: customeros,
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
    icon: kustomer,
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
    icon: looker,
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
    icon: mailchimp,
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
    icon: mailjetemail,
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
    icon: marketo,
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
    icon: microsoftteams,
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
    key: 'mixpanel',
    state: 'INACTIVE',
    icon: mixpanel,
    identifier: 'mixpanel',
    name: 'Mixpanel',
    fields: [
      {
        name: 'username',
        label: 'Username',
      },
      {
        name: 'secret',
        label: 'Secret',
      },
      {
        name: 'projectId',
        label: 'Project ID',
      },
      {
        name: 'projectSecret',
        label: 'Project Secret',
      },
      {
        name: 'projectTimezone',
        label: 'Project Timezone',
      },
      {
        name: 'region',
        label: 'Region',
      },
    ],
  },
  {
    key: 'monday',
    state: 'INACTIVE',
    identifier: 'monday',
    name: 'Monday',
    icon: customeros,
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
    icon: paypaltransaction,
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
    icon: notion,
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
    icon: oraclenetsuite,
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
    icon: orb,
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
    icon: orbit,
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
    icon: pagerduty,
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
    icon: customeros,
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
    icon: customeros,
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
    icon: pipedrive,
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
    icon: plaid,
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
    icon: plausible,
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
    icon: posthog,
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
    icon: qualaroo,
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
    icon: quickbooks,
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
    icon: recharge,
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
    icon: recruitee,
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
    icon: recurly,
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
    icon: retently,
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
    icon: salesforce,
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
    isFromIntegrationApp: true,
  },
  {
    key: 'salesloft',
    state: 'INACTIVE',
    identifier: 'salesloft',
    name: 'SalesLoft',
    icon: salesloft,
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
    icon: sendgrid,
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
    icon: sentry,
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
    icon: slack,
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
    icon: stripe,
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
    isFromIntegrationApp: true,
  },
  {
    key: 'surveymonkey',
    state: 'INACTIVE',
    identifier: 'surveymonkey',
    name: 'SurveyMonkey',
    icon: surveymonkey,
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
    icon: talkdesk,
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
    icon: tiktok,
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
    icon: customeros,
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
    icon: customeros,
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
    icon: customeros,
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
    icon: wrike,
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
    icon: customeros,
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
    icon: zendesk,
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
    icon: zendeskchat,
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
    icon: zendesktalk,
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
    icon: zendesksell,
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
    icon: zendesksunshine,
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
    icon: zenefits,
    fields: [
      {
        name: 'token',
        label: 'Token',
      },
    ],
  },
];
