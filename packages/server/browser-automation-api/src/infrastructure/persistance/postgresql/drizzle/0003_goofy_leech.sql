DO $$ BEGIN
 CREATE TYPE "public"."browser_config_session_status" AS ENUM('VALID', 'INVALID', 'EXPIRED');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
ALTER TABLE "browser_configs" ADD COLUMN "session_status" "browser_config_session_status" DEFAULT 'VALID' NOT NULL;