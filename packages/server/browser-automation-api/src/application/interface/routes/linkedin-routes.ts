import { Router } from "@/infrastructure";

import { CompanyController } from "../controllers/linkedin/company-controller";
import { ConnectController } from "../controllers/linkedin/connect-controller";
import { MessagesController } from "../controllers/linkedin/messages-controller";
import { ConnectionsController } from "../controllers/linkedin/connections-controller";
import { connectValidators, sendMessageValidators } from "../validators";

export class LinkedinRouter {
  public router = new Router().instance;
  private companyController = new CompanyController();
  private connectController = new ConnectController();
  private messagesController = new MessagesController();
  private connectionsController = new ConnectionsController();

  constructor() {
    this.router.post(
      "/company/people",
      this.companyController.scrapeCompanyPeople,
    );
    this.router.get(
      "/connections",
      this.connectionsController.scrapeConnections,
    );
    this.router.post(
      "/message",
      ...sendMessageValidators,
      this.messagesController.sendMessage,
    );
    this.router.post(
      "/connect",
      ...connectValidators,
      this.connectController.sendConnectionInvite,
    );
  }
}
