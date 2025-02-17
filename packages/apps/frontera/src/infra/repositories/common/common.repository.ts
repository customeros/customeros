import { Transport } from '@infra/transport';

import { BrowserAutomationConfigResponse } from './common.datum';
import GetSlackChannelsDocument from './queries/getSlackChannels.graphql';
import { GetSlackChannelsQuery } from './queries/getSlackChannels.generated';
import TenantImpersonateListDocument from './queries/impersonateList.graphql';
import { TenantImpersonateListQuery } from '../common/queries/impersonateList.generated';

export type WebhookIntegration = 'grain' | 'fathom';

export class CommonRepository {
  static instance: CommonRepository | null = null;
  private transport = Transport.getInstance();

  public static getInstance() {
    if (!CommonRepository.instance) {
      CommonRepository.instance = new CommonRepository();
    }

    return CommonRepository.instance;
  }

  async getSlackChannels() {
    return this.transport.graphql.request<GetSlackChannelsQuery>(
      GetSlackChannelsDocument,
    );
  }

  async getTenantImpersonateList() {
    return this.transport.graphql.request<TenantImpersonateListQuery>(
      TenantImpersonateListDocument,
    );
  }

  async getWebhookUrl(
    integration: WebhookIntegration,
    { headers }: { headers?: Record<string, string> } = {},
  ) {
    return this.transport.http.get(
      `/tenant-cos-api/webhooks/v1/hooks?integration=${integration}`,
      {
        headers,
      },
    );
  }

  async createWebhookUrl(
    integration: WebhookIntegration,
    { headers }: { headers?: Record<string, string> } = {},
  ) {
    return this.transport.http.post(
      `/tenant-cos-api/webhooks/v1/hooks`,
      {
        integration,
      },
      {
        headers,
      },
    );
  }

  async getBrowserAutomationConfig({
    headers,
  }: { headers?: Record<string, string> } = {}) {
    return this.transport.http.get<BrowserAutomationConfigResponse>(
      '/bas/browser/config',
      {
        headers,
      },
    );
  }
}
