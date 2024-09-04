-- Current sql file was generated after introspecting the database
-- If you want to run this migration please uncomment this code before executing migrations
/*
CREATE TABLE IF NOT EXISTS "app_keys" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"app_id" varchar(255) NOT NULL,
	"key" varchar(255) NOT NULL,
	"active" boolean NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "google_service_account_keys" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant_name" varchar(255) NOT NULL,
	"key" varchar(255) NOT NULL,
	"value" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "event_buffer" (
	"tenant" varchar(50),
	"uuid" varchar(250) PRIMARY KEY NOT NULL,
	"expiry_timestamp" timestamp with time zone,
	"created_date" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"event_type" varchar(250),
	"event_data" "bytea",
	"event_metadata" "bytea",
	"event_id" varchar(50),
	"event_timestamp" timestamp with time zone,
	"event_aggregate_type" varchar(250),
	"event_aggregate_id" varchar(250),
	"event_version" bigint
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "email_exclusion" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant" varchar(255) NOT NULL,
	"exclude_subject" varchar(255) NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "enrich_details_scrapin" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"flow" varchar(255) NOT NULL,
	"param1" varchar(1000) DEFAULT '' NOT NULL,
	"param2" varchar(1000),
	"param3" varchar(1000),
	"param4" varchar(1000),
	"all_params_json" text DEFAULT '' NOT NULL,
	"data" text DEFAULT '' NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"success" boolean DEFAULT false,
	"person_found" boolean DEFAULT false,
	"company_found" boolean DEFAULT false
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "user_email_import_state" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant" varchar(255) NOT NULL,
	"username" varchar(255) NOT NULL,
	"provider" varchar(255) NOT NULL,
	"state" varchar(50) NOT NULL,
	"start_date" timestamp with time zone,
	"stop_date" timestamp with time zone,
	"active" boolean NOT NULL,
	"cursor" varchar(255) NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "user_email_import_state_history" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"entity_id" text NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"tenant" varchar(255) NOT NULL,
	"username" varchar(255) NOT NULL,
	"provider" varchar(255) NOT NULL,
	"state" varchar(50) NOT NULL,
	"start_date" timestamp with time zone,
	"stop_date" timestamp with time zone,
	"active" boolean NOT NULL,
	"cursor" varchar(255) NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tenant_webhook_api_keys" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"tenant_name" varchar(255) NOT NULL,
	"key" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"enabled" boolean DEFAULT true
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tenant_webhooks" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"tenant_name" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"webhook_url" varchar(255) NOT NULL,
	"api_key" varchar(255) NOT NULL,
	"event" varchar(255) NOT NULL,
	"auth_header_name" varchar(255),
	"auth_header_value" varchar(255),
	"user_id" varchar(255),
	"user_first_name" varchar(255),
	"user_last_name" varchar(255),
	"user_email" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tracking_allowed_origin" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"tenant" varchar(255) NOT NULL,
	"origin" varchar(255) NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "table_view_definition" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"tenant" varchar(255) NOT NULL,
	"user_id" varchar(255),
	"table_id" varchar(255) DEFAULT '' NOT NULL,
	"table_type" varchar(255) NOT NULL,
	"table_name" varchar(255) NOT NULL,
	"position" bigint NOT NULL,
	"icon" varchar(255),
	"filters" text,
	"sorting" text,
	"columns" text,
	"is_preset" boolean DEFAULT false NOT NULL,
	"is_shared" boolean DEFAULT false NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "ai_location_mapping" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"input" text NOT NULL,
	"response_json" text NOT NULL,
	"ai_prompt_log_id" uuid
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tenant_settings" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant_name" varchar(255) NOT NULL,
	"smart_sheet_id" varchar(255),
	"smart_sheet_access_token" varchar(255),
	"jira_api_token" varchar(255),
	"jira_domain" varchar(255),
	"jira_email" varchar(255),
	"trello_api_token" varchar(255),
	"trello_api_key" varchar(255),
	"aha_api_url" varchar(255),
	"aha_api_key" varchar(255),
	"airtable_personal_access_token" varchar(255),
	"amplitude_api_key" varchar(255),
	"amplitude_secret_key" varchar(255),
	"asana_access_token" varchar(255),
	"baton_api_key" varchar(255),
	"babelforce_region_environment" varchar(255),
	"babelforce_access_key_id" varchar(255),
	"babelforce_access_token" varchar(255),
	"bigquery_service_account_key" varchar(255),
	"braintree_public_key" varchar(255),
	"braintree_private_key" varchar(255),
	"braintree_environment" varchar(255),
	"braintree_merchant_id" varchar(255),
	"callrail_account" varchar(255),
	"callrail_api_token" varchar(255),
	"chargebee_api_key" varchar(255),
	"chargebee_product_catalog" varchar(255),
	"chargify_api_key" varchar(255),
	"chargify_domain" varchar(255),
	"clickup_api_key" varchar(255),
	"closecom_api_key" varchar(255),
	"coda_auth_token" varchar(255),
	"coda_document_id" varchar(255),
	"confluence_api_token" varchar(255),
	"confluence_domain" varchar(255),
	"confluence_login_email" varchar(255),
	"courier_api_key" varchar(255),
	"customerio_api_key" varchar(255),
	"datadog_api_key" varchar(255),
	"datadog_application_key" varchar(255),
	"delighted_api_key" varchar(255),
	"dixa_api_token" varchar(255),
	"drift_api_token" varchar(255),
	"emailoctopus_api_key" varchar(255),
	"facebook_marketing_access_token" varchar(255),
	"fastbill_api_key" varchar(255),
	"fastbill_project_id" varchar(255),
	"flexport_api_key" varchar(255),
	"freshcaller_api_key" varchar(255),
	"freshdesk_api_key" varchar(255),
	"freshdesk_domain" varchar(255),
	"freshsales_api_key" varchar(255),
	"freshsales_domain" varchar(255),
	"freshservice_api_key" varchar(255),
	"freshservice_domain" varchar(255),
	"genesys_region" varchar(255),
	"genesys_client_id" varchar(255),
	"genesys_client_secret" varchar(255),
	"github_access_token" varchar(255),
	"gitlab_access_token" varchar(255),
	"gocardless_access_token" varchar(255),
	"gocardless_environment" varchar(255),
	"gocardless_version" varchar(255),
	"gong_api_key" varchar(255),
	"harvest_account_id" varchar(255),
	"harvest_access_token" varchar(255),
	"insightly_api_token" varchar(255),
	"instagram_access_token" varchar(255),
	"instatus_api_key" varchar(255),
	"intercom_access_token" varchar(255),
	"klaviyo_api_key" varchar(255),
	"kustomer_api_token" varchar(255),
	"looker_client_id" varchar(255),
	"looker_client_secret" varchar(255),
	"looker_domain" varchar(255),
	"mailchimp_api_key" varchar(255),
	"mailjet_email_api_key" varchar(255),
	"mailjet_email_api_secret" varchar(255),
	"marketo_client_id" varchar(255),
	"marketo_client_secret" varchar(255),
	"marketo_domain_url" varchar(255),
	"microsoft_teams_tenant_id" varchar(255),
	"microsoft_teams_client_id" varchar(255),
	"microsoft_teams_client_secret" varchar(255),
	"monday_api_token" varchar(255),
	"notion_internal_access_token" varchar(255),
	"notion_public_access_token" varchar(255),
	"notion_public_client_id" varchar(255),
	"notion_public_client_secret" varchar(255),
	"oracle_netsuite_account_id" varchar(255),
	"oracle_netsuite_consumer_key" varchar(255),
	"oracle_netsuite_consumer_secret" varchar(255),
	"oracle_netsuite_token_id" varchar(255),
	"oracle_netsuite_token_secret" varchar(255),
	"orb_api_key" varchar(255),
	"orbit_api_key" varchar(255),
	"pager_duty_apikey" varchar(255),
	"paypal_transaction_client_id" varchar(255),
	"paypal_transaction_secret" varchar(255),
	"paystack_secret_key" varchar(255),
	"paystack_lookback_window" varchar(255),
	"pendo_api_token" varchar(255),
	"pipedrive_api_token" varchar(255),
	"plaid_access_token" varchar(255),
	"plausible_api_key" varchar(255),
	"plausible_site_id" varchar(255),
	"post_hog_api_key" varchar(255),
	"post_hog_base_url" varchar(255),
	"qualaroo_api_key" varchar(255),
	"qualaroo_api_token" varchar(255),
	"quick_books_client_id" varchar(255),
	"quick_books_client_secret" varchar(255),
	"quick_books_realm_id" varchar(255),
	"quick_books_refresh_token" varchar(255),
	"recharge_api_token" varchar(255),
	"recruitee_company_id" varchar(255),
	"recruitee_api_key" varchar(255),
	"recurly_api_key" varchar(255),
	"retently_api_token" varchar(255),
	"salesloft_api_key" varchar(255),
	"sendgrid_api_key" varchar(255),
	"sentry_project" varchar(255),
	"sentry_host" varchar(255),
	"sentry_authentication_token" varchar(255),
	"sentry_organization" varchar(255),
	"slack_api_token" varchar(255),
	"slack_channel_filter" varchar(255),
	"slack_lookback_window" varchar(255),
	"stripe_account_id" varchar(255),
	"stripe_secret_key" varchar(255),
	"survey_sparrow_access_token" varchar(255),
	"survey_monkey_access_token" varchar(255),
	"talkdesk_api_key" varchar(255),
	"tik_tok_access_token" varchar(255),
	"todoist_api_token" varchar(255),
	"typeform_api_token" varchar(255),
	"vittally_api_key" varchar(255),
	"wrike_access_token" varchar(255),
	"wrike_host_url" varchar(255),
	"xero_client_id" varchar(255),
	"xero_client_secret" varchar(255),
	"xero_tenant_id" varchar(255),
	"xero_scopes" varchar(255),
	"zendesk_api_key" varchar(255),
	"zendesk_subdomain" varchar(255),
	"zendesk_admin_email" varchar(255),
	"zendesk_chat_subdomain" varchar(255),
	"zendesk_chat_access_key" varchar(255),
	"zendesk_talk_subdomain" varchar(255),
	"zendesk_talk_access_key" varchar(255),
	"zendesk_sell_api_token" varchar(255),
	"zendesk_sunshine_subdomain" varchar(255),
	"zendesk_sunshine_api_token" varchar(255),
	"zendesk_sunshine_email" varchar(255),
	"zenefits_token" varchar(255),
	"mixpanel_username" varchar(255),
	"mixpanel_secret" varchar(255),
	"mixpanel_project_id" varchar(255),
	"mixpanel_project_secret" varchar(255),
	"mixpanel_project_timezone" varchar(255),
	"mixpanel_region" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "currency_rates" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"currency" varchar(3) NOT NULL,
	"rate" numeric NOT NULL,
	"date" timestamp with time zone NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"source" varchar(255) NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "ai_prompt_log" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"app_source" varchar(50) NOT NULL,
	"provider" varchar(50) NOT NULL,
	"model" varchar(100) NOT NULL,
	"prompt_type" varchar(255) NOT NULL,
	"prompt_template" text,
	"tenant" varchar(100),
	"node_id" varchar(255),
	"node_label" varchar(100),
	"prompt" text NOT NULL,
	"raw_response" text NOT NULL,
	"post_process_error" boolean,
	"post_process_error_message" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "slack_channel" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"source" varchar(255),
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"tenant_name" varchar(255) NOT NULL,
	"channel_id" varchar(255) NOT NULL,
	"channel_name" varchar(255),
	"organization_id" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "invoice_numbers" (
	"invoice_number" varchar(16) PRIMARY KEY NOT NULL,
	"tenant" varchar(50),
	"created_date" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"attempts" bigint
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tech_limit" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"key" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"usage_count" bigint NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "external_app_keys" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"app" varchar(255) NOT NULL,
	"app_key" varchar(255) NOT NULL,
	"group1" varchar(255),
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"usage_count" bigint NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "enrich_details_better_contact" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"request_id" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"contact_first_name" varchar(255),
	"contact_last_name" varchar(255),
	"contact_linkedin_url" varchar(255),
	"company_name" varchar(255),
	"company_domain" varchar(255),
	"request" text,
	"response" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "raw_email" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"sent_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"external_system" varchar(255) NOT NULL,
	"tenant" varchar(255) NOT NULL,
	"username" varchar(255) NOT NULL,
	"state" varchar(255) NOT NULL,
	"provider_message_id" varchar(255) NOT NULL,
	"message_id" varchar(255) NOT NULL,
	"sent_to_event_store_state" varchar(50) NOT NULL,
	"sent_to_event_store_reason" text,
	"sent_to_event_store_error" text,
	"data" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "slack_settings" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant_name" varchar(255) NOT NULL,
	"app_id" varchar(255),
	"authed_user_id" varchar(255),
	"scope" varchar(255),
	"token_type" varchar(255),
	"access_token" varchar(255),
	"bot_user_id" varchar(255),
	"team_id" varchar(255),
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "personal_integrations" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"tenant_name" varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	"email" varchar(255) NOT NULL,
	"key" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "postmark_api_keys" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"tenant_name" varchar(255) NOT NULL,
	"key" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "personal_email_provider" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"provider_name" varchar(255) NOT NULL,
	"provider_domain" varchar(255) NOT NULL,
	"created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "sync_run_webhook" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant" varchar(50),
	"external_system" varchar(50),
	"app_source" varchar(50),
	"entity" varchar(50),
	"start_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"end_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"reason" text,
	"total" bigint,
	"completed" bigint,
	"skipped" bigint,
	"failed" bigint
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "slack_channel_notification" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"tenant" varchar(255) NOT NULL,
	"channel_id" varchar(255) NOT NULL,
	"workflow" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "api_cache" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"tenant" varchar(100) NOT NULL,
	"type" varchar(255) NOT NULL,
	"data" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "workflow" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"tenant" varchar(255) NOT NULL,
	"workflow_type" varchar(255) NOT NULL,
	"name" varchar(255),
	"condition" text,
	"live" boolean DEFAULT false,
	"action_param1" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "industry_mapping" (
	"id" bigserial PRIMARY KEY NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"input_industry" varchar(255) NOT NULL,
	"output_industry" varchar(255) NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tracking" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"tenant" varchar(255),
	"user_id" varchar(255) NOT NULL,
	"ip" varchar(255),
	"event_type" varchar(255),
	"event_data" text,
	"timestamp" bigint,
	"href" varchar(255),
	"origin" varchar(255),
	"search" varchar(255),
	"hostname" varchar(255),
	"pathname" varchar(255),
	"referrer" varchar(255),
	"user_agent" text,
	"language" varchar(255),
	"cookies_enabled" boolean,
	"screen_resolution" varchar(255),
	"state" varchar(50),
	"organization_id" varchar(255),
	"organization_name" varchar(255),
	"notified" boolean DEFAULT false,
	"organization_domain" varchar(255),
	"organization_website" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "enrich_details_prefilter_tracking" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"ip" varchar(255),
	"should_identify" boolean,
	"response" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "enrich_details_tracking" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"ip" varchar(255),
	"company_name" varchar(255),
	"company_domain" varchar(255),
	"company_website" varchar(255),
	"response" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tenant" (
	"name" varchar(255) PRIMARY KEY NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tenant_settings_opportunity_stage" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"visible" boolean NOT NULL,
	"val" text NOT NULL,
	"idx" bigint NOT NULL,
	"label" varchar(255) NOT NULL,
	"likelihood_rate" bigint DEFAULT 0 NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tenant_settings_email_exclusion" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant" varchar(255) NOT NULL,
	"exclude_subject" varchar(255),
	"exclude_body" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "tenant_settings_mailbox" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp,
	"user_name" varchar(255),
	"mailbox_username" varchar(255),
	"mailbox_password" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "email_lookup" (
	"id" varchar(64) PRIMARY KEY NOT NULL,
	"tenant" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"tracker_domain" varchar(255),
	"message_id" varchar(64) NOT NULL,
	"link_id" varchar(64) NOT NULL,
	"redirect_url" varchar(255) NOT NULL,
	"campaign" varchar(255) NOT NULL,
	"type" varchar(32) NOT NULL,
	"recipient_id" varchar(255),
	"track_opens" boolean NOT NULL,
	"track_clicks" boolean NOT NULL,
	"unsubscribe_url" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "email_tracking" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"tenant" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"timestamp" timestamp DEFAULT CURRENT_TIMESTAMP,
	"message_id" varchar(64) NOT NULL,
	"link_id" varchar(64),
	"recipient_id" varchar(255),
	"campaign" varchar(255),
	"event_type" varchar(255) NOT NULL,
	"ip" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "flow_sequence_step_template_variable" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp NOT NULL,
	"updated_at" timestamp NOT NULL,
	"name" varchar(255) NOT NULL,
	"value" varchar(255) NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "flows" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp NOT NULL,
	"updated_at" timestamp NOT NULL,
	"tenant" text NOT NULL,
	"name" varchar(255) NOT NULL,
	"description" text,
	"active" boolean DEFAULT false NOT NULL,
	"active_days_string" varchar(255),
	"active_time_window_start" varchar(255),
	"active_time_window_end" varchar(255),
	"pause_on_holidays" boolean,
	"respect_recipient_timezone" boolean,
	"minutes_delay_between_emails" bigint,
	"emails_per_mailbox_per_hour" bigint,
	"emails_per_mailbox_per_day" bigint
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "flow_sequence" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp NOT NULL,
	"updated_at" timestamp NOT NULL,
	"flow_id" uuid NOT NULL,
	"name" varchar(255) NOT NULL,
	"description" text NOT NULL,
	"active" boolean DEFAULT false NOT NULL,
	"personas_string" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "flow_sequence_step" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp NOT NULL,
	"updated_at" timestamp NOT NULL,
	"sequence_id" uuid NOT NULL,
	"active" boolean DEFAULT false NOT NULL,
	"order" bigint NOT NULL,
	"type" varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	"text" varchar(255),
	"template" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "flow_sequence_contact" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp NOT NULL,
	"updated_at" timestamp NOT NULL,
	"sequence_id" uuid NOT NULL,
	"first_name" text,
	"last_name" text,
	"email" text NOT NULL,
	"linkedin_url" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "flow_sequence_sender" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp NOT NULL,
	"updated_at" timestamp NOT NULL,
	"sequence_id" uuid NOT NULL,
	"mailbox_id" text NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "cache_ip_data" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"ip" varchar(255) NOT NULL,
	"data" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "cache_ip_hunter" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"ip" varchar(255) NOT NULL,
	"data" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "cache_email_validation" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"email" varchar(255) NOT NULL,
	"is_deliverable" boolean,
	"is_mailbox_full" boolean,
	"is_role_account" boolean,
	"is_free_account" boolean,
	"smtp_success" boolean,
	"response_code" varchar(255),
	"error_code" varchar(255),
	"description" text,
	"retry_validation" boolean,
	"smtp_response" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "cache_email_validation_domain" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"domain" varchar(255) NOT NULL,
	"is_catch_all" boolean,
	"is_firewalled" boolean,
	"can_connect_smtp" boolean,
	"provider" varchar(255),
	"firewall" varchar(255)
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "scraper_configs" (
	"id" serial PRIMARY KEY NOT NULL,
	"user_id" varchar(36) NOT NULL,
	"tenant" text NOT NULL,
	"created_at" timestamp,
	"updated_at" timestamp
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "customer_os_ids" (
	"tenant" varchar(50) NOT NULL,
	"customer_os_id" varchar(30) NOT NULL,
	"entity" varchar(30),
	"entity_id" varchar(50),
	"created_date" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"attempts" bigint,
	CONSTRAINT "customer_os_ids_pkey" PRIMARY KEY("tenant","customer_os_id")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "oauth_token" (
	"provider" varchar(255) NOT NULL,
	"tenant_name" varchar(255) NOT NULL,
	"email_address" varchar(255) NOT NULL,
	"type" varchar(50),
	"player_identity_id" varchar(255) NOT NULL,
	"access_token" text,
	"refresh_token" text,
	"needs_manual_refresh" boolean DEFAULT false,
	"id_token" text,
	"expires_at" timestamp,
	"scope" text,
	"gmail_sync_enabled" boolean DEFAULT false,
	"google_calendar_sync_enabled" boolean DEFAULT false,
	CONSTRAINT "oauth_token_pkey" PRIMARY KEY("provider","tenant_name","email_address")
);
--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "idx_key" ON "app_keys" USING btree ("key");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "idx_tenant_api_keys" ON "google_service_account_keys" USING btree ("tenant_name","key");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "idx_enrich_details_scrapin_param1" ON "enrich_details_scrapin" USING btree ("param1");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "uq_one_state_per_tenant_and_user" ON "user_email_import_state" USING btree ("tenant","username","provider","state");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "name_domain_idx" ON "tracking_allowed_origin" USING btree ("tenant","origin");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "idx_key_unique" ON "tech_limit" USING btree ("key");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "idx_external_app_key_unique" ON "external_app_keys" USING btree ("app","app_key","group1");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "idx_raw_email_external_system" ON "raw_email" USING btree ("external_system","tenant","username","message_id");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "idx_tenant_uk" ON "slack_settings" USING btree ("tenant_name");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "idx_provider_domain" ON "personal_email_provider" USING btree ("provider_domain");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "ip_unique" ON "enrich_details_prefilter_tracking" USING btree ("ip");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "idx_cache_ip_data_ip" ON "cache_ip_data" USING btree ("ip");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "idx_cache_ip_hunter_ip" ON "cache_ip_hunter" USING btree ("ip");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "idx_cache_email_validation_email" ON "cache_email_validation" USING btree ("email");--> statement-breakpoint
CREATE UNIQUE INDEX IF NOT EXISTS "idx_cache_email_validation_domain_domain" ON "cache_email_validation_domain" USING btree ("domain");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "idx_primary" ON "oauth_token" USING btree ("provider","tenant_name","email_address");
*/