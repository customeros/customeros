import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

import { TenantSettings } from '@graphql/types';

export class TenantStore {
  value: TenantSettings | null = null;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;

  constructor(public root: RootStore, public transportLayer: Transport) {
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
        await this.transportLayer.graphql.request<TENANT_SETTINGS_QUERY_RESULT>(
          TENANT_SETTINGS_QUERY,
        );
      runInAction(() => {
        this.value = repsonse.tenantSettings;
        this.isBootstrapped = true;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
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
