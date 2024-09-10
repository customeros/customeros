import { logger, Server } from "@/infrastructure";
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
import { setTimeout } from "timers/promises";

const PORT = 3000;

export class App {
  private server = new Server(PORT);
  private authMiddleware = new AuthMiddleware(
    new UserRepository(),
    new BrowserConfigsRepository(),
    new TenantWebhookApiKeysRepository(),
  );
  private scheduleService = ScheduleService.getInstance();

  constructor() {}

  public init() {
    this.server.instance.use(this.authMiddleware.getValidators());
    this.server.instance.use("/browser", new BrowserRouter().router);
    this.server.instance.use("/linkedin", new LinkedinRouter().router);
    this.server.instance.use("/proxy", new ProxyRouter().router);

    this.scheduleService.pollBrowserAutomationRuns();
    this.registerShutdownHooks();
  }

  private blockIO() {
    if (process.stdin.setRawMode) {
      process.stdin.setRawMode(true);
      process.stdin.pause();
    }
    process.stdout.write("\x1B[?25l");
  }

  private unblockIO() {
    if (process.stdin.setRawMode) {
      process.stdin.setRawMode(false);
      process.stdin.resume();
    }
    process.stdout.write("\x1B[?25h");
  }

  private registerShutdownHooks() {
    process.on("SIGTERM", () => {
      this.blockIO();
      this.server.stop(async () => {
        await this.scheduleService.shutdown();
        this.unblockIO();
        process.exit(0);
      });
    });
    process.on("SIGINT", () => {
      this.blockIO();
      this.server.stop(async () => {
        await this.scheduleService.shutdown();
        this.unblockIO();
        process.exit(0);
      });
    });
  }
}
