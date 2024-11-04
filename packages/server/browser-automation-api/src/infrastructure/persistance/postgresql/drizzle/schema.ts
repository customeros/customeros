import {
  pgTable,
  varchar,
  timestamp,
  bigserial,
  bigint,
  uuid,
  index,
  text,
  boolean,
  uniqueIndex,
  numeric,
  integer,
  date,
  serial,
  foreignKey,
  unique,
  primaryKey,
  pgEnum,
  customType,
} from "drizzle-orm/pg-core";
import { sql } from "drizzle-orm";

const bytea = customType<{
  data: Buffer;
  notNull: false;
  default: false;
}>({
  dataType() {
    return "bytea";
  },
});

export const browserAutomationRunStatus = pgEnum(
  "browser_automation_run_status",
  [
    "SCHEDULED",
    "RUNNING",
    "COMPLETED",
    "FAILED",
    "CANCELLED",
    "PROCESSED",
    "RETRYING",
  ],
);
export const browserAutomationRunTrigger = pgEnum(
  "browser_automation_run_trigger",
  ["MANUAL", "SCHEDULER"],
);
export const browserAutomationRunType = pgEnum("browser_automation_run_type", [
  "FIND_CONNECTIONS",
  "SEND_CONNECTION_REQUEST",
  "SEND_MESSAGE",
  "FIND_COMPANY_PEOPLE",
  "DOWNLOAD_CONNECTIONS",
  "GET_MESSAGES",
  "CHECK_CONNECTION_STATUS",
]);
export const browserConfigSessionStatus = pgEnum(
  "browser_config_session_status",
  ["VALID", "INVALID", "EXPIRED"],
);

export const tenant = pgTable("tenant", {
  name: varchar("name", { length: 255 }).primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
});

export const postmarkApiKeys = pgTable("postmark_api_keys", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  tenantName: varchar("tenant_name", { length: 255 }).notNull(),
  key: varchar("key", { length: 255 }).notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }),
});

export const eventBuffer = pgTable("event_buffer", {
  tenant: varchar("tenant", { length: 50 }),
  uuid: varchar("uuid", { length: 250 }).primaryKey().notNull(),
  expiryTimestamp: timestamp("expiry_timestamp", {
    withTimezone: true,
    mode: "string",
  }),
  createdDate: timestamp("created_date", {
    withTimezone: true,
    mode: "string",
  }).default(sql`CURRENT_TIMESTAMP`),
  eventType: varchar("event_type", { length: 250 }),
  // TODO: failed to parse database type 'bytea'
  eventData: bytea("event_data"),
  // TODO: failed to parse database type 'bytea'
  eventMetadata: bytea("event_metadata"),
  eventId: varchar("event_id", { length: 50 }),
  eventTimestamp: timestamp("event_timestamp", {
    withTimezone: true,
    mode: "string",
  }),
  eventAggregateType: varchar("event_aggregate_type", { length: 250 }),
  eventAggregateId: varchar("event_aggregate_id", { length: 250 }),
  // You can use { mode: "bigint" } if numbers are exceeding js number limitations
  eventVersion: bigint("event_version", { mode: "number" }),
});

export const tenantSettingsEmailExclusion = pgTable(
  "tenant_settings_email_exclusion",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    excludeSubject: varchar("exclude_subject", { length: 255 }),
    excludeBody: varchar("exclude_body", { length: 255 }),
  },
);

export const enrichDetailsScrapin = pgTable(
  "enrich_details_scrapin",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    flow: varchar("flow", { length: 255 }).notNull(),
    param1: varchar("param1", { length: 1000 }).default(""),
    param2: varchar("param2", { length: 1000 }),
    param3: varchar("param3", { length: 1000 }),
    param4: varchar("param4", { length: 1000 }),
    allParamsJson: text("all_params_json").default("").notNull(),
    data: text("data").default("").notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    success: boolean("success").default(false),
    personFound: boolean("person_found").default(false),
    companyFound: boolean("company_found").default(false),
  },
  (table) => {
    return {
      idxEnrichDetailsScrapinParam1: index(
        "idx_enrich_details_scrapin_param1",
      ).using("btree", table.param1.asc().nullsLast()),
    };
  },
);

export const enrichDetailsBrandfetch = pgTable("enrich_details_brandfetch", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  domain: varchar("domain", { length: 255 }).default("").notNull(),
  data: text("data").default("").notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  success: boolean("success").default(false),
});

export const rawEmail = pgTable(
  "raw_email",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    sentAt: timestamp("sent_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    externalSystem: varchar("external_system", { length: 255 }).notNull(),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    username: varchar("username", { length: 255 }).notNull(),
    state: varchar("state", { length: 255 }).notNull(),
    providerMessageId: varchar("provider_message_id", {
      length: 255,
    }).notNull(),
    messageId: varchar("message_id", { length: 255 }).notNull(),
    sentToEventStoreState: varchar("sent_to_event_store_state", {
      length: 50,
    }).notNull(),
    sentToEventStoreReason: text("sent_to_event_store_reason"),
    sentToEventStoreError: text("sent_to_event_store_error"),
    data: text("data"),
  },
  (table) => {
    return {
      idxRawEmailExternalSystem: index("idx_raw_email_external_system").using(
        "btree",
        table.externalSystem.asc().nullsLast(),
        table.tenant.asc().nullsLast(),
        table.username.asc().nullsLast(),
        table.messageId.asc().nullsLast(),
      ),
    };
  },
);

export const slackSettings = pgTable(
  "slack_settings",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    tenantName: varchar("tenant_name", { length: 255 }).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    appId: varchar("app_id", { length: 255 }),
    authedUserId: varchar("authed_user_id", { length: 255 }),
    scope: varchar("scope", { length: 255 }),
    tokenType: varchar("token_type", { length: 255 }),
    accessToken: varchar("access_token", { length: 255 }),
    botUserId: varchar("bot_user_id", { length: 255 }),
    teamId: varchar("team_id", { length: 255 }),
  },
  (table) => {
    return {
      idxTenantUk: index("idx_tenant_uk").using(
        "btree",
        table.tenantName.asc().nullsLast(),
      ),
    };
  },
);

export const aiLocationMapping = pgTable("ai_location_mapping", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  input: text("input").notNull(),
  responseJson: text("response_json").notNull(),
  aiPromptLogId: uuid("ai_prompt_log_id"),
});

