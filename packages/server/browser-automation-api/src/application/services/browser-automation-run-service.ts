import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { BrowserAutomationRunsRepository } from "@/infrastructure/persistance/postgresql/repositories";
import {
  BrowserAutomationRun,
  type BrowserAutomationRunPayload,
} from "@/domain/models/browser-automation-run";

export class BrowserAutomationRunService {
  constructor(
    private browserAutomationRunsRepository: BrowserAutomationRunsRepository,
  ) {}

  async createRun(payload: BrowserAutomationRunPayload) {
    try {
      const newRun = await this.browserAutomationRunsRepository.insert(payload);

      if (newRun) {
        return new BrowserAutomationRun(newRun);
      }
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error creating automation run", {
        source: "BrowserAutomationRunService",
      });
      throw error;
    }
  }
}
