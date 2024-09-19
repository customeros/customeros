import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { BrowserAutomationRunService } from "@/application/services/browser-automation-run-service";
import { BrowserAutomationRunsRepository } from "@/infrastructure/persistance/postgresql/repositories";

export class ConnectionsController {
  private browserAutomationRunService = new BrowserAutomationRunService(
    new BrowserAutomationRunsRepository(),
  );

  constructor() {
    this.scrapeConnections = this.scrapeConnections.bind(this);
  }

  async scrapeConnections(req: Request, res: Response) {
    const validationErrors = validationResult(req);

    if (!validationErrors.isEmpty()) {
      return res.status(400).send({
        success: false,
        message: "Validation failed",
        errors: validationErrors.array(),
      });
    }

    if (!res.locals.browserConfig) {
      return res.status(400).send({
        success: false,
        message: "Browser config not found",
      });
    }

    try {
      const newAutomationRun = await this.browserAutomationRunService.createRun(
        {
          browserConfigId: res.locals.browserConfig.id,
          tenant: res.locals.tenantName,
          type: "FIND_CONNECTIONS",
          userId: res.locals.user.id,
          payload: "",
        },
      );

      res.send({
        success: true,
        message: "Browser automation scheduled successfully",
        data: newAutomationRun?.toDTO(),
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in ConnectController", {
        error: error.message,
        details: error.details,
        source: "ConnectionsController",
      });
      res.status(500).send({
        success: false,
        message: "Failed to get connections",
        error: error.message,
      });
    }
  }
}
