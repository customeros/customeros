import { Request, Response } from "express";
import { validationResult } from "express-validator";

import { logger } from "@/infrastructure";
import { ErrorParser, StandardError } from "@/util/error";
import { ProxyService } from "@/application/services/proxy-service";

import { ProxyResponseDTO } from "../dto/proxy-response-dto";

export class ProxyController {
  private proxyService = new ProxyService();

  constructor() {
    this.getProxy = this.getProxy.bind(this);
    this.addProxy = this.addProxy.bind(this);
    this.rotateProxy = this.rotateProxy.bind(this);
    this.getProxyPool = this.getProxyPool.bind(this);
    this.getAssignedProxy = this.getAssignedProxy.bind(this);
  }

  async getProxyPool(_req: Request, res: Response) {
    try {
      const proxyPool = await this.proxyService.getProxyPool();

      res.send({
        success: true,
        data: proxyPool?.map((values) => new ProxyResponseDTO(values)),
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in ProxyPoolController", {
        error: error.message,
        details: error.details,
      });
      res.status(500).send({
        success: false,
        message: "Failed to get proxy pool",
        error: error.message,
      });
    }
  }

  async getProxy(req: Request, res: Response) {
    const validationErrors = validationResult(req);

    if (!validationErrors.isEmpty()) {
      return res.status(400).send({
        success: false,
        message: "Validation failed",
        errors: validationErrors.array(),
      });
    }

    const proxyId = (() => {
      const { id } = req.params;
      return parseInt(id);
    })();

    try {
      const result = await this.proxyService.getProxy(proxyId);
      if (!result) {
        return res.status(404).send({
          success: false,
          message: "Proxy not found",
        });
      }

      res.send({
        success: true,
        data: new ProxyResponseDTO(result),
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in ProxyPoolController", {
        error: error.message,
        details: error.details,
      });
      res.status(500).send({
        success: false,
        message: "Failed to get proxy",
        error: error.message,
      });
    }
  }

  async addProxy(req: Request, res: Response) {
    const userId = res.locals.user.id;
    const tenant = res.locals.tenantName;

    const validationErrors = validationResult(req);

    if (!validationErrors.isEmpty()) {
      return res.status(400).send({
        success: false,
        message: "Validation failed",
        errors: validationErrors.array(),
      });
    }

    try {
      const proxy = await this.proxyService.createProxy(req.body);
      if (!proxy) {
        throw new StandardError({
          code: "INTERNAL_ERROR",
          message: "Failed to add new proxy to the pool",
          severity: "critical",
        });
      }

      res.send({
        success: true,
        message: "Proxy added to pool successfully",
        data: new ProxyResponseDTO(proxy),
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in ProxyPoolController", {
        error: error.message,
        details: error.details,
      });
      res.status(500).send({
        success: false,
        message: "Failed to add new proxy to the pool",
        error: error.message,
      });
    }
  }

  async getAssignedProxy(req: Request, res: Response) {
    const userId = res.locals.user.id;
    const tenant = res.locals.tenantName;
    const validationErrors = validationResult(req);

    if (!validationErrors.isEmpty()) {
      return res.status(400).send({
        success: false,
        message: "Validation failed",
        errors: validationErrors.array(),
      });
    }

    try {
      let assignedProxy = await this.proxyService.getAssignedProxy(userId);

      if (!assignedProxy) {
        assignedProxy = await this.proxyService.assignProxy(userId, tenant);
      }

      res.send({
        success: true,
        data: assignedProxy,
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in AssignedProxyController", {
        error: error.message,
        details: error.details,
      });
      res.status(500).send({
        success: false,
        message: "Failed to get assigned proxy",
        error: error.message,
      });
    }
  }

  async rotateProxy(req: Request, res: Response) {
    const userId = res.locals.user.id;
    const validationErrors = validationResult(req);

    if (!validationErrors.isEmpty()) {
      return res.status(400).send({
        success: false,
        message: "Validation failed",
        errors: validationErrors.array(),
      });
    }

    try {
      const proxyPool = await this.proxyService.rotateProxy(userId);

      res.send({
        success: true,
        message: "Assigned proxy rotated successfully",
        data: proxyPool,
      });
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in AssignedProxyController", {
        error: error.message,
        details: error.details,
      });
      res.status(500).send({
        success: false,
        message: "Failed to rotate assigned proxy",
        error: error.message,
      });
    }
  }
}
