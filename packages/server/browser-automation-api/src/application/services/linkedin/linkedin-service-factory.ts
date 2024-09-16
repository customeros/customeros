import { Proxy } from "@/domain/models/proxy";
import { logger } from "@/infrastructure/logger";
import { ErrorParser, StandardError } from "@/util/error";
import { BrowserConfig } from "@/domain/models/browser-config";
import type { BrowserAutomationRun } from "@/domain/models/browser-automation-run";
import type {
  ProxyPoolRepository,
  BrowserConfigsRepository,
  AssignedProxiesRepository,
} from "@/infrastructure/persistance/postgresql/repositories";

import { LinkedinService } from "./linkedin-service";

export class LinkedinServiceFactory {
  constructor(
    private configsRepository: BrowserConfigsRepository,
    private proxyPoolRepository: ProxyPoolRepository,
    private assignedProxiesRepository: AssignedProxiesRepository,
  ) {}

  async createForRun(browserAutomationRun: BrowserAutomationRun) {
    try {
      const browserConfig = await this.configsRepository.selectById(
        browserAutomationRun.browserConfigId,
      );
      if (!browserConfig) {
        throw new StandardError({
          code: "APPLICATION_ERROR",
          message: `Failed to find browser config with id: ${browserAutomationRun.browserConfigId}.`,
          severity: "critical",
        });
      }

      const assignedProxy = await this.assignedProxiesRepository.selectByUserId(
        browserConfig.userId,
      );
      if (!assignedProxy) {
        throw new StandardError({
          code: "APPLICATION_ERROR",
          message: `Assigned proxy for user with id ${browserConfig.userId} not found`,
          severity: "critical",
        });
      }

      const proxy = await this.proxyPoolRepository.selectById(
        assignedProxy.proxyPoolId,
      );
      if (!proxy) {
        throw new StandardError({
          code: "APPLICATION_ERROR",
          message: `Proxy with id ${assignedProxy.proxyPoolId} not found`,
          severity: "critical",
        });
      }

      const proxyHeader = Proxy.toBrowserHeader(proxy);

      return new LinkedinService(new BrowserConfig(browserConfig), proxyHeader);
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error(error.message, {
        source: "LinkedinServiceFactory",
      });
      throw error;
    }
  }
}
