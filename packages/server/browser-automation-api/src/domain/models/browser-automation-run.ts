import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import {
  type BrowserAutomationRunType,
  type BrowserAutomationRunTable,
  type BrowserAutomationRunInsert,
  type BrowserAutomationRunStatus,
  type BrowserAutomationRunTrigger,
  BrowserAutomationRunsRepository,
} from "@/infrastructure/persistance/postgresql/repositories";

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

  start() {
    this.status = "RUNNING";
    this.startedAt = new Date().toISOString();
  }

  complete() {
    this.status = "COMPLETED";
    this.finishedAt = new Date().toISOString();
    this.runDuration = this.getDuration();
  }

  fail() {
    this.status = "FAILED";
    this.finishedAt = new Date().toISOString();
    this.runDuration = this.getDuration();
  }

  private getDuration(): number | null {
    if (this.startedAt && this.finishedAt) {
      return (
        new Date(this.finishedAt).getTime() - new Date(this.startedAt).getTime()
      );
    }

    return null;
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

  static parsePayload(payload: string | null) {
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
