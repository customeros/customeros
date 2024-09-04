import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { type JobParams } from "@/infrastructure/scheduler";
import { type LinkedinService } from "@/application/services/linkedin/linkedin-service";
import {
  type BrowserAutomationRunType,
  type BrowserAutomationRunTable,
  type BrowserAutomationRunInsert,
  type BrowserAutomationRunStatus,
  type BrowserAutomationRunTrigger,
  BrowserAutomationRunsRepository,
} from "@/infrastructure/persistance/postgresql/repositories/browser-automation-runs-repository";
import { BrowserAutomationRunErrorsRepository } from "@/infrastructure/persistance/postgresql/repositories/browser-automation-run-errors-repository";
import { type BrowserAutomationRunResultsRepository } from "@/infrastructure/persistance/postgresql/repositories/browser-automation-run-results-repository";

export type BrowserAutomationRunPayload = Pick<
  BrowserAutomationRunInsert,
  | "type"
  | "userId"
  | "tenant"
  | "payload"
  | "priority"
  | "triggeredBy"
  | "browserConfigId"
>;

export class BrowserAutomationRun {
  id: number;
  browserConfigId: number;
  userId: string;
  tenant: string;
  type: BrowserAutomationRunType;
  payload: string | null;
  status: BrowserAutomationRunStatus;
  scheduledAt: string | null;
  createdAt: string | null;
  updatedAt: string | null;
  startedAt: string | null;
  finishedAt: string | null;
  runDuration: number | null;
  retryCount: number | null;
  triggeredBy: BrowserAutomationRunTrigger;
  priority: number | null;
  logLocation: string | null;

  constructor(values: BrowserAutomationRunTable) {
    this.id = values.id;
    this.browserConfigId = values.browserConfigId;
    this.userId = values.userId;
    this.tenant = values.tenant;
    this.type = values.type;
    this.payload = values.payload;
    this.status = values.status;
    this.scheduledAt = values.scheduledAt;
    this.createdAt = values.createdAt;
    this.updatedAt = values.updatedAt;
    this.startedAt = values.startedAt;
    this.finishedAt = values.finishedAt;
    this.runDuration = values.runDuration;
    this.retryCount = values.retryCount;
    this.triggeredBy = values.triggeredBy;
    this.priority = values.priority;
    this.logLocation = values.logLocation;
  }

  toDTO(): BrowserAutomationRunTable {
    return {
      id: this.id,
      browserConfigId: this.browserConfigId,
      userId: this.userId,
      tenant: this.tenant,
      type: this.type,
      payload: this.payload,
      status: this.status,
      scheduledAt: this.scheduledAt,
      createdAt: this.createdAt,
      updatedAt: this.updatedAt,
      startedAt: this.startedAt,
      finishedAt: this.finishedAt,
      runDuration: this.runDuration,
      retryCount: this.retryCount,
      triggeredBy: this.triggeredBy,
      priority: this.priority,
      logLocation: this.logLocation,
    };
  }

  async updateStatus(
    status: BrowserAutomationRunStatus,
    repository: BrowserAutomationRunsRepository,
  ) {
    try {
      this.status = status;
      await repository.updateById(this.toDTO());
    } catch (err) {
      BrowserAutomationRun.handleError(err);
    }
  }

  toJobParams(
    linkedinService: LinkedinService,
    automationRepository: BrowserAutomationRunsRepository,
    resultsRepository: BrowserAutomationRunResultsRepository,
    errorsRepository: BrowserAutomationRunErrorsRepository,
  ): JobParams | null {
    const payload = BrowserAutomationRun.parsePayload(this.payload);

    switch (this.type) {
      case "SEND_MESSAGE": {
        if (!payload) return null;

        return {
          cronTime: "* * * * * *",
          start: true,
          runOnce: true,
          onTick: async () => {
            await linkedinService.sendMessage(
              payload?.profileUrl,
              payload?.message,
              {
                dryRun: payload?.dryRun,
                onStart: async () => {
                  await this.updateStatus("RUNNING", automationRepository);
                  logger.info("Sending message", {
                    profileUrl: payload?.profileUrl,
                  });
                },
                onSuccess: async () => {
                  await this.updateStatus("COMPLETED", automationRepository);
                  await resultsRepository.insert({
                    runId: this.id,
                    type: "SEND_MESSAGE",
                    resultData: JSON.stringify({
                      profileUrl: payload?.profileUrl,
                      message: "Message sent successfully",
                    }),
                  });
                  logger.info("Message sent", {
                    profileUrl: payload?.profileUrl,
                  });
                },
                onError: async (err) => {
                  await this.updateStatus("FAILED", automationRepository);
                  await errorsRepository.insert({
                    runId: this.id,
                    errorMessage: err.message,
                    errorDetails: err.details,
                    errorType: err.code,
                  });
                  logger.error("Failed to send message", {
                    profileUrl: payload?.profileUrl,
                  });
                },
              },
            );
          },
        };
      }
      case "SEND_CONNECTION_REQUEST": {
        if (!payload) return null;

        return {
          cronTime: "* * * * * *",
          start: true,
          runOnce: true,
          onTick: async () =>
            await linkedinService.sendInvite(
              payload?.profileUrl,
              payload?.message,
              {
                dryRun: payload?.dryRun,
                onStart: async () => {
                  this.updateStatus("RUNNING", automationRepository);
                  logger.info("Sending connection invite", {
                    profileUrl: payload?.profileUrl,
                  });
                },
                onSuccess: async () => {
                  this.updateStatus("COMPLETED", automationRepository);
                  await resultsRepository.insert({
                    runId: this.id,
                    type: "SEND_CONNECTION_REQUEST",
                    resultData: JSON.stringify({
                      profileUrl: payload?.profileUrl,
                      message: "Connection invite sent successfully",
                    }),
                  });
                  logger.info("Connection request sent", {
                    profileUrl: payload?.profileUrl,
                  });
                },
                onError: async (err) => {
                  this.updateStatus("FAILED", automationRepository);
                  await errorsRepository.insert({
                    runId: this.id,
                    errorMessage: err.message,
                    errorDetails: err.details,
                    errorType: err.code,
                  });
                  logger.error("Failed to send connection invite", {
                    profileUrl: payload?.profileUrl,
                  });
                },
              },
            ),
        };
      }
      case "FIND_CONNECTIONS": {
        if (!payload) return null;

        return {
          cronTime: "0 0 0 * * *",
          start: true,
          runOnce: true,
          onTick: async () => {
            await linkedinService.scrapeConnections({
              dryRun: payload?.dryRun,
              onStart: async () => {
                await this.updateStatus("RUNNING", automationRepository);
                logger.info("Scraping connections");
              },
              onSuccess: async () => {
                this.updateStatus("COMPLETED", automationRepository);
                logger.info("Connections scraped");
              },
              onError: async () => {
                this.updateStatus("FAILED", automationRepository);
                logger.error("Failed to scrape connections");
              },
            });
          },
        };
      }
      default:
        return null;
    }
  }

  static async create(
    values: BrowserAutomationRunPayload,
    browserAutomationRunsRepository: BrowserAutomationRunsRepository,
  ) {
    try {
      return browserAutomationRunsRepository.insert(values);
    } catch (err) {
      BrowserAutomationRun.handleError(err);
    }
  }

  private static parsePayload(payload: string | null) {
    try {
      if (!payload) return null;
      return JSON.parse(payload);
    } catch (err) {
      BrowserAutomationRun.handleError(err);
    }
  }

  static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in BrowserConfig", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