export const tenantWebhookApiKeys = pgTable(
  "tenant_webhook_api_keys",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    tenantName: varchar("tenant_name", { length: 255 }).notNull(),
    key: varchar("key", { length: 255 }).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }),
    enabled: boolean("enabled").default(true),
  },
  (table) => {
    return {
      uk: uniqueIndex("tenant_webhook_api_keys_uk").using(
        "btree",
        table.key.asc().nullsLast(),
      ),
    };
  },
);

export const googleServiceAccountKeys = pgTable(
  "google_service_account_keys",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    tenantName: varchar("tenant_name", { length: 255 }).notNull(),
    key: varchar("key", { length: 255 }).notNull(),
    value: text("value"),
  },
  (table) => {
    return {
      idxTenantApiKeys: index("idx_tenant_api_keys").using(
        "btree",
        table.tenantName.asc().nullsLast(),
        table.key.asc().nullsLast(),
      ),
    };
  },
);

export const aiPromptLog = pgTable("ai_prompt_log", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  createdAt: timestamp("created_at", {
    withTimezone: true,
    mode: "string",
  }).default(sql`CURRENT_TIMESTAMP`),
  appSource: varchar("app_source", { length: 50 }).notNull(),
  provider: varchar("provider", { length: 50 }).notNull(),
  model: varchar("model", { length: 100 }).notNull(),
  promptType: varchar("prompt_type", { length: 255 }).notNull(),
  promptTemplate: text("prompt_template"),
  tenant: varchar("tenant", { length: 100 }),
  nodeId: varchar("node_id", { length: 255 }),
  nodeLabel: varchar("node_label", { length: 100 }),
  prompt: text("prompt").notNull(),
  rawResponse: text("raw_response").notNull(),
  postProcessError: boolean("post_process_error"),
  postProcessErrorMessage: text("post_process_error_message"),
});

export const syncRunWebhook = pgTable("sync_run_webhook", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  tenant: varchar("tenant", { length: 50 }),
  externalSystem: varchar("external_system", { length: 50 }),
  appSource: varchar("app_source", { length: 50 }),
  entity: varchar("entity", { length: 50 }),
  startAt: timestamp("start_at", {
    withTimezone: true,
    mode: "string",
  }).default(sql`CURRENT_TIMESTAMP`),
  endAt: timestamp("end_at", { withTimezone: true, mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  reason: text("reason"),
  // You can use { mode: "bigint" } if numbers are exceeding js number limitations
  total: bigint("total", { mode: "number" }),
  // You can use { mode: "bigint" } if numbers are exceeding js number limitations
  completed: bigint("completed", { mode: "number" }),
  // You can use { mode: "bigint" } if numbers are exceeding js number limitations
  skipped: bigint("skipped", { mode: "number" }),
  // You can use { mode: "bigint" } if numbers are exceeding js number limitations
  failed: bigint("failed", { mode: "number" }),
});

export const trackingAllowedOrigin = pgTable(
  "tracking_allowed_origin",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    createdAt: timestamp("created_at", {
      withTimezone: true,
      mode: "string",
    }).default(sql`CURRENT_TIMESTAMP`),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    origin: varchar("origin", { length: 255 }).notNull(),
  },
  (table) => {
    return {
      nameDomainIdx: uniqueIndex("name_domain_idx").using(
        "btree",
        table.tenant.asc().nullsLast(),
        table.origin.asc().nullsLast(),
      ),
    };
  },
);

export const invoiceNumbers = pgTable("invoice_numbers", {
  invoiceNumber: varchar("invoice_number", { length: 16 })
    .primaryKey()
    .notNull(),
  tenant: varchar("tenant", { length: 50 }),
  createdDate: timestamp("created_date", {
    withTimezone: true,
    mode: "string",
  }).default(sql`CURRENT_TIMESTAMP`),
  // You can use { mode: "bigint" } if numbers are exceeding js number limitations
  attempts: bigint("attempts", { mode: "number" }),
});

export const currencyRates = pgTable("currency_rates", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  currency: varchar("currency", { length: 3 }).notNull(),
  rate: numeric("rate").notNull(),
  date: timestamp("date", { withTimezone: true, mode: "string" }).notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }),
  source: varchar("source", { length: 255 }).notNull(),
});

export const techLimit = pgTable(
  "tech_limit",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    key: varchar("key", { length: 255 }).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    usageCount: bigint("usage_count", { mode: "number" }).notNull(),
  },
  (table) => {
    return {
      idxKeyUnique: uniqueIndex("idx_key_unique").using(
        "btree",
        table.key.asc().nullsLast(),
      ),
    };
  },
);

export const externalAppKeys = pgTable(
  "external_app_keys",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    app: varchar("app", { length: 255 }).notNull(),
    appKey: varchar("app_key", { length: 255 }).notNull(),
    group1: varchar("group1", { length: 255 }),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    usageCount: bigint("usage_count", { mode: "number" }).notNull(),
  },
  (table) => {
    return {
      idxExternalAppKeyUnique: uniqueIndex("idx_external_app_key_unique").using(
        "btree",
        table.app.asc().nullsLast(),
        table.appKey.asc().nullsLast(),
        table.group1.asc().nullsLast(),
      ),
    };
  },
);

export const enrichDetailsBetterContact = pgTable(
  "enrich_details_better_contact",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    requestId: varchar("request_id", { length: 255 }).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }),
    contactFirstName: varchar("contact_first_name", { length: 255 }),
    contactLastName: varchar("contact_last_name", { length: 255 }),
    contactLinkedinUrl: varchar("contact_linkedin_url", { length: 255 }),
    companyName: varchar("company_name", { length: 255 }),
    companyDomain: varchar("company_domain", { length: 255 }),
    enrichPhoneNumber: boolean("enrich_phone_number").default(false),
    request: text("request"),
    response: text("response"),
  },
);

export const userEmailImportState = pgTable(
  "user_email_import_state",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    username: varchar("username", { length: 255 }).notNull(),
    provider: varchar("provider", { length: 255 }).notNull(),
    state: varchar("state", { length: 50 }).notNull(),
    startDate: timestamp("start_date", { withTimezone: true, mode: "string" }),
    stopDate: timestamp("stop_date", { withTimezone: true, mode: "string" }),
    active: boolean("active").notNull(),
    cursor: varchar("cursor", { length: 255 }).notNull(),
  },
  (table) => {
    return {
      uqOneStatePerTenantAndUser: uniqueIndex(
        "uq_one_state_per_tenant_and_user",
      ).using(
        "btree",
        table.tenant.asc().nullsLast(),
        table.username.asc().nullsLast(),
        table.provider.asc().nullsLast(),
        table.state.asc().nullsLast(),
      ),
    };
  },
);

