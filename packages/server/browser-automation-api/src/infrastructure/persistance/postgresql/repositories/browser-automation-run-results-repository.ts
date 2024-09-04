import { eq } from "drizzle-orm";

import { ErrorParser, StandardError } from "@/util/error";
import { logger } from "@/infrastructure/logger";

import { db } from "../db";
import { browserAutomationRunResults } from "../drizzle/schema";

export type BrowserAutomationRunResultsInsert =
  typeof browserAutomationRunResults.$inferInsert;
export type BrowserAutomationRunResultsTable =
  typeof browserAutomationRunResults.$inferSelect;

export class BrowserAutomationRunResultsRepository {
  constructor() {}

  async selectAllByRunId(runId: number) {
    try {
      const results = await db
        .select()
        .from(browserAutomationRunResults)
        .where(eq(browserAutomationRunResults.runId, runId));

      return results;
    } catch (err) {
      BrowserAutomationRunResultsRepository.handleError(err);
    }
  }

  async insert(values: BrowserAutomationRunResultsInsert) {
    try {
      const result = await db
        .insert(browserAutomationRunResults)
        .values(values)
        .returning();

      return result;
    } catch (err) {
      BrowserAutomationRunResultsRepository.handleError(err);
    }
  }

  async updateById(values: BrowserAutomationRunResultsInsert) {
    if (!values.id) {
      throw new StandardError({
        code: "INVARIANT_ERROR",
        message:
          "Invariant Violation: BrowserAutomationRunResult.id is required",
        severity: "critical",
        details: "BrowserAutomationRunResult ID is required to update",
      });
    }

    try {
      const result = await db
        .update(browserAutomationRunResults)
        .set(values)
        .where(eq(browserAutomationRunResults.id, values.id))
        .returning();

      return result;
    } catch (err) {
      BrowserAutomationRunResultsRepository.handleError(err);
    }
  }

  private static handleError(err: any) {
    const error = ErrorParser.parse(err);
    logger.error("Error in Postgresql BrowserAutomationRunResultsRepository", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
