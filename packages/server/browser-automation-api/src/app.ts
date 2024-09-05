import { Server } from "@/infrastructure";
import { AuthMiddleware } from "@/infrastructure/middleware";
import { UserRepository } from "@/infrastructure/persistance/neo4j/repositories";
import {
  BrowserConfigsRepository,
  TenantWebhookApiKeysRepository,
} from "@/infrastructure/persistance/postgresql/repositories";
import {
  ProxyRouter,
  BrowserRouter,
  LinkedinRouter,
} from "@/application/interface/routes";
import { ScheduleService } from "@/application/services/schedule-service";

const PORT = 3000;

export class App {
  constructor() {
    const server = new Server(PORT);
    const authMiddleware = new AuthMiddleware(
      new UserRepository(),
      new BrowserConfigsRepository(),
      new TenantWebhookApiKeysRepository(),
    );

    server.instance.use(authMiddleware.getValidators());
    server.instance.use("/browser", new BrowserRouter().router);
    server.instance.use("/linkedin", new LinkedinRouter().router);
    server.instance.use("/proxy", new ProxyRouter().router);

    ScheduleService.getInstance().pollBrowserAutomationRuns();
  }
}
