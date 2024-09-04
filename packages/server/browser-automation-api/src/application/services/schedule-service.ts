import { ErrorParser } from "@/util/error";
import { logger } from "@/infrastructure/logger";
import { Scheduler } from "@/infrastructure/scheduler";
import { BrowserAutomationRun } from "@/domain/models/browser-automation-run";
import {
  BrowserConfigsRepository,
  type BrowserAutomationRunType,
  BrowserAutomationRunsRepository,
} from "@/infrastructure/persistance/postgresql/repositories";
import { BrowserConfig } from "@/domain/models/browser-config";
import { BrowserAutomationRunErrorsRepository } from "@/infrastructure/persistance/postgresql/repositories/browser-automation-run-errors-repository";
import { BrowserAutomationRunResultsRepository } from "@/infrastructure/persistance/postgresql/repositories/browser-automation-run-results-repository";

import { LinkedinService } from "./linkedin/linkedin-service";

export class ScheduleService {
  private static instance: ScheduleService;
  private scheduler: Scheduler = Scheduler.getInstance();
  private browserConfigsRepository = new BrowserConfigsRepository();
  private browserAutomationRunsRepository =
    new BrowserAutomationRunsRepository();
  private browserAutomationRunResultsRepository =
    new BrowserAutomationRunResultsRepository();
  private browserAutomationRunErrorsRepository =
    new BrowserAutomationRunErrorsRepository();

  constructor() {
    this.getUnscheduledRuns = this.getUnscheduledRuns.bind(this);
    this.scheduleAutomationRun = this.scheduleAutomationRun.bind(this);
    this.pollBrowserAutomationRuns = this.pollBrowserAutomationRuns.bind(this);
  }

  public static getInstance() {
    if (!ScheduleService.instance) {
      ScheduleService.instance = new ScheduleService();
    }

    return ScheduleService.instance;
  }

  public async createAutomationRun(
    browserConfig: BrowserConfig,
    type: BrowserAutomationRunType,
    payload: Record<string, unknown>,
  ) {
    try {
      return await BrowserAutomationRun.create(
        {
          browserConfigId: browserConfig.id,
          tenant: browserConfig.tenant,
          userId: browserConfig.userId,
          type,
          payload: JSON.stringify(payload),
        },
        this.browserAutomationRunsRepository,
      );
    } catch (err) {
      ScheduleService.handleError(err);
    }
  }

  private async scheduleAutomationRun(
    browserAutomationRun: BrowserAutomationRun,
  ) {
    const browserConfig = await this.browserConfigsRepository.selectById(
      browserAutomationRun.browserConfigId,
    );

    if (!browserConfig) {
      logger.warn(
        `Failed to find browser config with id: ${browserAutomationRun.browserConfigId}`,
      );
      return;
    }

    const linkedinService = new LinkedinService(
      new BrowserConfig(browserConfig),
    );

    const jobParams = browserAutomationRun.toJobParams(
      linkedinService,
      this.browserAutomationRunsRepository,
      this.browserAutomationRunResultsRepository,
      this.browserAutomationRunErrorsRepository,
    );

    if (!jobParams) {
      logger.warn("Failed to create job params for automation run.");
      return;
    }

    this.scheduler.schedule(browserAutomationRun.id, jobParams);
  }

  private async getUnscheduledRuns() {
    try {
      const browserAutomationRuns =
        await this.browserAutomationRunsRepository.selectAllScheduled();

      if (!browserAutomationRuns?.length) {
        return;
      }

      logger.info("Found scheduled runs. Queuing them for execution.");

      browserAutomationRuns?.forEach((browserAutomationRun) => {
        const run = new BrowserAutomationRun(browserAutomationRun);
        this.scheduleAutomationRun(run);
      });
    } catch (err) {
      ScheduleService.handleError(err);
    }
  }

  public async pollBrowserAutomationRuns() {
    logger.info("Browser automation poll-worker started.");
    this.scheduler.schedule("poll-worker", {
      cronTime: "*/20 * * * * *",
      onTick: this.getUnscheduledRuns,
      start: true,
    });
  }

  private static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in ScheduleService", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
