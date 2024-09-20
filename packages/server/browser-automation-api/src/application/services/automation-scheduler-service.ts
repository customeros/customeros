import { ErrorParser } from "@/util/error";
import { logger } from "@/infrastructure/logger";
import { Scheduler } from "@/infrastructure/scheduler";
import { BrowserAutomationRun } from "@/domain/models/browser-automation-run";
import {
  ProxyPoolRepository,
  BrowserConfigsRepository,
  AssignedProxiesRepository,
  BrowserAutomationRunsRepository,
  BrowserAutomationRunErrorsRepository,
  BrowserAutomationRunResultsRepository,
} from "@/infrastructure/persistance/postgresql/repositories";

import { AutomationRunnerService } from "./automation-runner-service";
import { LinkedinServiceFactory } from "./linkedin/linkedin-service-factory";

export class AutomationSchedulerService {
  private static instance: AutomationSchedulerService;
  public scheduler: Scheduler = Scheduler.getInstance();
  private linkedinServiceFactory: LinkedinServiceFactory;
  private automationRunnerService: AutomationRunnerService;
  private browserAutomationRunsRepository: BrowserAutomationRunsRepository;

  constructor() {
    this.browserAutomationRunsRepository =
      new BrowserAutomationRunsRepository();
    this.linkedinServiceFactory = new LinkedinServiceFactory(
      new BrowserConfigsRepository(),
      new ProxyPoolRepository(),
      new AssignedProxiesRepository(),
    );
    this.automationRunnerService = new AutomationRunnerService(
      this.linkedinServiceFactory,
      this.browserAutomationRunsRepository,
      new BrowserAutomationRunResultsRepository(),
      new BrowserAutomationRunErrorsRepository(),
      new BrowserConfigsRepository(),
    );

    this.getRetryingRuns = this.getRetryingRuns.bind(this);
    this.getUnscheduledRuns = this.getUnscheduledRuns.bind(this);
    this.scheduleAutomationRun = this.scheduleAutomationRun.bind(this);
    this.pollBrowserAutomationRuns = this.pollBrowserAutomationRuns.bind(this);
  }

  public static getInstance() {
    if (!AutomationSchedulerService.instance) {
      AutomationSchedulerService.instance = new AutomationSchedulerService();
    }

    return AutomationSchedulerService.instance;
  }

  private async scheduleAutomationRun(
    browserAutomationRun: BrowserAutomationRun,
  ) {
    this.scheduler.schedule(browserAutomationRun.id, {
      cronTime: "* * * * * *",
      start: true,
      runOnce: true,
      onTick: async (completeTick) => {
        await this.automationRunnerService.runAutomation(browserAutomationRun);
        completeTick();
      },
    });
  }

  private async getUnscheduledRuns(completeTick: () => void) {
    try {
      const browserAutomationRuns =
        await this.browserAutomationRunsRepository.selectAllScheduled();

      if (!browserAutomationRuns?.length) {
        return;
      }

      logger.info("Found scheduled runs, queuing them for execution.", {
        source: "ScheduleService",
      });

      browserAutomationRuns?.forEach((browserAutomationRun) => {
        const run = new BrowserAutomationRun(browserAutomationRun);
        this.scheduleAutomationRun(run);
      });
    } catch (err) {
      AutomationSchedulerService.handleError(err);
    } finally {
      completeTick();
    }
  }

  private async getRetryingRuns(completeTick: () => void) {
    try {
      const browserAutomationRuns =
        await this.browserAutomationRunsRepository.selectAllRetrying();

      if (!browserAutomationRuns?.length) {
        return;
      }

      logger.info(
        "Found retrying runs, queuing them for continuing execution.",
        {
          source: "ScheduleService",
        },
      );

      browserAutomationRuns?.forEach((browserAutomationRun) => {
        const run = new BrowserAutomationRun(browserAutomationRun);
        this.scheduleAutomationRun(run);
      });
    } catch (err) {
      AutomationSchedulerService.handleError(err);
    } finally {
      completeTick();
    }
  }

  public async pollBrowserAutomationRuns() {
    logger.info("Browser automation pollers started.", {
      source: "ScheduleService",
    });
    this.scheduler.schedule("poll-unscheduled-worker", {
      cronTime: "*/20 * * * * *",
      onTick: this.getUnscheduledRuns,
      start: true,
    });
    this.scheduler.schedule("poll-retrying-worker", {
      cronTime: "*/10 * * * * *",
      onTick: this.getRetryingRuns,
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

  public async shutdown() {
    try {
      const isStopped = await this.scheduler.stopJobs();
      if (isStopped) {
        logger.info("All jobs have been stopped", {
          source: "ScheduleService",
        });
        return true;
      }
    } catch (err) {
      AutomationSchedulerService.handleError(err);
    }
  }
}
