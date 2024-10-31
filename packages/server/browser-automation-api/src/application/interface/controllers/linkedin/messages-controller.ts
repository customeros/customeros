import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { BrowserAutomationRunService } from "@/application/services/browser-automation-run-service";
import { BrowserAutomationRunsRepository } from "@/infrastructure/persistance/postgresql/repositories";

export class MessagesController {
  private browserAutomationRunService = new BrowserAutomationRunService(
    new BrowserAutomationRunsRepository(),
  );

  constructor() {
    this.sendMessage = this.sendMessage.bind(this);
    this.retrieveMessages = this.retrieveMessages.bind(this);
  }

  async sendMessage(req: Request, res: Response) {
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
          type: "SEND_MESSAGE",
          userId: res.locals.user.id,
          payload: JSON.stringify(req.body),
        },
      );

      res.send({
        success: true,
        message: "Browser automation scheduled successfully",
        data: newAutomationRun?.toDTO(),
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in MessagesController", {
        error: error.message,
        details: error.details,
        source: "MessagesController",
      });
      res.status(500).send({
        success: false,
        message: "Failed to send message",
        error: error.message,
      });
    }
  }

  async retrieveMessages(req: Request, res: Response) {
    try {
      const newAutomationRun = await this.browserAutomationRunService.createRun(
        {
          browserConfigId: res.locals.browserConfig.id,
          tenant: res.locals.tenantName,
          type: "GET_MESSAGES",
          userId: res.locals.user.id,
          payload: JSON.stringify(req.body),
        },
      );

      res.send({
        success: true,
        message: "Browser automation scheduled successfully",
        data: newAutomationRun?.toDTO(),
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in MessagesController", {
        error: error.message,
        details: error.details,
        source: "MessagesController",
      });
      res.status(500).send({
        success: false,
        message: "Failed to get messages",
        error: error.message,
      });
    }
  }
}
