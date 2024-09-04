import { eq } from "drizzle-orm";

import { ErrorParser, StandardError } from "@/util/error";
import { logger } from "@/infrastructure/logger";

import { db } from "../db";
import { assignedProxies } from "../drizzle/schema";

export type AssignedProxyInsert = typeof assignedProxies.$inferInsert;
export type AssignedProxiesTable = typeof assignedProxies.$inferSelect;

export class AssignedProxiesRepository {
  constructor() {}

  async selectByUserId(userId: string) {
    try {
      const result = await db
        .select()
        .from(assignedProxies)
        .where(eq(assignedProxies.userId, userId));

      return result?.length > 0 ? result?.[0] : null;
    } catch (err) {
      AssignedProxiesRepository.handleError(err);
    }
  }

  async insert(values: AssignedProxyInsert) {
    try {
      const result = await db
        .insert(assignedProxies)
        .values(values)
        .returning();

      return result?.[0];
    } catch (err) {
      AssignedProxiesRepository.handleError(err);
    }
  }

  async updateByUserId(values: AssignedProxyInsert) {
    console.log(values);
    try {
      const result = await db
        .update(assignedProxies)
        .set(values)
        .where(eq(assignedProxies.userId, values.userId))
        .returning();

      return result?.[0];
    } catch (err) {
      AssignedProxiesRepository.handleError(err);
    }
  }

  private static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in Postgresql AssignedProxiesRepository", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
