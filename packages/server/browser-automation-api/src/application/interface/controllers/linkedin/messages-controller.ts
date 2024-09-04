import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { ScheduleService } from "@/application/services/schedule-service";

export class MessagesController {
  private scheduleService = ScheduleService.getInstance();

  constructor() {
    this.sendMessage = this.sendMessage.bind(this);
  }

  async sendMessage(req: Request, res: Response) {
    const { profileUrl, message, dryRun } = req?.body;
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
        "SEND_MESSAGE",
        { profileUrl, message, dryRun },
      );

      res.send({
        success: true,
        message: "Browser automation scheduled successfully",
        data: automationRun,
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in MessagesController", {
        error: error.message,
        details: error.details,
      });
      res.status(500).send({
        success: false,
        message: "Failed to send message",
        error: error.message,
      });
    }
  }
}
