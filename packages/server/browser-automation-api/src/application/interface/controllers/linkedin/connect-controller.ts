import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { BrowserAutomationRunService } from "@/application/services/browser-automation-run-service";
import { BrowserAutomationRunsRepository } from "@/infrastructure/persistance/postgresql/repositories";

export class ConnectController {
  private browserAutomationRunService = new BrowserAutomationRunService(
    new BrowserAutomationRunsRepository(),
  );

  constructor() {
    this.sendConnectionInvite = this.sendConnectionInvite.bind(this);
  }

  async sendConnectionInvite(req: Request, res: Response) {
    const validationErrors = validationResult(req);

    if (!validationErrors.isEmpty()) {
      return res.status(400).send({
        success: false,
        message: "Validation failed",
        errors: validationErrors.array(),
      });
    }

    try {
      const newAutomationRun = await this.browserAutomationRunService.createRun(
        {
          browserConfigId: res.locals.browserConfig.id,
          tenant: res.locals.tenantName,
          type: "SEND_CONNECTION_REQUEST",
          userId: res.locals.user.id,
          payload: req.body,
        },
      );

      res.send({
        success: true,
        message: "Connection request sent successfully",
        data: newAutomationRun?.toDTO(),
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in ConnectController", {
        error: error.message,
        details: error.details,
        source: "ConnectController",
      });
      res.status(500).send({
        success: false,
        message: "Failed to send connection request",
        error: error.message,
      });
    }
  }
}
