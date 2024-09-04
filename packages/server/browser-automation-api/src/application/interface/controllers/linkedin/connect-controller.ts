import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { ScheduleService } from "@/application/services/schedule-service";

export class ConnectController {
  private scheduleService = ScheduleService.getInstance();

  constructor() {
    this.sendConnectionInvite = this.sendConnectionInvite.bind(this);
  }

  async sendConnectionInvite(req: Request, res: Response) {
    const { profileUrl, message, dryRun } = req.body;
    const validationErrors = validationResult(req);

    if (!validationErrors.isEmpty()) {
      return res.status(400).send({
        success: false,
        message: "Validation failed",
        errors: validationErrors.array(),
      });
    }

    try {
      const automationRun = await this.scheduleService.createAutomationRun(
        res.locals.browserConfig,
        "SEND_CONNECTION_REQUEST",
        { profileUrl, message, dryRun },
      );

      res.send({
        success: true,
        message: "Connection request sent successfully",
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
        message: "Failed to send connection request",
        error: error.message,
      });
    }
  }
}
