import { ErrorParser, StandardError } from "@/util/error";
import { BrowserAutomationRun } from "@/domain/models/browser-automation-run";
import type {
  BrowserConfigsRepository,
  BrowserAutomationRunsRepository,
  BrowserAutomationRunErrorsRepository,
  BrowserAutomationRunResultsRepository,
} from "@/infrastructure/persistance/postgresql/repositories";

import type { LinkedinServiceFactory } from "./linkedin/linkedin-service-factory";
import { logger } from "@/infrastructure";

export class AutomationRunnerService {
  constructor(
    private linkedinServiceFactory: LinkedinServiceFactory,
    private runsRepository: BrowserAutomationRunsRepository,
    private resultsRepository: BrowserAutomationRunResultsRepository,
    private errorsRepository: BrowserAutomationRunErrorsRepository,
    private configRepository: BrowserConfigsRepository,
  ) {}

  async runAutomation(browserAutomationRun: BrowserAutomationRun) {
    try {
      const linkedinService =
        await this.linkedinServiceFactory.createForRun(browserAutomationRun);

      browserAutomationRun.start();
      await this.runsRepository.updateById(browserAutomationRun.toDTO());

      let result;
      let errorValue;

      const payload = BrowserAutomationRun.parsePayload(
        browserAutomationRun.payload,
      );

      switch (browserAutomationRun.type) {
        case "FIND_CONNECTIONS":
          const [res, err] = await linkedinService.scrapeConnections();
          result = res;
          errorValue = err;

          break;
        case "DOWNLOAD_CONNECTIONS":
          result = await linkedinService.downloadConnections();
          break;
        case "FIND_COMPANY_PEOPLE":
          result = await linkedinService.scrapeCompanyPeople(payload);
          break;
        case "SEND_CONNECTION_REQUEST":
          result = await linkedinService.sendInvite(payload);
          break;
        case "SEND_MESSAGE":
          result = await linkedinService.sendMessage(payload);
          break;
        case "GET_MESSAGES":
          result = await linkedinService.retrieveMessages(payload);
          break;
        case "CHECK_CONNECTION_STATUS":
          result = await linkedinService.checkConnectionStatus(payload);
          break;
        default:
          throw new StandardError({
            code: "APPLICATION_ERROR",
            message: `Unknown automation run type: ${browserAutomationRun.type}`,
            severity: "critical",
          });
      }

      if (errorValue) {
        browserAutomationRun.retry();

        await this.errorsRepository.insert({
          runId: browserAutomationRun.id,
          errorMessage: errorValue.message,
          errorDetails: errorValue.details,
          errorCode: errorValue.reference,
          errorType: errorValue.code,
        });
        logger.error("Automation run failed", {
          error: errorValue.message,
          details: errorValue.reference ?? errorValue.details,
          source: "AutomationRunnerService",
        });
        if (errorValue.reference === "S001") {
          await this.configRepository.updateByUserId({
            userId: browserAutomationRun.userId,
            tenant: browserAutomationRun.tenant,
            sessionStatus: "INVALID",
          });
        }
      } else {
        browserAutomationRun.complete();
      }

      await this.runsRepository.updateById(browserAutomationRun.toDTO());
      await this.resultsRepository.insert({
        runId: browserAutomationRun.id,
        type: browserAutomationRun.type,
        resultData: JSON.stringify(result),
      });
    } catch (err) {
      const error = ErrorParser.parse(err);

      browserAutomationRun.fail();
      await this.runsRepository.updateById(browserAutomationRun.toDTO());

      await this.errorsRepository.insert({
        runId: browserAutomationRun.id,
        errorMessage: error.message,
        errorDetails: error.details,
        errorCode: error.reference,
        errorType: error.code,
      });
      if (error.reference === "S001") {
        await this.configRepository.updateByUserId({
          userId: browserAutomationRun.userId,
          tenant: browserAutomationRun.tenant,
          sessionStatus: "INVALID",
        });
      }

      logger.error("Automation run failed", {
        error: error.message,
        details: error.reference ?? error.details,
        source: "AutomationRunnerService",
      });
    } finally {
      return;
    }
  }
}
