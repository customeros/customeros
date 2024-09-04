import { eq } from "drizzle-orm";

import { ErrorParser, StandardError } from "@/util/error";
import { logger } from "@/infrastructure/logger";

import { db } from "../db";
import { proxyPool } from "../drizzle/schema";

export type ProxyPoolInsert = typeof proxyPool.$inferInsert;
export type ProxyPoolTable = typeof proxyPool.$inferSelect;

export class ProxyPoolRepository {
  constructor() {}

  async selectAll() {
    try {
      const result = await db.select().from(proxyPool);

      return result;
    } catch (err) {
      ProxyPoolRepository.handleError(err);
    }
  }

  async selectById(id: number) {
    try {
      const result = await db
        .select()
        .from(proxyPool)
        .where(eq(proxyPool.id, id));

      return result?.length > 0 ? result?.[0] : null;
    } catch (err) {
      ProxyPoolRepository.handleError(err);
    }
  }

  async insert(values: ProxyPoolInsert) {
    try {
      const result = await db.insert(proxyPool).values(values).returning();

      return result?.[0];
    } catch (err) {
      ProxyPoolRepository.handleError(err);
    }
  }

  async updateById(values: ProxyPoolInsert) {
    if (!values?.id) {
      throw new StandardError({
        code: "INTERNAL_ERROR",
        message: "ProxyPoolRepository.updateById: id is required",
        severity: "high",
      });
    }

    try {
      const result = await db
        .update(proxyPool)
        .set(values)
        .where(eq(proxyPool.id, values.id))
        .returning();

      return result?.[0];
    } catch (err) {
      ProxyPoolRepository.handleError(err);
    }
  }

  private static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in Postgresql ProxyPoolRepository", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
