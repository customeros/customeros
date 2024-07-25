import { Transport } from '@store/transport';

import SlackChannelsDocument from './slackChannels.graphql';
import TenantSettingsDocument from './tenantSettings.graphql';
import UpdateTenantSettingsDocument from './updateTenantSettings.graphql';
import UpdateOpportunityStageDocument from './updateOpportunityStage.graphql';
import {
  SlackChannelsQuery,
  SlackChannelsQueryVariables,
} from './slackChannels.generated';
import {
  TenantSettingsQuery,
  TenantSettingsQueryVariables,
} from './tenantSettings.generated';
import {
  UpdateTenantSettingsMutation,
  UpdateTenantSettingsMutationVariables,
} from './updateTenantSettings.generated';
import {
  UpdateOpportunityStageMutation,
  UpdateOpportunityStageMutationVariables,
} from './updateOpportunityStage.generated';

export class SettingsService {
  private static instance: SettingsService;
  private transport: Transport;

  private constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport) {
    if (!SettingsService.instance) {
      SettingsService.instance = new SettingsService(transport);
    }

    return SettingsService.instance;
  }

  async getSlackChannels(payload: SlackChannelsQueryVariables) {
    return this.transport.graphql.request<
      SlackChannelsQuery,
      SlackChannelsQueryVariables
    >(SlackChannelsDocument, payload);
  }

  async getTenantSettings() {
    return this.transport.graphql.request<
      TenantSettingsQuery,
      TenantSettingsQueryVariables
    >(TenantSettingsDocument);
  }

  async updateTenantSettings(payload: UpdateTenantSettingsMutationVariables) {
    return this.transport.graphql.request<
      UpdateTenantSettingsMutation,
      UpdateTenantSettingsMutationVariables
    >(UpdateTenantSettingsDocument, payload);
  }

  async updateOpportunityStage(
    payload: UpdateOpportunityStageMutationVariables,
  ) {
    return this.transport.graphql.request<
      UpdateOpportunityStageMutation,
      UpdateOpportunityStageMutationVariables
    >(UpdateOpportunityStageDocument, payload);
  }
}
