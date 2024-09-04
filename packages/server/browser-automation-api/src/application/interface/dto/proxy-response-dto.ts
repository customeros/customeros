import { type ProxyPool } from "@/domain/models/proxy-pool";

export class ProxyResponseDTO {
  id: number;
  url: string;
  createdAt: string | null;
  updatedAt: string | null;

  constructor(values: ProxyPool) {
    this.id = values.id;
    this.url = values.url;
    this.createdAt = values.createdAt;
    this.updatedAt = values.updatedAt;
  }
}
