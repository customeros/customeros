import { eq } from "drizzle-orm";

import { ErrorParser, StandardError } from "@/util/error";
import { logger } from "@/infrastructure/logger";

import { db } from "../db";
import { browserAutomationRunErrors } from "../drizzle/schema";

export type BrowserAutomationRunErrorsInsert =
  typeof browserAutomationRunErrors.$inferInsert;
export type BrowserAutomationRunErrorsTable =
  typeof browserAutomationRunErrors.$inferSelect;

export class BrowserAutomationRunErrorsRepository {
  constructor() {}

  async selectAllByRunId(runId: number) {
    try {
      const results = await db
        .select()
        .from(browserAutomationRunErrors)
        .where(eq(browserAutomationRunErrors.runId, runId));

      return results;
    } catch (err) {
      BrowserAutomationRunErrorsRepository.handleError(err);
    }
  }

  async insert(values: BrowserAutomationRunErrorsInsert) {
    try {
      const result = await db
        .insert(browserAutomationRunErrors)
        .values(values)
        .returning();

      return result;
    } catch (err) {
      BrowserAutomationRunErrorsRepository.handleError(err);
    }
  }

  async updateById(values: BrowserAutomationRunErrorsInsert) {
    if (!values.id) {
      throw new StandardError({
        code: "INVARIANT_ERROR",
        message:
          "Invariant Violation: BrowserAutomationRunErrors.id is required",
        severity: "critical",
        details: "BrowserAutomationRunErrors ID is required to update",
      });
    }

    try {
      const result = await db
        .update(browserAutomationRunErrors)
        .set(values)
        .where(eq(browserAutomationRunErrors.id, values.id))
        .returning();

      return result;
    } catch (err) {
      BrowserAutomationRunErrorsRepository.handleError(err);
    }
  }

  private static handleError(err: any) {
    const error = ErrorParser.parse(err);
    logger.error("Error in Postgresql BrowserAutomationRunErrorsRepository", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
