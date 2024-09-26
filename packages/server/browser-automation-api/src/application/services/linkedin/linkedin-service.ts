import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { BrowserConfig } from "@/domain/models/browser-config";
import { LinkedinAutomationService } from "@/infrastructure/scraper/services/linkedin-automation-service";

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

  async sendInvite(payload: unknown) {
    const { profileUrl, message, dryRun } = payload as {
      profileUrl: string;
      message: string;
      dryRun?: boolean;
    };

    try {
      logger.info("Sending connection invite...", {
        source: "LinkedinService",
      });

      await this.linkedinAutomationService.sendConenctionInvite(
        profileUrl,
        message,
        { dryRun },
      );

      logger.info("Connection invite sent.", {
        source: "LinkedinService",
      });
      return { profileUrl, message: "Connection invite sent successfully" };
    } catch (err) {
      logger.info("Failed to send connection invite.", {
        source: "LinkedinService",
      });
      throw LinkedinService.handleError(err);
    }
  }

  async scrapeConnections(payload?: { lastPageVisited?: number }) {
    try {
      logger.info("Scraping connections...", {
        source: "LinkedinService",
      });

      const result = await this.linkedinAutomationService.getConnections(
        payload?.lastPageVisited,
      );

      logger.info("Connections scraped.", {
        source: "LinkedinService",
      });
      return result;
    } catch (err) {
      logger.info("Failed to scrape connections", {
        source: "LinkedinService",
      });
      throw LinkedinService.handleError(err);
    }
  }

  async scrapeCompanyPeople(payload: unknown) {
    const { companyName, dryRun } = payload as {
      companyName: string;
      dryRun?: boolean;
    };

    try {
      logger.info("Scraping company people...", {
        source: "LinkedinService",
      });

      const result =
        await this.linkedinAutomationService.getCompanyPeople(companyName);

      logger.info("Company people scraped.", {
        source: "LinkedinService",
      });
      return result;
    } catch (err) {
      logger.info("Failed to scrape company people", {
        source: "LinkedinService",
      });
      throw LinkedinService.handleError(err);
    }
  }

  async sendMessage(payload: unknown) {
    const { profileUrl, message, dryRun } = payload as {
      profileUrl: string;
      message: string;
      dryRun?: boolean;
    };

    try {
      logger.info("Sending message", {
        source: "LinkedinService",
      });

      await this.linkedinAutomationService.sendMessageToConnection(
        profileUrl,
        message,
        { dryRun },
      );

      logger.info("Message sent", {
        source: "LinkedinService",
      });
      return { profileUrl, message: "Message sent successfully" };
    } catch (err) {
      logger.info("Failed to send message", {
        source: "LinkedinService",
      });
      throw LinkedinService.handleError(err);
    }
  }

  private static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in LinkedinService", {
      error: error.message,
      details: error.reference ?? error.details,
      source: "LinkedinService",
    });
    return error;
  }
}