export const userEmailImportStateHistory = pgTable(
  "user_email_import_state_history",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    entityId: text("entity_id").notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    username: varchar("username", { length: 255 }).notNull(),
    provider: varchar("provider", { length: 255 }).notNull(),
    state: varchar("state", { length: 50 }).notNull(),
    startDate: timestamp("start_date", { withTimezone: true, mode: "string" }),
    stopDate: timestamp("stop_date", { withTimezone: true, mode: "string" }),
    active: boolean("active").notNull(),
    cursor: varchar("cursor", { length: 255 }).notNull(),
  },
);

export const slackChannelNotification = pgTable("slack_channel_notification", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }),
  tenant: varchar("tenant", { length: 255 }).notNull(),
  channelId: varchar("channel_id", { length: 255 }).notNull(),
  workflow: varchar("workflow", { length: 255 }),
});

export const personalIntegrations = pgTable(
  "personal_integrations",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    tenantName: varchar("tenant_name", { length: 255 }).notNull(),
    name: varchar("name", { length: 255 }).notNull(),
    email: varchar("email", { length: 255 }).notNull(),
    key: varchar("key", { length: 255 }).notNull(),
    active: boolean("active").default(true),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }),
  },
  (table) => {
    return {
      idxKey: uniqueIndex("idx_key").using(
        "btree",
        table.key.asc().nullsLast(),
      ),
    };
  },
);

export const tenantWebhooks = pgTable("tenant_webhooks", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  tenantName: varchar("tenant_name", { length: 255 }).notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }),
  webhookUrl: varchar("webhook_url", { length: 255 }).notNull(),
  apiKey: varchar("api_key", { length: 255 }).notNull(),
  event: varchar("event", { length: 255 }).notNull(),
  authHeaderName: varchar("auth_header_name", { length: 255 }),
  authHeaderValue: varchar("auth_header_value", { length: 255 }),
  userId: varchar("user_id", { length: 255 }),
  userFirstName: varchar("user_first_name", { length: 255 }),
  userLastName: varchar("user_last_name", { length: 255 }),
  userEmail: varchar("user_email", { length: 255 }),
});

export const personalEmailProvider = pgTable(
  "personal_email_provider",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    providerName: varchar("provider_name", { length: 255 }).notNull(),
    providerDomain: varchar("provider_domain", { length: 255 }).notNull(),
    createdAt: timestamp("created_at", {
      withTimezone: true,
      mode: "string",
    }).default(sql`CURRENT_TIMESTAMP`),
  },
  (table) => {
    return {
      idxProviderDomain: index("idx_provider_domain").using(
        "btree",
        table.providerDomain.asc().nullsLast(),
      ),
    };
  },
);

export const slackChannel = pgTable("slack_channel", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  source: varchar("source", { length: 255 }),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }),
  tenantName: varchar("tenant_name", { length: 255 }).notNull(),
  channelId: varchar("channel_id", { length: 255 }).notNull(),
  channelName: varchar("channel_name", { length: 255 }),
  organizationId: varchar("organization_id", { length: 255 }),
});

export const apiCache = pgTable("api_cache", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  createdAt: timestamp("created_at", { withTimezone: true, mode: "string" })
    .default(sql`CURRENT_TIMESTAMP`)
    .notNull(),
  tenant: varchar("tenant", { length: 100 }).notNull(),
  type: varchar("type", { length: 255 }).notNull(),
  data: text("data").notNull(),
});

export const workflow = pgTable("workflow", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  tenant: varchar("tenant", { length: 255 }).notNull(),
  workflowType: varchar("workflow_type", { length: 255 }).notNull(),
  name: varchar("name", { length: 255 }),
  condition: text("condition"),
  live: boolean("live").default(false),
  actionParam1: varchar("action_param1", { length: 255 }),
});

export const industryMapping = pgTable("industry_mapping", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  inputIndustry: varchar("input_industry", { length: 255 }).notNull(),
  outputIndustry: varchar("output_industry", { length: 255 }).notNull(),
});

export const tracking = pgTable("tracking", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  createdAt: timestamp("created_at", {
    withTimezone: true,
    mode: "string",
  }).default(sql`CURRENT_TIMESTAMP`),
  tenant: varchar("tenant", { length: 255 }),
  userId: varchar("user_id", { length: 255 }).notNull(),
  ip: varchar("ip", { length: 255 }),
  eventType: varchar("event_type", { length: 255 }),
  eventData: text("event_data"),
  // You can use { mode: "bigint" } if numbers are exceeding js number limitations
  timestamp: bigint("timestamp", { mode: "number" }),
  href: varchar("href", { length: 1000 }),
  origin: varchar("origin", { length: 255 }),
  search: varchar("search", { length: 1000 }),
  hostname: varchar("hostname", { length: 255 }),
  pathname: varchar("pathname", { length: 255 }),
  referrer: varchar("referrer", { length: 2000 }),
  userAgent: text("user_agent"),
  language: varchar("language", { length: 255 }),
  cookiesEnabled: boolean("cookies_enabled"),
  screenResolution: varchar("screen_resolution", { length: 255 }),
  state: varchar("state", { length: 50 }),
  organizationId: varchar("organization_id", { length: 255 }),
  organizationName: varchar("organization_name", { length: 255 }),
  organizationDomain: varchar("organization_domain", { length: 255 }),
  organizationWebsite: varchar("organization_website", { length: 255 }),
  notified: boolean("notified").default(false),
  notificationTry: integer("notification_try").default(0),
});

export const enrichDetailsPrefilterTracking = pgTable(
  "enrich_details_prefilter_tracking",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }),
    ip: varchar("ip", { length: 255 }),
    shouldIdentify: boolean("should_identify"),
    skipIdentifyReason: varchar("skip_identify_reason", { length: 255 }),
    response: text("response"),
  },
  (table) => {
    return {
      ipUnique: uniqueIndex("ip_unique").using(
        "btree",
        table.ip.asc().nullsLast(),
      ),
    };
  },
);

export const enrichDetailsTracking = pgTable("enrich_details_tracking", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }),
  ip: varchar("ip", { length: 255 }),
  companyName: varchar("company_name", { length: 255 }),
  companyDomain: varchar("company_domain", { length: 255 }),
  companyWebsite: varchar("company_website", { length: 255 }),
  response: text("response"),
});

