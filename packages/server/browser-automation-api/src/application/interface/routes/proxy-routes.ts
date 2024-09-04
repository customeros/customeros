import { Router } from "@/infrastructure";

import { ProxyController } from "../controllers/proxy-controller";
import { getProxyValidator, addProxyValidator } from "../validators";

export class ProxyRouter {
  public router = new Router().instance;
  private proxyController = new ProxyController();

  constructor() {
    this.router.get("/pool", this.proxyController.getProxyPool);
    this.router.get("/assigned", this.proxyController.getAssignedProxy);
    this.router.get("/rotate", this.proxyController.rotateProxy);
    this.router.post("/add", addProxyValidator, this.proxyController.addProxy);
    this.router.get("/:id", getProxyValidator, this.proxyController.getProxy);
  }
}
