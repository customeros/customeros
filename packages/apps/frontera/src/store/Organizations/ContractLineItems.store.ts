import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { makePayload } from '@store/util.ts';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Operation, GroupOperation } from '@store/types';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { ContractLineItemStore } from '@store/Organizations/ContractLineItem.store.ts';

import {
  ContractInput,
  ContractStatus,
  ServiceLineItem,
  ContractUpdateInput,
  ContractRenewalInput,
  ServiceLineItemBulkUpdateInput,
} from '@graphql/types';

import { ContractStore } from './Contract.store';

export class ContractLineItemsStore implements GroupStore<ServiceLineItem> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<ServiceLineItem>> = new Map();
  organizationId: string = '';
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<ServiceLineItem>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: `ContractLineItems`,
      getItemId: (item: ServiceLineItem) => item?.metadata?.id,
      ItemStore: ContractLineItemStore,
    });
  }

  bulkUpdate = async (payload: ContractInput) => {
    try {
      this.isLoading = true;
      const { serviceLineItem_BulkUpdate } =
        await this.transport.graphql.request<
          SERVICE_LINE_UPDATE_BULK_RESPONSE,
          SERVICE_LINE_UPDATE_BULK_PAYLOAD
        >(SERVICE_LINE_UPDATE_BULK_MUTATION, {
          input: {
            ...payload,
          },
        });

      this.sync({ action: 'APPEND', ids: serviceLineItem_BulkUpdate });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };
}

type SERVICE_LINE_UPDATE_BULK_PAYLOAD = {
  input: ServiceLineItemBulkUpdateInput;
};
type SERVICE_LINE_UPDATE_BULK_RESPONSE = {
  serviceLineItem_BulkUpdate: [string];
};
const SERVICE_LINE_UPDATE_BULK_MUTATION = gql`
  mutation createContract($input: ServiceLineItemBulkUpdateInput!) {
    serviceLineItem_BulkUpdate(input: $input) {
      id
    }
  }
`;
