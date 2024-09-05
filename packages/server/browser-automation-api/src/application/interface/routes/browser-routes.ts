import { Router } from "@/infrastructure";

import { BrowserController } from "../controllers/browser-controller";
import { getBrowserRunByIdValidators } from "../validators/browser-validators";

export class BrowserRouter {
  public router = new Router().instance;
  private browserController = new BrowserController();

  constructor() {
    this.router.get("/config", this.browserController.getBrowserConfig);
    this.router.post("/config", this.browserController.createBrowserConfig);
    this.router.patch("/config", this.browserController.updateBrowserConfig);
    this.router.get("/runs", this.browserController.getBrowserAutomationRuns);
    this.router.get(
      "/run/:id",
      getBrowserRunByIdValidators,
      this.browserController.getBrowserAutomationRun,
    );
  }
}
