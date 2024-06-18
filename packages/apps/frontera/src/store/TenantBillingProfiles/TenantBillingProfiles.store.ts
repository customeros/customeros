import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Store } from '@store/store.ts';
import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { GroupOperation } from '@store/types.ts';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store.ts';

import { TenantBillingProfile } from '@graphql/types';

import { TenantBillingProfileStore } from './TenantBillingProfile.store.ts';

export class TenantBillingProfilesStore
  implements GroupStore<TenantBillingProfile>
{
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<TenantBillingProfile>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<TenantBillingProfile>();
  totalElements = 0;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'TenantBillingProfiles',
      getItemId: (item) => item?.id,
      ItemStore: TenantBillingProfileStore,
    });
  }
  toArray() {
    return Array.from(this.value.values());
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    try {
      const { tenantBillingProfiles } =
        await this.transport.graphql.request<TENANT_BILLING_PROFILES_RESPONSE>(
          TENANT_BILLING_PROFILES_QUERY,
        );

      runInAction(() => {
        this.load(tenantBillingProfiles);
      });
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = tenantBillingProfiles.length;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}

type TENANT_BILLING_PROFILES_RESPONSE = {
  tenantBillingProfiles: TenantBillingProfile[];
};
const TENANT_BILLING_PROFILES_QUERY = gql`
  query getTenantBillingProfiles {
    tenantBillingProfiles {
      id
      createdAt
      updatedAt
      source
      sourceOfTruth
      appSource
      phone
      addressLine1
      addressLine2
      addressLine3
      locality
      country
      region
      zip
      legalName
      vatNumber
      sendInvoicesFrom
      sendInvoicesBcc
      canPayWithBankTransfer
      canPayWithPigeon
      check
      email
    }
  }
`;
