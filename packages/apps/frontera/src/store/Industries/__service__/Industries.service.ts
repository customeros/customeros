import { Transport } from '@infra/transport.ts';

import GetIndustriesDocument from './getTenantIndustries.graphql';
import {
  GetTenantIndustriesQuery,
  GetTenantIndustriesQueryVariables,
} from './getTenantIndustries.generated.ts';

class IndustriesService {
  private static instance: IndustriesService | null = null;
  private transport = Transport.getInstance();

  constructor() {}

  static getInstance(): IndustriesService {
    if (!IndustriesService.instance) {
      IndustriesService.instance = new IndustriesService();
    }

    return IndustriesService.instance;
  }

  async getIndustries() {
    return this.transport.graphql.request<
      GetTenantIndustriesQuery,
      GetTenantIndustriesQueryVariables
    >(GetIndustriesDocument, {});
  }
}

export { IndustriesService };
