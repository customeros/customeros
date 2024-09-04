import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import {
  AssignedProxyInsert,
  AssignedProxiesTable,
  type AssignedProxiesRepository,
} from "@/infrastructure/persistance/postgresql/repositories/assigned-proxies-repository";

export type AssignedProxyPayload = Pick<
  AssignedProxyInsert,
  "proxyPoolId" | "userId" | "tenant"
>;

export class AssignedProxy {
  id: number;
  proxyPoolId: number;
  userId: string;
  tenant: string;
  createdAt: string | null = null;
  updatedAt: string | null = null;

  constructor(values: AssignedProxiesTable) {
    this.id = values.id;
    this.proxyPoolId = values.proxyPoolId;
    this.userId = values.userId;
    this.tenant = values.tenant;
    this.createdAt = values.createdAt;
    this.updatedAt = values.updatedAt;
  }

  static async create(
    values: AssignedProxyPayload,
    assignedProxiesRepository: AssignedProxiesRepository,
  ) {
    try {
      return await assignedProxiesRepository.insert(values);
    } catch (err) {
      AssignedProxy.handleError(err);
    }
  }

  static async update(
    values: AssignedProxyPayload,
    assignedProxiesRepository: AssignedProxiesRepository,
  ) {
    try {
      return await assignedProxiesRepository.updateByUserId(values);
    } catch (err) {
      AssignedProxy.handleError(err);
    }
  }

  static handleError(err: any) {
    const error = ErrorParser.parse(err);
    logger.error("Error in AssignedProxy", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
