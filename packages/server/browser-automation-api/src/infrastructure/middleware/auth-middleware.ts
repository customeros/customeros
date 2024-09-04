import { Request, Response, NextFunction } from "express";
import { header, validationResult } from "express-validator";

import { ErrorParser } from "@/util/error";
import { logger } from "@/infrastructure/logger";
import { type UserRepository } from "@/infrastructure/persistance/neo4j/repositories";
import {
  type BrowserConfigsRepository,
  type TenantWebhookApiKeysRepository,
} from "@/infrastructure/persistance/postgresql/repositories";

export class AuthMiddleware {
  constructor(
    private userRepository: UserRepository,
    private browserConfigsRepository: BrowserConfigsRepository,
    private tenantWebhookApiKeysRepository: TenantWebhookApiKeysRepository,
  ) {
    this.checkApiKey = this.checkApiKey.bind(this);
    this.getValidators = this.getValidators.bind(this);
  }

  private async checkApiKey(req: Request, res: Response, next: NextFunction) {
    const headersValidationErrors = validationResult(req);
    if (!headersValidationErrors.isEmpty()) {
      return res.status(400).json({
        success: false,
        message: "Authorization headers required.",
        errors: headersValidationErrors.array(),
      });
    }

    try {
      const username = req.header("X-OPENLINE-USERNAME") as string;
      const apiKey = req.header("X-OPENLINE-API-KEY") as string;

      const tenantApiKey =
        await this.tenantWebhookApiKeysRepository.selectByApiKey(apiKey);
      const tenantName = tenantApiKey?.tenantName ?? "";

      if (!tenantApiKey || tenantApiKey.key !== apiKey) {
        return res.status(401).json({
          success: false,
          message: "Unauthorized",
        });
      }

      if (!tenantName) {
        return res.status(401).json({
          success: false,
          message: "Unauthorized. Tenant not found.",
        });
      }

      const user = await this.userRepository.getUserByEmail(
        tenantName,
        username,
      );

      if (!user) {
        return res.status(401).json({
          success: false,
          message: "Unauthorized. User not found.",
        });
      }

      const browserConfig = await this.browserConfigsRepository.selectByUserId(
        user.id,
      );

      res.locals.tenantName = tenantName;
      res.locals.user = user;
      res.locals.browserConfig = browserConfig;
      next();
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in AuthMiddleware", {
        error: error.message,
        details: error.details,
      });
      return res.status(500).json({
        success: false,
        message: "Internal server error",
      });
    }
  }

  getValidators() {
    return [
      header("X-OPENLINE-API-KEY").exists().withMessage("API key is required"),
      header("X-OPENLINE-USERNAME")
        .exists()
        .withMessage("Username is required"),
      this.checkApiKey,
    ];
  }
}
