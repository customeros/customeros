import { eq } from "drizzle-orm";

import { ErrorParser } from "@/util/error";
import { logger } from "@/infrastructure/logger";

import { db } from "../db";
import { browserConfigs } from "../drizzle/schema";

export type BrowserConfigInsert = typeof browserConfigs.$inferInsert;
export type BrowserConfigTable = typeof browserConfigs.$inferSelect;
export type BrowserConfigSessionStatus = BrowserConfigTable["sessionStatus"];

export class BrowserConfigsRepository {
  constructor() {}

  async selectByUserId(userId: string) {
    try {
      const result = await db
        .select()
        .from(browserConfigs)
        .where(eq(browserConfigs.userId, userId));

      return result?.length > 0 ? result?.[0] : null;
    } catch (err) {
      BrowserConfigsRepository.handleError(err);
    }
  }

  async selectById(id: number) {
    try {
      const result = await db
        .select()
        .from(browserConfigs)
        .where(eq(browserConfigs.id, id));

      return result?.length > 0 ? result?.[0] : null;
    } catch (err) {
      BrowserConfigsRepository.handleError(err);
    }
  }

  async insert(values: BrowserConfigInsert) {
    try {
      const result = await db.insert(browserConfigs).values(values).returning();

      return result?.[0];
    } catch (err) {
      BrowserConfigsRepository.handleError(err);
    }
  }

  async updateByUserId(values: BrowserConfigInsert) {
    try {
      const result = await db
        .update(browserConfigs)
        .set(values)
        .where(eq(browserConfigs.userId, values.userId))
        .returning();

      return result?.[0];
    } catch (err) {
      BrowserConfigsRepository.handleError(err);
    }
  }

  private static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in Postgresql BrowserConfigsRepository", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
