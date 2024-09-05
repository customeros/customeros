import sample from "lodash/sample";

import { logger } from "@/infrastructure/logger";
import { ErrorParser, StandardError } from "@/util/error";
import { AssignedProxy } from "@/domain/models/assigned-proxy";
import { Proxy, ProxyPoolPayload } from "@/domain/models/proxy";
import { ProxyPoolRepository } from "@/infrastructure/persistance/postgresql/repositories/proxy-pool-repository";
import { AssignedProxiesRepository } from "@/infrastructure/persistance/postgresql/repositories/assigned-proxies-repository";

export class ProxyService {
  private proxyPoolRepository: ProxyPoolRepository;
  private assignedProxiesRepository: AssignedProxiesRepository;

  constructor() {
    this.proxyPoolRepository = new ProxyPoolRepository();
    this.assignedProxiesRepository = new AssignedProxiesRepository();
  }

  async getProxyPool(): Promise<Proxy[] | undefined> {
    try {
      const values = await this.proxyPoolRepository.selectAll();
      return (values ?? []).map((value) => new Proxy(value));
    } catch (err) {
      ProxyService.handleError(err);
    }
  }

  async getProxy(id: number) {
    try {
      const values = await this.proxyPoolRepository.selectById(id);
      if (!values) {
        return null;
      }
      return new Proxy(values);
    } catch (err) {
      ProxyService.handleError(err);
    }
  }

  async createProxy(payload: ProxyPoolPayload) {
    try {
      return await Proxy.create(payload, this.proxyPoolRepository);
    } catch (err) {
      ProxyService.handleError(err);
    }
  }

  async getAssignedProxy(userId: string) {
    try {
      const values =
        await this.assignedProxiesRepository.selectByUserId(userId);
      if (!values) {
        return null;
      }
      return new AssignedProxy(values);
    } catch (err) {
      ProxyService.handleError(err);
    }
  }

  async assignProxy(userId: string, tenant: string) {
    try {
      const proxyPool = await this.proxyPoolRepository.selectAll();
      if (!proxyPool) {
        throw new StandardError({
          code: "INVARIANT_ERROR",
          message: "Invariant violation: No proxy pool available",
          severity: "critical",
          details: "Proxy pool is empty",
        });
      }

      const proxyPoolId = sample(proxyPool)?.id as number;

      return await AssignedProxy.create(
        {
          userId,
          tenant,
          proxyPoolId,
        },
        this.assignedProxiesRepository,
      );
    } catch (err) {
      ProxyService.handleError(err);
    }
  }

  async rotateProxy(userId: string) {
    try {
      const assignedProxy = await this.getAssignedProxy(userId);
      if (!assignedProxy) {
        throw new StandardError({
          code: "APPLICATION_ERROR",
          message: "Assigned proxy not found",
          severity: "low",
          details: "No assigned proxy found for the user",
        });
      }

      const pool = await this.proxyPoolRepository.selectAll();
      if (!pool) {
        throw new StandardError({
          code: "INVARIANT_ERROR",
          message: "Invariant violation: No proxy pool available",
          severity: "critical",
          details: "Proxy pool is empty",
        });
      }

      const proxyPoolId =
        sample(pool.filter((p) => p.id !== assignedProxy?.proxyPoolId))?.id ??
        assignedProxy?.proxyPoolId;

      return await AssignedProxy.update(
        {
          ...assignedProxy,
          proxyPoolId,
        },
        this.assignedProxiesRepository,
      );
    } catch (err) {
      ProxyService.handleError(err);
    }
  }

  static handleError(err: unknown) {
    const error = ErrorParser.parse(err);
    logger.error("Error in AssignedProxyService", {
      message: error.message,
      details: error.details,
    });
    throw error;
  }
}
