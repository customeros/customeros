import { eq } from "drizzle-orm";

import { ErrorParser, StandardError } from "@/util/error";
import { logger } from "@/infrastructure/logger";

import { db } from "../db";
import { browserAutomationRuns } from "../drizzle/schema";

export type BrowserAutomationRunInsert =
  typeof browserAutomationRuns.$inferInsert;
export type BrowserAutomationRunTable =
  typeof browserAutomationRuns.$inferSelect;
export type BrowserAutomationRunStatus = BrowserAutomationRunTable["status"];
export type BrowserAutomationRunTrigger =
  BrowserAutomationRunTable["triggeredBy"];
export type BrowserAutomationRunType = BrowserAutomationRunTable["type"];

export class BrowserAutomationRunsRepository {
  constructor() {}

  async selectAllByUserId(userId: string) {
    try {
      const results = await db
        .select()
        .from(browserAutomationRuns)
        .where(eq(browserAutomationRuns.userId, userId));

      return results;
    } catch (err) {
      BrowserAutomationRunsRepository.handleError(err);
    }
  }

  async selectById(id: number) {
    try {
      const results = await db
        .select()
        .from(browserAutomationRuns)
        .where(eq(browserAutomationRuns.id, id));

      return results?.[0];
    } catch (err) {
      BrowserAutomationRunsRepository.handleError(err);
    }
  }

  async selectAllScheduled() {
    try {
      const results = await db
        .select()
        .from(browserAutomationRuns)
        .where(eq(browserAutomationRuns.status, "SCHEDULED"));

      return results;
    } catch (err) {
      BrowserAutomationRunsRepository.handleError(err);
    }
  }

  async insert(values: BrowserAutomationRunInsert) {
    try {
      const result = await db
        .insert(browserAutomationRuns)
        .values(values)
        .returning();

      return result?.[0];
    } catch (err) {
      BrowserAutomationRunsRepository.handleError(err);
    }
  }

  async updateById(values: BrowserAutomationRunInsert) {
    if (!values.id) {
      throw new StandardError({
        code: "INVARIANT_ERROR",
        message: "Invariant Violation: BrowserAutomationRun.id is required",
        severity: "critical",
        details: "BrowserAutomationRun ID is required to update",
      });
    }

    try {
      const result = await db
        .update(browserAutomationRuns)
        .set(values)
        .where(eq(browserAutomationRuns.id, values.id))
        .returning();

      return result?.[0];
    } catch (err) {
      BrowserAutomationRunsRepository.handleError(err);
    }
  }

  private static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in Postgresql BrowserAutomationRunsRepository", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