export const tenantSettings = pgTable("tenant_settings", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  tenantName: varchar("tenant_name", { length: 255 }).notNull(),
  smartSheetId: varchar("smart_sheet_id", { length: 255 }),
  smartSheetAccessToken: varchar("smart_sheet_access_token", { length: 255 }),
  jiraApiToken: varchar("jira_api_token", { length: 255 }),
  jiraDomain: varchar("jira_domain", { length: 255 }),
  jiraEmail: varchar("jira_email", { length: 255 }),
  trelloApiToken: varchar("trello_api_token", { length: 255 }),
  trelloApiKey: varchar("trello_api_key", { length: 255 }),
  ahaApiUrl: varchar("aha_api_url", { length: 255 }),
  ahaApiKey: varchar("aha_api_key", { length: 255 }),
  airtablePersonalAccessToken: varchar("airtable_personal_access_token", {
    length: 255,
  }),
  amplitudeApiKey: varchar("amplitude_api_key", { length: 255 }),
  amplitudeSecretKey: varchar("amplitude_secret_key", { length: 255 }),
  asanaAccessToken: varchar("asana_access_token", { length: 255 }),
  batonApiKey: varchar("baton_api_key", { length: 255 }),
  babelforceRegionEnvironment: varchar("babelforce_region_environment", {
    length: 255,
  }),
  babelforceAccessKeyId: varchar("babelforce_access_key_id", { length: 255 }),
  babelforceAccessToken: varchar("babelforce_access_token", { length: 255 }),
  bigqueryServiceAccountKey: varchar("bigquery_service_account_key", {
    length: 255,
  }),
  braintreePublicKey: varchar("braintree_public_key", { length: 255 }),
  braintreePrivateKey: varchar("braintree_private_key", { length: 255 }),
  braintreeEnvironment: varchar("braintree_environment", { length: 255 }),
  braintreeMerchantId: varchar("braintree_merchant_id", { length: 255 }),
  callrailAccount: varchar("callrail_account", { length: 255 }),
  callrailApiToken: varchar("callrail_api_token", { length: 255 }),
  chargebeeApiKey: varchar("chargebee_api_key", { length: 255 }),
  chargebeeProductCatalog: varchar("chargebee_product_catalog", {
    length: 255,
  }),
  chargifyApiKey: varchar("chargify_api_key", { length: 255 }),
  chargifyDomain: varchar("chargify_domain", { length: 255 }),
  clickupApiKey: varchar("clickup_api_key", { length: 255 }),
  closecomApiKey: varchar("closecom_api_key", { length: 255 }),
  codaAuthToken: varchar("coda_auth_token", { length: 255 }),
  codaDocumentId: varchar("coda_document_id", { length: 255 }),
  confluenceApiToken: varchar("confluence_api_token", { length: 255 }),
  confluenceDomain: varchar("confluence_domain", { length: 255 }),
  confluenceLoginEmail: varchar("confluence_login_email", { length: 255 }),
  courierApiKey: varchar("courier_api_key", { length: 255 }),
  customerioApiKey: varchar("customerio_api_key", { length: 255 }),
  datadogApiKey: varchar("datadog_api_key", { length: 255 }),
  datadogApplicationKey: varchar("datadog_application_key", { length: 255 }),
  delightedApiKey: varchar("delighted_api_key", { length: 255 }),
  dixaApiToken: varchar("dixa_api_token", { length: 255 }),
  driftApiToken: varchar("drift_api_token", { length: 255 }),
  emailoctopusApiKey: varchar("emailoctopus_api_key", { length: 255 }),
  facebookMarketingAccessToken: varchar("facebook_marketing_access_token", {
    length: 255,
  }),
  fastbillApiKey: varchar("fastbill_api_key", { length: 255 }),
  fastbillProjectId: varchar("fastbill_project_id", { length: 255 }),
  flexportApiKey: varchar("flexport_api_key", { length: 255 }),
  freshcallerApiKey: varchar("freshcaller_api_key", { length: 255 }),
  freshdeskApiKey: varchar("freshdesk_api_key", { length: 255 }),
  freshdeskDomain: varchar("freshdesk_domain", { length: 255 }),
  freshsalesApiKey: varchar("freshsales_api_key", { length: 255 }),
  freshsalesDomain: varchar("freshsales_domain", { length: 255 }),
  freshserviceApiKey: varchar("freshservice_api_key", { length: 255 }),
  freshserviceDomain: varchar("freshservice_domain", { length: 255 }),
  genesysRegion: varchar("genesys_region", { length: 255 }),
  genesysClientId: varchar("genesys_client_id", { length: 255 }),
  genesysClientSecret: varchar("genesys_client_secret", { length: 255 }),
  githubAccessToken: varchar("github_access_token", { length: 255 }),
  gitlabAccessToken: varchar("gitlab_access_token", { length: 255 }),
  gocardlessAccessToken: varchar("gocardless_access_token", { length: 255 }),
  gocardlessEnvironment: varchar("gocardless_environment", { length: 255 }),
  gocardlessVersion: varchar("gocardless_version", { length: 255 }),
  gongApiKey: varchar("gong_api_key", { length: 255 }),
  harvestAccountId: varchar("harvest_account_id", { length: 255 }),
  harvestAccessToken: varchar("harvest_access_token", { length: 255 }),
  insightlyApiToken: varchar("insightly_api_token", { length: 255 }),
  instagramAccessToken: varchar("instagram_access_token", { length: 255 }),
  instatusApiKey: varchar("instatus_api_key", { length: 255 }),
  intercomAccessToken: varchar("intercom_access_token", { length: 255 }),
  klaviyoApiKey: varchar("klaviyo_api_key", { length: 255 }),
  kustomerApiToken: varchar("kustomer_api_token", { length: 255 }),
  lookerClientId: varchar("looker_client_id", { length: 255 }),
  lookerClientSecret: varchar("looker_client_secret", { length: 255 }),
  lookerDomain: varchar("looker_domain", { length: 255 }),
  mailchimpApiKey: varchar("mailchimp_api_key", { length: 255 }),
  mailjetEmailApiKey: varchar("mailjet_email_api_key", { length: 255 }),
  mailjetEmailApiSecret: varchar("mailjet_email_api_secret", { length: 255 }),
  marketoClientId: varchar("marketo_client_id", { length: 255 }),
  marketoClientSecret: varchar("marketo_client_secret", { length: 255 }),
  marketoDomainUrl: varchar("marketo_domain_url", { length: 255 }),
  microsoftTeamsTenantId: varchar("microsoft_teams_tenant_id", { length: 255 }),
  microsoftTeamsClientId: varchar("microsoft_teams_client_id", { length: 255 }),
  microsoftTeamsClientSecret: varchar("microsoft_teams_client_secret", {
    length: 255,
  }),
  mondayApiToken: varchar("monday_api_token", { length: 255 }),
  notionInternalAccessToken: varchar("notion_internal_access_token", {
    length: 255,
  }),
  notionPublicAccessToken: varchar("notion_public_access_token", {
    length: 255,
  }),
  notionPublicClientId: varchar("notion_public_client_id", { length: 255 }),
  notionPublicClientSecret: varchar("notion_public_client_secret", {
    length: 255,
  }),
  oracleNetsuiteAccountId: varchar("oracle_netsuite_account_id", {
    length: 255,
  }),
  oracleNetsuiteConsumerKey: varchar("oracle_netsuite_consumer_key", {
    length: 255,
  }),
  oracleNetsuiteConsumerSecret: varchar("oracle_netsuite_consumer_secret", {
    length: 255,
  }),
  oracleNetsuiteTokenId: varchar("oracle_netsuite_token_id", { length: 255 }),
  oracleNetsuiteTokenSecret: varchar("oracle_netsuite_token_secret", {
    length: 255,
  }),
  orbApiKey: varchar("orb_api_key", { length: 255 }),
  orbitApiKey: varchar("orbit_api_key", { length: 255 }),
  pagerDutyApikey: varchar("pager_duty_apikey", { length: 255 }),
  paypalTransactionClientId: varchar("paypal_transaction_client_id", {
    length: 255,
  }),
  paypalTransactionSecret: varchar("paypal_transaction_secret", {
    length: 255,
  }),
  paystackSecretKey: varchar("paystack_secret_key", { length: 255 }),
  paystackLookbackWindow: varchar("paystack_lookback_window", { length: 255 }),
  pendoApiToken: varchar("pendo_api_token", { length: 255 }),
  pipedriveApiToken: varchar("pipedrive_api_token", { length: 255 }),
  plaidAccessToken: varchar("plaid_access_token", { length: 255 }),
  plausibleApiKey: varchar("plausible_api_key", { length: 255 }),
  plausibleSiteId: varchar("plausible_site_id", { length: 255 }),
  postHogApiKey: varchar("post_hog_api_key", { length: 255 }),
  postHogBaseUrl: varchar("post_hog_base_url", { length: 255 }),
  qualarooApiKey: varchar("qualaroo_api_key", { length: 255 }),
  qualarooApiToken: varchar("qualaroo_api_token", { length: 255 }),
  quickBooksClientId: varchar("quick_books_client_id", { length: 255 }),
  quickBooksClientSecret: varchar("quick_books_client_secret", { length: 255 }),
  quickBooksRealmId: varchar("quick_books_realm_id", { length: 255 }),
  quickBooksRefreshToken: varchar("quick_books_refresh_token", { length: 255 }),
  rechargeApiToken: varchar("recharge_api_token", { length: 255 }),
  recruiteeCompanyId: varchar("recruitee_company_id", { length: 255 }),
  recruiteeApiKey: varchar("recruitee_api_key", { length: 255 }),
  recurlyApiKey: varchar("recurly_api_key", { length: 255 }),
  retentlyApiToken: varchar("retently_api_token", { length: 255 }),
  salesloftApiKey: varchar("salesloft_api_key", { length: 255 }),
  sendgridApiKey: varchar("sendgrid_api_key", { length: 255 }),
  sentryProject: varchar("sentry_project", { length: 255 }),
  sentryHost: varchar("sentry_host", { length: 255 }),
  sentryAuthenticationToken: varchar("sentry_authentication_token", {
    length: 255,
  }),
  sentryOrganization: varchar("sentry_organization", { length: 255 }),
  slackApiToken: varchar("slack_api_token", { length: 255 }),
  slackChannelFilter: varchar("slack_channel_filter", { length: 255 }),
  slackLookbackWindow: varchar("slack_lookback_window", { length: 255 }),
  stripeAccountId: varchar("stripe_account_id", { length: 255 }),
  stripeSecretKey: varchar("stripe_secret_key", { length: 255 }),
  surveySparrowAccessToken: varchar("survey_sparrow_access_token", {
    length: 255,
  }),
  surveyMonkeyAccessToken: varchar("survey_monkey_access_token", {
    length: 255,
  }),
  talkdeskApiKey: varchar("talkdesk_api_key", { length: 255 }),
  tikTokAccessToken: varchar("tik_tok_access_token", { length: 255 }),
  todoistApiToken: varchar("todoist_api_token", { length: 255 }),
  typeformApiToken: varchar("typeform_api_token", { length: 255 }),
  vittallyApiKey: varchar("vittally_api_key", { length: 255 }),
  wrikeAccessToken: varchar("wrike_access_token", { length: 255 }),
  wrikeHostUrl: varchar("wrike_host_url", { length: 255 }),
  xeroClientId: varchar("xero_client_id", { length: 255 }),
  xeroClientSecret: varchar("xero_client_secret", { length: 255 }),
  xeroTenantId: varchar("xero_tenant_id", { length: 255 }),
  xeroScopes: varchar("xero_scopes", { length: 255 }),
  zendeskApiKey: varchar("zendesk_api_key", { length: 255 }),
  zendeskSubdomain: varchar("zendesk_subdomain", { length: 255 }),
  zendeskAdminEmail: varchar("zendesk_admin_email", { length: 255 }),
  zendeskChatSubdomain: varchar("zendesk_chat_subdomain", { length: 255 }),
  zendeskChatAccessKey: varchar("zendesk_chat_access_key", { length: 255 }),
  zendeskTalkSubdomain: varchar("zendesk_talk_subdomain", { length: 255 }),
  zendeskTalkAccessKey: varchar("zendesk_talk_access_key", { length: 255 }),
  zendeskSellApiToken: varchar("zendesk_sell_api_token", { length: 255 }),
  zendeskSunshineSubdomain: varchar("zendesk_sunshine_subdomain", {
    length: 255,
  }),
  zendeskSunshineApiToken: varchar("zendesk_sunshine_api_token", {
    length: 255,
  }),
  zendeskSunshineEmail: varchar("zendesk_sunshine_email", { length: 255 }),
  zenefitsToken: varchar("zenefits_token", { length: 255 }),
  mixpanelUsername: varchar("mixpanel_username", { length: 255 }),
  mixpanelSecret: varchar("mixpanel_secret", { length: 255 }),
  mixpanelProjectId: varchar("mixpanel_project_id", { length: 255 }),
  mixpanelProjectSecret: varchar("mixpanel_project_secret", { length: 255 }),
  mixpanelProjectTimezone: varchar("mixpanel_project_timezone", {
    length: 255,
  }),
  mixpanelRegion: varchar("mixpanel_region", { length: 255 }),
  linkedinCredential: varchar("linkedin_credential", { length: 255 }),
  linkedinCredentialIv: varchar("linkedin_credential_iv", { length: 32 }),
  linkedinPassword: varchar("linkedin_password", { length: 255 }),
  linkedinPasswordIv: varchar("linkedin_password_iv", { length: 32 }),
});

