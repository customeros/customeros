import { logger } from "@/infrastructure";
import { ErrorParser, StandardError } from "@/util/error";

import {
  type BrowserAutomationRunResultsInsert,
  type BrowserAutomationRunResultsTable,
  BrowserAutomationRunResultsRepository,
} from "@/infrastructure/persistance/postgresql/repositories/browser-automation-run-results-repository";

export class BrowserAutomationRunResult {
  id: number;
  runId: number;
  type: string;
  resultData: string | null;
  createdAt: string | null;
  isProcessed: boolean | null;

  constructor(values: BrowserAutomationRunResultsTable) {
    this.id = values.id;
    this.runId = values.runId;
    this.type = values.type;
    this.resultData = values.resultData;
    this.createdAt = values.createdAt;
    this.isProcessed = values.isProcessed;
  }

  static async create(
    values: BrowserAutomationRunResultsInsert,
    repository: BrowserAutomationRunResultsRepository,
  ): Promise<BrowserAutomationRunResult | undefined> {
    try {
      const result = await repository.insert(values);
      if (!result) {
        throw new StandardError({
          code: "INTERNAL_ERROR",
          severity: "critical",
          message: "Error creating BrowserAutomationRunResult",
        });
      }

      return new BrowserAutomationRunResult(result?.[0]);
    } catch (err) {
      BrowserAutomationRunResult.handleError(err);
    }
  }

  private static handleError(err: any) {
    const error = ErrorParser.parse(err);
    logger.error("Error creating BrowserAutomationRunResult", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
