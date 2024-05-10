import { gql } from 'graphql-request';
import { makeAutoObservable } from 'mobx';

import { TenantSettings } from '@graphql/types';
import { TransportLayer } from '@store/transport';
import { RootStore } from '@store/root';

export class TenantStore {
  value: TenantSettings | null = null;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    this.load();
  }

  async load() {
    try {
      this.isLoading = true;
      const repsonse =
        await this.transportLayer.client.request<TENANT_SETTINGS_QUERY_RESULT>(
          TENANT_SETTINGS_QUERY,
        );
      this.value = repsonse.tenantSettings;
      this.isBootstrapped = true;
    } catch (err) {
      this.error = (err as Error).message;
    } finally {
      this.isLoading = false;
    }
  }
}

type TENANT_SETTINGS_QUERY_RESULT = {
  tenantSettings: TenantSettings;
};
const TENANT_SETTINGS_QUERY = gql`
  query TenantSettings {
    tenantSettings {
      logoUrl
      logoRepositoryFileId
      baseCurrency
      billingEnabled
    }
  }
`;
