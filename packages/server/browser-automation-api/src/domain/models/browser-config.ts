import { logger } from "@/infrastructure";
import { ErrorParser, StandardError } from "@/util/error";
import {
  BrowserConfigTable,
  BrowserConfigInsert,
  BrowserConfigsRepository,
  BrowserConfigSessionStatus,
} from "@/infrastructure/persistance/postgresql/repositories/browser-configs-repository";

export type BrowserConfigPayload = Pick<
  BrowserConfigInsert,
  "cookies" | "userAgent" | "tenant" | "userId"
>;

export class BrowserConfig {
  id: number;
  userId: string;
  tenant: string;
  cookies: string | null = null;
  userAgent: string | null = null;
  createdAt: string | null = null;
  updatedAt: string | null = null;
  sessionStatus: BrowserConfigSessionStatus;

  constructor(values: BrowserConfigTable) {
    this.id = values.id;
    this.userId = values.userId;
    this.tenant = values.tenant;
    this.cookies = values.cookies;
    this.userAgent = values.userAgent;
    this.createdAt = values.createdAt;
    this.updatedAt = values.updatedAt;
    this.sessionStatus = values.sessionStatus;
  }

  static async create(
    values: BrowserConfigPayload,
    browserConfigsRepository: BrowserConfigsRepository,
  ) {
    try {
      const existingBrowserConfig =
        await browserConfigsRepository.selectByUserId(values.userId);

      if (existingBrowserConfig) {
        throw new StandardError({
          code: "INVARIANT_ERROR",
          message: "Browser config already exists",
          severity: "medium",
          details: "Invariant violation: User already has a browser config.",
        });
      }

      return browserConfigsRepository.insert(values);
    } catch (err) {
      BrowserConfig.handleError(err);
    }
  }

  static async update(
    values: BrowserConfigPayload,
    browserConfigsRepository: BrowserConfigsRepository,
  ) {
    try {
      const existingBrowserConfig =
        await browserConfigsRepository.selectByUserId(values.userId);

      if (!existingBrowserConfig) {
        throw new StandardError({
          code: "INVARIANT_ERROR",
          message: "Browser config does not exists",
          severity: "medium",
          details: "Invariant violation: User does not have a browser config.",
        });
      }

      return await browserConfigsRepository.updateByUserId(values);
    } catch (err) {
      BrowserConfig.handleError(err);
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
