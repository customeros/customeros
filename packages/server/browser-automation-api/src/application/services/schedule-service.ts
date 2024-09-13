import { ErrorParser } from "@/util/error";
import { Proxy } from "@/domain/models/proxy";
import { logger } from "@/infrastructure/logger";
import { Scheduler } from "@/infrastructure/scheduler";
import { BrowserAutomationRun } from "@/domain/models/browser-automation-run";
import {
  BrowserConfigsRepository,
  type BrowserAutomationRunType,
  BrowserAutomationRunsRepository,
} from "@/infrastructure/persistance/postgresql/repositories";
import { LinkedinService } from "./linkedin/linkedin-service";
import { BrowserConfig } from "@/domain/models/browser-config";
import { ProxyPoolRepository } from "@/infrastructure/persistance/postgresql/repositories/proxy-pool-repository";
import { AssignedProxiesRepository } from "@/infrastructure/persistance/postgresql/repositories/assigned-proxies-repository";
import { BrowserAutomationRunErrorsRepository } from "@/infrastructure/persistance/postgresql/repositories/browser-automation-run-errors-repository";
import { BrowserAutomationRunResultsRepository } from "@/infrastructure/persistance/postgresql/repositories/browser-automation-run-results-repository";

export class ScheduleService {
  private static instance: ScheduleService;
  public scheduler: Scheduler = Scheduler.getInstance();
  private browserConfigsRepository = new BrowserConfigsRepository();
  private browserAutomationRunsRepository =
    new BrowserAutomationRunsRepository();
  private browserAutomationRunResultsRepository =
    new BrowserAutomationRunResultsRepository();
  private browserAutomationRunErrorsRepository =
    new BrowserAutomationRunErrorsRepository();
  private assignedProxiesRepository = new AssignedProxiesRepository();
  private ProxyPoolRepository = new ProxyPoolRepository();

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
    payload?: Record<string, unknown>,
  ) {
    try {
      return await BrowserAutomationRun.create(
        {
          browserConfigId: browserConfig.id,
          tenant: browserConfig.tenant,
          userId: browserConfig.userId,
          type,
          payload: payload ? JSON.stringify(payload) : undefined,
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
        `Failed to find browser config with id: ${browserAutomationRun.browserConfigId}.`,
        {
          source: "ScheduleService",
        },
      );
      return;
    }

    const assignedProxy = await this.assignedProxiesRepository.selectByUserId(
      browserConfig?.userId,
    );

    if (!assignedProxy) {
      logger.warn(
        `Failed to find assigned proxy for user with id: ${browserConfig.userId}.`,
        {
          source: "ScheduleService",
        },
      );
      return;
    }

    const proxy = await this.ProxyPoolRepository.selectById(
      assignedProxy?.proxyPoolId,
    );

    if (!proxy) {
      logger.warn(
        `Failed to find proxy with id: ${assignedProxy.proxyPoolId}.`,
        {
          source: "ScheduleService",
        },
      );
      return;
    }

    const linkedinService = new LinkedinService(
      new BrowserConfig(browserConfig),
      Proxy.toBrowserHeader(proxy),
    );

    const jobParams = browserAutomationRun.toJobParams(
      linkedinService,
      this.browserAutomationRunsRepository,
      this.browserAutomationRunResultsRepository,
      this.browserAutomationRunErrorsRepository,
      this.browserConfigsRepository,
    );

    if (!jobParams) {
      logger.warn("Failed to create job params for automation run.", {
        source: "ScheduleService",
      });
      return;
    }

    this.scheduler.schedule(browserAutomationRun.id, jobParams);
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
      ScheduleService.handleError(err);
    } finally {
      completeTick();
    }
  }

  public async pollBrowserAutomationRuns() {
    logger.info("Browser automation poll-worker started.", {
      source: "ScheduleService",
    });
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
      ScheduleService.handleError(err);
    }
  }
}
