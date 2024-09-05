DO $$ BEGIN
 CREATE TYPE "public"."browser_automation_run_status" AS ENUM('SCHEDULED', 'RUNNING', 'COMPLETED', 'FAILED', 'CANCELLED');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 CREATE TYPE "public"."browser_automation_run_type" AS ENUM('FIND_CONNECTIONS', 'SEND_CONNECTION_REQUEST', 'SEND_MESSAGE');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 CREATE TYPE "public"."browser_automation_run_trigger" AS ENUM('MANUAL', 'SCHEDULER');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "assigned_proxies" (
	"id" serial PRIMARY KEY NOT NULL,
	"proxy_pool_id" integer NOT NULL,
	"user_id" varchar(36) NOT NULL,
	"tenant" text NOT NULL,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "browser_automation_run_errors" (
	"id" serial PRIMARY KEY NOT NULL,
	"run_id" integer NOT NULL,
	"occurred_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"error_type" varchar(100) NOT NULL,
	"error_message" text NOT NULL,
	"error_details" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "browser_automation_run_results" (
	"id" serial PRIMARY KEY NOT NULL,
	"run_id" integer NOT NULL,
	"type" varchar(50) NOT NULL,
	"result_data" text,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"is_processed" boolean DEFAULT false
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "browser_automation_runs" (
	"id" serial PRIMARY KEY NOT NULL,
	"browser_config_id" integer NOT NULL,
	"user_id" varchar(36) NOT NULL,
	"tenant" text NOT NULL,
	"type" "browser_automation_run_type" NOT NULL,
	"payload" text,
	"status" "browser_automation_run_status" DEFAULT 'SCHEDULED' NOT NULL,
	"scheduled_at" timestamp,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"started_at" timestamp,
	"finished_at" timestamp,
	"run_duration" integer,
	"retry_count" integer DEFAULT 0,
	"triggered_by" "browser_automation_run_trigger",
	"priority" integer DEFAULT 0,
	"log_location" text
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "browser_configs" (
	"id" serial PRIMARY KEY NOT NULL,
	"user_id" varchar(36) NOT NULL,
	"tenant" text NOT NULL,
	"cookies" text,
	"user_agent" text,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT "browser_configs_user_id_unique" UNIQUE("user_id")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "proxy_pool" (
	"id" serial PRIMARY KEY NOT NULL,
	"url" text NOT NULL,
	"username" text NOT NULL,
	"password" text NOT NULL,
	"enabled" boolean DEFAULT true,
	"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
	"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "browser_automation_run_errors" ADD CONSTRAINT "browser_automation_run_errors_run_id_browser_automation_runs_id_fk" FOREIGN KEY ("run_id") REFERENCES "public"."browser_automation_runs"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "browser_automation_run_results" ADD CONSTRAINT "browser_automation_run_results_run_id_browser_automation_runs_id_fk" FOREIGN KEY ("run_id") REFERENCES "public"."browser_automation_runs"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