export const tenantSettingsOpportunityStage = pgTable(
  "tenant_settings_opportunity_stage",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }),
    visible: boolean("visible").notNull(),
    val: text("val").notNull(),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    idx: bigint("idx", { mode: "number" }).notNull(),
    label: varchar("label", { length: 255 }).notNull(),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    likelihoodRate: bigint("likelihood_rate", { mode: "number" })
      .default(0)
      .notNull(),
  },
);

export const tenantSettingsMailbox = pgTable("tenant_settings_mailbox", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  tenant: varchar("tenant", { length: 255 }).notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }),
  mailboxUsername: varchar("mailbox_username", { length: 255 }),
  mailboxPassword: varchar("mailbox_password", { length: 255 }),
  domain: varchar("domain", { length: 255 }),
  userName: varchar("user_name", { length: 255 }),
  minMinutesBetweenEmails: integer("min_minutes_between_emails"),
  maxMinutesBetweenEmails: integer("max_minutes_between_emails"),
});

export const emailLookup = pgTable("email_lookup", {
  id: varchar("id", { length: 64 }).primaryKey().notNull(),
  tenant: varchar("tenant", { length: 255 }).notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  trackerDomain: varchar("tracker_domain", { length: 255 }),
  messageId: varchar("message_id", { length: 64 }).notNull(),
  linkId: varchar("link_id", { length: 64 }).notNull(),
  redirectUrl: varchar("redirect_url", { length: 255 }).notNull(),
  campaign: varchar("campaign", { length: 255 }).notNull(),
  type: varchar("type", { length: 32 }).notNull(),
  recipientId: varchar("recipient_id", { length: 255 }),
  trackOpens: boolean("track_opens").notNull(),
  trackClicks: boolean("track_clicks").notNull(),
  unsubscribeUrl: varchar("unsubscribe_url", { length: 255 }),
});

