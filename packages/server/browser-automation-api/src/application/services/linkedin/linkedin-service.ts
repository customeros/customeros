import { logger } from "@/infrastructure";
import { Proxy } from "@/domain/models/proxy";
import { ErrorParser, StandardError } from "@/util/error";
import { BrowserConfig } from "@/domain/models/browser-config";
import { LinkedinAutomationService } from "@/infrastructure/scraper/services/linkedin-automation-service";

type LinkedinServiceMethodOptions = {
  dryRun?: boolean;
  onStart?: () => void;
  onSuccess?: () => void;
  onError?: (err: StandardError) => void;
};

export class LinkedinService {
  private linkedinAutomationService: LinkedinAutomationService;

  constructor(
    private browserConfig: BrowserConfig,
    private proxyHeader: string,
  ) {
    this.linkedinAutomationService = new LinkedinAutomationService(
      JSON.parse(this.browserConfig.cookies ?? "{}"),
      this.browserConfig.userAgent as string,
      this.proxyHeader,
    );
  }

  async sendInvite(
    profileUrl: string,
    message: string,
    options?: LinkedinServiceMethodOptions,
  ) {
    try {
      options?.onStart?.();
      await this.linkedinAutomationService.sendConenctionInvite(
        profileUrl,
        message,
        options,
      );
      options?.onSuccess?.();
    } catch (err) {
      LinkedinService.handleError(err, (error) => {
        options?.onError?.(error);
      });
    }
  }

  async scrapeConnections(options?: LinkedinServiceMethodOptions) {
    try {
      options?.onStart?.();
      const result = await this.linkedinAutomationService.getConnections();
      options?.onSuccess?.();
      return result;
    } catch (err) {
      LinkedinService.handleError(err, (error) => {
        options?.onError?.(error);
      });
    }
  }

  async sendMessage(
    profileUrl: string,
    message: string,
    options?: LinkedinServiceMethodOptions,
  ) {
    try {
      options?.onStart?.();
      await this.linkedinAutomationService.sendMessageToConnection(
        profileUrl,
        message,
      );
      options?.onSuccess?.();
    } catch (err) {
      LinkedinService.handleError(err, (error) => {
        options?.onError?.(error);
      });
    }
  }

  private static handleError(arr: unknown, cb?: (err: StandardError) => void) {
    const error = ErrorParser.parse(arr);
    logger.error("Error in LinkedinService", {
      error: error.message,
      details: error.details,
    });
    cb?.(error);
  }
}
