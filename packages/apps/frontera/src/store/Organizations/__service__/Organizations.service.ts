import { Transport } from '@store/transport';

import UpdateOnboardingStatusDocument from './updateOnboardingStatus.graphql';
import {
  UpdateOnboardingStatusMutation,
  UpdateOnboardingStatusMutationVariables,
} from './updateOnboardingStatus.generated';

export class OrganizationsService {
  private static instance: OrganizationsService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  public static getInstance(transport: Transport): OrganizationsService {
    if (!OrganizationsService.instance) {
      OrganizationsService.instance = new OrganizationsService(transport);
    }

    return OrganizationsService.instance;
  }

  async updateOnboardingStatus(
    payload: UpdateOnboardingStatusMutationVariables,
  ) {
    return this.transport.graphql.request<
      UpdateOnboardingStatusMutation,
      UpdateOnboardingStatusMutationVariables
    >(UpdateOnboardingStatusDocument, payload);
  }
}