export const emailTracking = pgTable("email_tracking", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  tenant: varchar("tenant", { length: 255 }).notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  timestamp: timestamp("timestamp", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  messageId: varchar("message_id", { length: 64 }).notNull(),
  linkId: varchar("link_id", { length: 64 }),
  recipientId: varchar("recipient_id", { length: 255 }),
  campaign: varchar("campaign", { length: 255 }),
  eventType: varchar("event_type", { length: 255 }).notNull(),
  ip: varchar("ip", { length: 255 }),
});

export const cacheIpData = pgTable(
  "cache_ip_data",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    ip: varchar("ip", { length: 255 }).notNull(),
    data: text("data"),
  },
  (table) => {
    return {
      idxCacheIpDataIp: uniqueIndex("idx_cache_ip_data_ip").using(
        "btree",
        table.ip.asc().nullsLast(),
      ),
    };
  },
);

export const cacheIpHunter = pgTable(
  "cache_ip_hunter",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    ip: varchar("ip", { length: 255 }).notNull(),
    data: text("data"),
  },
  (table) => {
    return {
      idxCacheIpHunterIp: uniqueIndex("idx_cache_ip_hunter_ip").using(
        "btree",
        table.ip.asc().nullsLast(),
      ),
    };
  },
);

export const appKeys = pgTable("app_keys", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  appId: varchar("app_id", { length: 255 }).notNull(),
  key: varchar("key", { length: 255 }).notNull(),
  active: boolean("active").notNull(),
});

export const cacheEmailValidation = pgTable(
  "cache_email_validation",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    email: varchar("email", { length: 255 }).notNull(),
    normalizedEmail: varchar("normalized_email", { length: 255 }),
    username: varchar("username", { length: 255 }),
    domain: varchar("domain", { length: 255 }),
    isMailboxFull: boolean("is_mailbox_full"),
    isRoleAccount: boolean("is_role_account"),
    isFreeAccount: boolean("is_free_account"),
    smtpSuccess: boolean("smtp_success"),
    responseCode: varchar("response_code", { length: 255 }),
    errorCode: varchar("error_code", { length: 255 }),
    description: text("description"),
    tlsRequired: boolean("tls_required"),
    retryValidation: boolean("retry_validation"),
    deliverable: varchar("deliverable", { length: 16 }),
    healthIsGreylisted: boolean("health_is_greylisted"),
    healthIsBlacklisted: boolean("health_is_blacklisted"),
    healthServerIp: varchar("health_server_ip", { length: 255 }),
    healthFromEmail: varchar("health_from_email", { length: 255 }),
    healthRetryAfter: integer("health_retry_after"),
    alternateEmail: varchar("alternate_email", { length: 255 }),
    error: varchar("error", { length: 255 }),
    data: text("data"),
  },
  (table) => {
    return {
      idxCacheEmailValidationEmail: uniqueIndex(
        "idx_cache_email_validation_email",
      ).using("btree", table.email.asc().nullsLast()),
    };
  },
);

export const cacheEmailValidationDomain = pgTable(
  "cache_email_validation_domain",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    domain: varchar("domain", { length: 255 }).notNull(),
    isCatchAll: boolean("is_catch_all"),
    isFirewalled: boolean("is_firewalled"),
    provider: varchar("provider", { length: 255 }),
    firewall: varchar("firewall", { length: 255 }),
    hasMxRecord: boolean("has_mx_record"),
    hasSpfRecord: boolean("has_spf_record"),
    error: varchar("error", { length: 255 }),
    data: text("data"),
    canConnectSmtp: boolean("can_connect_smtp"),
    tlsRequired: boolean("tls_required"),
    responseCode: varchar("response_code", { length: 255 }),
    errorCode: varchar("error_code", { length: 255 }),
    description: text("description"),
    healthIsGreylisted: boolean("health_is_greylisted"),
    healthIsBlacklisted: boolean("health_is_blacklisted"),
    healthServerIp: varchar("health_server_ip", { length: 255 }),
    healthFromEmail: varchar("health_from_email", { length: 255 }),
    healthRetryAfter: integer("health_retry_after"),
    isPrimaryDomain: boolean("is_primary_domain"),
    primaryDomain: varchar("primary_domain", { length: 255 }),
  },
  (table) => {
    return {
      idxCacheEmailValidationDomainDomain: uniqueIndex(
        "idx_cache_email_validation_domain_domain",
      ).using("btree", table.domain.asc().nullsLast()),
    };
  },
);

