import { eq } from "drizzle-orm";

import { logger } from "@/infrastructure/logger";
import { ErrorParser } from "@/util/error";

import { db } from "../db";
import { tenantWebhookApiKeys } from "../drizzle/schema";

export class TenantWebhookApiKeysRepository {
  constructor() {
    this.selectByApiKey = this.selectByApiKey.bind(this);
  }

  async selectByApiKey(key: string) {
    try {
      const result = await db
        .select()
        .from(tenantWebhookApiKeys)
        .where(eq(tenantWebhookApiKeys.key, key));

      return result?.length > 0 ? result?.[0] : null;
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in Postgresql TenantWebhookApiKeysRepository", {
        error: error.message,
        details: error.details,
      });
      throw error;
    }
  }
}
