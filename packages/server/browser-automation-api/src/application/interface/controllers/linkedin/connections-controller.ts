import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { ScheduleService } from "@/application/services/schedule-service";

export class ConnectionsController {
  private scheduleService = ScheduleService.getInstance();

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
      const automationRun = await this.scheduleService.createAutomationRun(
        res.locals.browserConfig,
        "FIND_CONNECTIONS",
      );
      res.send({
        success: true,
        message: "Browser automation scheduled successfully",
        data: automationRun,
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in ConnectController", {
        error: error.message,
        details: error.details,
      });
      res.status(500).send({
        success: false,
        message: "Failed to get connections",
        error: error.message,
      });
    }
  }
}