export const cacheEmailScrubby = pgTable("cache_email_scrubby", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  checkedAt: timestamp("checked_at", { mode: "string" }),
  email: varchar("email", { length: 255 }).notNull(),
  status: varchar("status", { length: 255 }),
});

export const statsApiCalls = pgTable(
  "stats_api_calls",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    api: varchar("api", { length: 255 }).notNull(),
    day: date("day").notNull(),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    calls: bigint("calls", { mode: "number" }).notNull(),
  },
  (table) => {
    return {
      idxStatsApiCallsUnique: uniqueIndex("idx_stats_api_calls_unique").using(
        "btree",
        table.tenant.asc().nullsLast(),
        table.api.asc().nullsLast(),
        table.day.asc().nullsLast(),
      ),
    };
  },
);

export const cosApiEnrichPersonTempResult = pgTable(
  "cos_api_enrich_person_temp_result",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    tenant: varchar("tenant", { length: 255 }),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    scrapinRecordId: bigint("scrapin_record_id", { mode: "number" }),
    bettercontactRecordId: varchar("bettercontact_record_id", { length: 255 }),
  },
);

export const organizationWebsiteHostingPlatform = pgTable(
  "organization_website_hosting_platform",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    urlPattern: varchar("url_pattern", { length: 255 }).notNull(),
    createdAt: timestamp("created_at", {
      withTimezone: true,
      mode: "string",
    }).default(sql`CURRENT_TIMESTAMP`),
  },
);

export const cacheEmailTrueinbox = pgTable("cache_email_trueinbox", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  email: varchar("email", { length: 255 }).notNull(),
  result: varchar("result", { length: 255 }),
  data: text("data"),
});

export const cacheEmailEnrow = pgTable("cache_email_enrow", {
  id: uuid("id").defaultRandom().primaryKey().notNull(),
  requestId: varchar("request_id", { length: 255 }).notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }),
  email: varchar("email", { length: 255 }).notNull(),
  qualification: varchar("qualification", { length: 255 }),
  data: text("data"),
});

export const emailValidationRecord = pgTable(
  "email_validation_record",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    requestId: varchar("request_id", { length: 255 }).notNull(),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    email: varchar("email", { length: 255 }).notNull(),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    priority: bigint("priority", { mode: "number" }).default(0),
    data: text("data").notNull(),
    verifyCatchAll: boolean("verify_catch_all").default(false),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }),
  },
  (table) => {
    return {
      idxEmailRequestId: uniqueIndex("idx_email_request_id").using(
        "btree",
        table.requestId.asc().nullsLast(),
        table.email.asc().nullsLast(),
      ),
    };
  },
);

export const emailValidationRequestBulk = pgTable(
  "email_validation_request_bulk",
  {
    requestId: uuid("request_id").defaultRandom().primaryKey().notNull(),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    totalEmails: bigint("total_emails", { mode: "number" }).notNull(),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    deliverableEmails: bigint("deliverable_emails", { mode: "number" }).default(
      0,
    ),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    undeliverableEmails: bigint("undeliverable_emails", {
      mode: "number",
    }).default(0),
    status: varchar("status", { length: 50 }).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }),
    fileName: varchar("file_name", { length: 255 }).notNull(),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    priority: bigint("priority", { mode: "number" }).default(0),
    verifyCatchAll: boolean("verify_catch_all").default(false),
    fileStoreId: varchar("file_store_id", { length: 255 }),
  },
);

export const apiBillableEvent = pgTable(
  "api_billable_event",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    event: varchar("event", { length: 255 }).notNull(),
    externalId: varchar("external_id", { length: 255 }),
    referenceData: text("reference_data"),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
  },
  (table) => {
    return {
      idxApiBillableEventExternalId: index(
        "idx_api_billable_event_external_id",
      ).using("btree", table.externalId.asc().nullsLast()),
      idxTenantEvent: index("idx_tenant_event").using(
        "btree",
        table.event.asc().nullsLast(),
      ),
    };
  },
);

export const mailstackDomain = pgTable(
  "mailstack_domain",
  {
    id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    domain: varchar("domain", { length: 255 }).notNull(),
    configured: boolean("configured").default(false).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }),
    active: boolean("active").default(true).notNull(),
    dkimPublic: text("dkim_public"),
    dkimPrivate: text("dkim_private"),
  },
  (table) => {
    return {
      idxMailstackDomainDomain: uniqueIndex(
        "idx_mailstack_domain_domain",
      ).using("btree", table.domain.asc().nullsLast()),
    };
  },
);

export const assignedProxies = pgTable("assigned_proxies", {
  id: serial("id").primaryKey().notNull(),
  proxyPoolId: integer("proxy_pool_id").notNull(),
  userId: varchar("user_id", { length: 36 }).notNull(),
  tenant: text("tenant").notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
});

export const proxyPool = pgTable("proxy_pool", {
  id: serial("id").primaryKey().notNull(),
  url: text("url").notNull(),
  username: text("username").notNull(),
  password: text("password").notNull(),
  enabled: boolean("enabled").default(true),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
});

export const browserAutomationRuns = pgTable("browser_automation_runs", {
  id: serial("id").primaryKey().notNull(),
  browserConfigId: integer("browser_config_id").notNull(),
  userId: varchar("user_id", { length: 36 }).notNull(),
  tenant: text("tenant").notNull(),
  type: browserAutomationRunType("type").notNull(),
  payload: text("payload"),
  status: browserAutomationRunStatus("status").default("SCHEDULED").notNull(),
  scheduledAt: timestamp("scheduled_at", { mode: "string" }),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  startedAt: timestamp("started_at", { mode: "string" }),
  finishedAt: timestamp("finished_at", { mode: "string" }),
  runDuration: integer("run_duration"),
  retryCount: integer("retry_count").default(0),
  triggeredBy: browserAutomationRunTrigger("triggered_by"),
  priority: integer("priority").default(0),
  logLocation: text("log_location"),
});

