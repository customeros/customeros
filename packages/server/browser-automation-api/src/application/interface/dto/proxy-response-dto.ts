import { type Proxy } from "@/domain/models/proxy";

export class ProxyResponseDTO {
  id: number;
  url: string;
  createdAt: string | null;
  updatedAt: string | null;

  constructor(values: Proxy) {
    this.id = values.id;
    this.url = values.url;
    this.createdAt = values.createdAt;
    this.updatedAt = values.updatedAt;
  }
}
