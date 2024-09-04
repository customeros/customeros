import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import { LinkedinService } from "@/application/services/linkedin/linkedin-service";

export class ConnectionsController {
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

    const linkedinService = new LinkedinService(res.locals.browserConfig);

    try {
      const connectionUrls = await linkedinService.scrapeConnections();
      res.send({ success: true, connectionUrls });
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
