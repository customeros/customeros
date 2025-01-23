import type { Transport } from '@infra/transport.ts';

import GlobalCacheQueryDocument from './getGlobalCache.graphql';
import { GlobalCacheQuery } from './getGlobalCache.generated.ts';
import UpdateOnboardingDetailsMutationDocument from './updateUserOnboardingDetails.graphql';
import {
  UserUpdateOnboardingDetailsMutation,
  UserUpdateOnboardingDetailsMutationVariables,
} from './updateUserOnboardingDetails.generated.ts';

class GlobalCacheService {
  private static instance: GlobalCacheService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): GlobalCacheService {
    if (!GlobalCacheService.instance) {
      GlobalCacheService.instance = new GlobalCacheService(transport);
    }

    return GlobalCacheService.instance;
  }

  async getGlobalCache() {
    return this.transport.graphql.request<GlobalCacheQuery>(
      GlobalCacheQueryDocument,
    );
  }

  async updateOnboardingDetails(
    payload: UserUpdateOnboardingDetailsMutationVariables,
  ) {
    return this.transport.graphql.request<
      UserUpdateOnboardingDetailsMutation,
      UserUpdateOnboardingDetailsMutationVariables
    >(UpdateOnboardingDetailsMutationDocument, payload);
  }
}

export { GlobalCacheService };
