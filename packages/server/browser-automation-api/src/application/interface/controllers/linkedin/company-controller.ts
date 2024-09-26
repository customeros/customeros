import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { BrowserAutomationRunService } from "@/application/services/browser-automation-run-service";
import { BrowserAutomationRunsRepository } from "@/infrastructure/persistance/postgresql/repositories";

export class CompanyController {
  private browserAutomationRunService = new BrowserAutomationRunService(
    new BrowserAutomationRunsRepository(),
  );

  constructor() {
    this.scrapeCompanyPeople = this.scrapeCompanyPeople.bind(this);
  }

  async scrapeCompanyPeople(req: Request, res: Response) {
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
          type: "FIND_COMPANY_PEOPLE",
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
      logger.error("Error in CompanyController", {
        error: error.message,
        details: error.details,
        source: "CompanyController",
      });
      res.status(500).send({
        success: false,
        message: "Failed to get company people",
        error: error.message,
      });
    }
  }
}
