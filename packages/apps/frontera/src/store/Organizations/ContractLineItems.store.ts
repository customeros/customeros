import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { ContractLineItemStore } from '@store/Organizations/ContractLineItem.store.ts';

import { Contract, ContractInput, ServiceLineItem } from '@graphql/types';

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
      channelName: `ContractLineItem:${this.root.session.value.tenant}`,
      getItemId: (item: ServiceLineItem) => item?.metadata?.id,
      ItemStore: ContractLineItemStore,
    });
  }

  create = async (payload: ContractInput) => {
    const newContract = new ContractStore(this.root, this.transport);
    const tempId = newContract.value.metadata.id;
    const { name, organizationId, ...rest } = payload;

    if (payload) {
      merge(newContract.value, {
        contractName: payload.name,
        ...rest,
      });
    }

    this.value.set(tempId, newContract);

    try {
      const { contract_Create } = await this.transport.graphql.request<
        CREATE_CONTRACT_RESPONSE,
        CREATE_CONTRACT_PAYLOAD
      >(CREATE_CONTRACT_MUTATION, {
        input: {
          ...payload,
        },
      });
      runInAction(() => {
        this.value.delete(tempId);
        const serverId = contract_Create.metadata.id;

        newContract.value.metadata.id = serverId;
        this.value.set(serverId, newContract);

        this.sync({ action: 'APPEND', ids: [serverId] });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.value.delete(tempId);
      });
    }
  };
}

type CREATE_CONTRACT_PAYLOAD = {
  input: ContractInput;
};
type CREATE_CONTRACT_RESPONSE = {
  contract_Create: {
    metadata: {
      id: string;
    };
  };
};
const CREATE_CONTRACT_MUTATION = gql`
  mutation createContract($input: ContractInput!) {
    contract_Create(input: $input) {
      metadata {
        id
      }
    }
  }
`;