export const emailMessage = pgTable(
  "email_message",
  {
    id: uuid("id").defaultRandom().primaryKey().notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    status: varchar("status", { length: 50 }).notNull(),
    sentAt: timestamp("sent_at", { mode: "string" }),
    error: text("error"),
    uniqueInternalIdentifier: varchar("unique_internal_identifier", {
      length: 255,
    }),
    tenant: varchar("tenant", { length: 255 }).notNull(),
    producerId: varchar("producer_id", { length: 255 }).notNull(),
    producerType: varchar("producer_type", { length: 255 }).notNull(),
    from: varchar("from", { length: 255 }),
    fromProvider: varchar("from_provider", { length: 255 }),
    to: text("to"),
    cc: text("cc"),
    bcc: text("bcc"),
    subject: text("subject"),
    content: text("content"),
    replyTo: text("reply_to"),
    providerMessageId: varchar("provider_message_id", {
      length: 255,
    }).notNull(),
    providerThreadId: varchar("provider_thread_id", { length: 255 }).notNull(),
    providerInReplyTo: varchar("provider_in_reply_to", {
      length: 255,
    }).notNull(),
    providerReferences: varchar("provider_references", {
      length: 255,
    }).notNull(),
  },
  (table) => {
    return {
      uniqueInternalIdentifier: index("unique_internal_identifier").using(
        "btree",
        table.uniqueInternalIdentifier.asc().nullsLast(),
      ),
    };
  },
);

export const browserAutomationRunResults = pgTable(
  "browser_automation_run_results",
  {
    id: serial("id").primaryKey().notNull(),
    runId: integer("run_id").notNull(),
    type: varchar("type", { length: 50 }).notNull(),
    resultData: text("result_data"),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    isProcessed: boolean("is_processed").default(false),
  },
  (table) => {
    return {
      browserAutomationRunResultsRunIdBrowserAutomationRunsI: foreignKey({
        columns: [table.runId],
        foreignColumns: [browserAutomationRuns.id],
        name: "browser_automation_run_results_run_id_browser_automation_runs_i",
      }).onDelete("cascade"),
    };
  },
);

export const browserAutomationRunErrors = pgTable(
  "browser_automation_run_errors",
  {
    id: serial("id").primaryKey().notNull(),
    runId: integer("run_id").notNull(),
    occurredAt: timestamp("occurred_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    errorType: varchar("error_type", { length: 100 }).notNull(),
    errorMessage: text("error_message").notNull(),
    errorDetails: text("error_details"),
    errorCode: varchar("error_code", { length: 50 }),
  },
  (table) => {
    return {
      browserAutomationRunErrorsRunIdBrowserAutomationRunsId: foreignKey({
        columns: [table.runId],
        foreignColumns: [browserAutomationRuns.id],
        name: "browser_automation_run_errors_run_id_browser_automation_runs_id",
      }).onDelete("cascade"),
    };
  },
);

export const browserConfigs = pgTable(
  "browser_configs",
  {
    id: serial("id").primaryKey().notNull(),
    userId: varchar("user_id", { length: 36 }).notNull(),
    tenant: text("tenant").notNull(),
    cookies: text("cookies"),
    userAgent: text("user_agent"),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    sessionStatus: browserConfigSessionStatus("session_status")
      .default("VALID")
      .notNull(),
  },
  (table) => {
    return {
      browserConfigsUserIdUnique: unique("browser_configs_user_id_unique").on(
        table.userId,
      ),
    };
  },
);

export const tableViewDefinition = pgTable("table_view_definition", {
  id: bigserial("id", { mode: "bigint" }).primaryKey().notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  tenant: varchar("tenant", { length: 255 }).notNull(),
  userId: varchar("user_id", { length: 255 }),
  tableId: varchar("table_id", { length: 255 }).default("").notNull(),
  tableType: varchar("table_type", { length: 255 }).notNull(),
  tableName: varchar("table_name", { length: 255 }).notNull(),
  // You can use { mode: "bigint" } if numbers are exceeding js number limitations
  position: bigint("position", { mode: "number" }).notNull(),
  icon: varchar("icon", { length: 255 }),
  filters: text("filters"),
  sorting: text("sorting"),
  columns: text("columns"),
  isPreset: boolean("is_preset").default(false).notNull(),
  isShared: boolean("is_shared").default(false).notNull(),
});

export const customerOsIds = pgTable(
  "customer_os_ids",
  {
    tenant: varchar("tenant", { length: 50 }).notNull(),
    customerOsId: varchar("customer_os_id", { length: 30 }).notNull(),
    entity: varchar("entity", { length: 30 }),
    entityId: varchar("entity_id", { length: 50 }),
    createdDate: timestamp("created_date", {
      withTimezone: true,
      mode: "string",
    }).default(sql`CURRENT_TIMESTAMP`),
    // You can use { mode: "bigint" } if numbers are exceeding js number limitations
    attempts: bigint("attempts", { mode: "number" }),
  },
  (table) => {
    return {
      customerOsIdsPkey: primaryKey({
        columns: [table.tenant, table.customerOsId],
        name: "customer_os_ids_pkey",
      }),
    };
  },
);

export const oauthToken = pgTable(
  "oauth_token",
  {
    provider: varchar("provider", { length: 255 }).notNull(),
    tenantName: varchar("tenant_name", { length: 255 }).notNull(),
    emailAddress: varchar("email_address", { length: 255 }).notNull(),
    type: varchar("type", { length: 50 }),
    playerIdentityId: varchar("player_identity_id", { length: 255 }).notNull(),
    accessToken: text("access_token"),
    refreshToken: text("refresh_token"),
    needsManualRefresh: boolean("needs_manual_refresh").default(false),
    idToken: text("id_token"),
    expiresAt: timestamp("expires_at", { mode: "string" }),
    scope: text("scope"),
    gmailSyncEnabled: boolean("gmail_sync_enabled").default(false),
    googleCalendarSyncEnabled: boolean("google_calendar_sync_enabled").default(
      false,
    ),
  },
  (table) => {
    return {
      idxPrimary: index("idx_primary").using(
        "btree",
        table.provider.asc().nullsLast(),
        table.tenantName.asc().nullsLast(),
        table.emailAddress.asc().nullsLast(),
      ),
      oauthTokenPkey: primaryKey({
        columns: [table.provider, table.tenantName, table.emailAddress],
        name: "oauth_token_pkey",
      }),
    };
  },
);

