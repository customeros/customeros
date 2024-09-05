import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";
import {
  ProxyPoolTable,
  ProxyPoolInsert,
  type ProxyPoolRepository,
} from "@/infrastructure/persistance/postgresql/repositories/proxy-pool-repository";

export type ProxyPoolPayload = Pick<
  ProxyPoolInsert,
  "url" | "username" | "password" | "enabled"
>;

export class Proxy {
  id: number;
  url: string;
  username: string;
  password: string;
  createdAt: string | null = null;
  updatedAt: string | null = null;

  constructor(values: ProxyPoolTable) {
    this.id = values.id;
    this.url = values.url;
    this.username = values.username;
    this.password = values.password;
    this.createdAt = values.createdAt;
    this.updatedAt = values.updatedAt;
  }

  static async create(
    values: ProxyPoolPayload,
    proxyPoolRepository: ProxyPoolRepository,
  ) {
    try {
      return await proxyPoolRepository.insert(values);
    } catch (err) {
      Proxy.handleError(err);
    }
  }

  static async update(
    values: ProxyPoolPayload,
    proxyPoolRepository: ProxyPoolRepository,
  ) {
    try {
      return await proxyPoolRepository.updateById(values);
    } catch (err) {
      Proxy.handleError(err);
    }
  }

  static toBrowserHeader(values: ProxyPoolTable) {
    return JSON.stringify({
      proxy: {
        server: values.url,
        username: values.username,
        password: values.password,
      },
    });
  }

  static handleError(err: any) {
    const error = ErrorParser.parse(err);
    logger.error("Error in ProxyPool", {
      error: error.message,
      details: error.details,
    });
    throw error;
  }
}
