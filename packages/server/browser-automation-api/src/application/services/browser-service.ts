import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import {
  BrowserConfig,
  BrowserConfigPayload,
} from "@/domain/models/browser-config";
import { BrowserAutomationRun } from "@/domain/models/browser-automation-run";
import { BrowserConfigsRepository } from "@/infrastructure/persistance/postgresql/repositories/browser-configs-repository";
import { BrowserAutomationRunsRepository } from "@/infrastructure/persistance/postgresql/repositories/browser-automation-runs-repository";

export class BrowserService {
  private browserConfigsRepository: BrowserConfigsRepository;
  private browserAutomationRunRepository: BrowserAutomationRunsRepository;

  constructor() {
    this.browserConfigsRepository = new BrowserConfigsRepository();
    this.browserAutomationRunRepository = new BrowserAutomationRunsRepository();
  }

  async getBrowserConfig(
    userId: string,
  ): Promise<BrowserConfig | null | undefined> {
    try {
      const values = await this.browserConfigsRepository.selectByUserId(userId);
      if (!values) {
        return null;
      }
      return new BrowserConfig(values);
    } catch (err) {
      BrowserService.handleError(err);
    }
  }

  async createBrowserConfig(payload: BrowserConfigPayload) {
    try {
      return await BrowserConfig.create(payload, this.browserConfigsRepository);
    } catch (err) {
      BrowserService.handleError(err);
    }
  }

  async updateBrowserConfig(payload: BrowserConfigPayload) {
    try {
      return await BrowserConfig.update(payload, this.browserConfigsRepository);
    } catch (err) {
      BrowserService.handleError(err);
    }
  }

  async getBrowserAutomationRuns(
    userId: string,
  ): Promise<BrowserAutomationRun[] | undefined> {
    try {
      const values =
        await this.browserAutomationRunRepository.selectAllByUserId(userId);
      if (!values) {
        return [];
      }
      return values.map((value) => new BrowserAutomationRun(value));
    } catch (err) {
      BrowserService.handleError(err);
    }
  }

  async getBrowserAutomationRun(
    id: number,
  ): Promise<BrowserAutomationRun | undefined> {
    try {
      const value = await this.browserAutomationRunRepository.selectById(id);
      if (!value) {
        return undefined;
      }
      return new BrowserAutomationRun(value);
    } catch (err) {
      BrowserService.handleError(err);
    }
  }

  private static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in BrowserService", {
      message: error.message,
      details: error.details,
    });
    throw error;
  }
}
