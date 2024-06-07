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
import { ContractStore } from '@store/Organizations/Contract.store.ts';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { ContractLineItemStore } from '@store/Organizations/ContractLineItem.store.ts';

import {
  ServiceLineItem,
  ContractUpdateInput,
  ContractRenewalInput,
  ServiceLineItemInput,
  ServiceLineItemUpdateInput,
  ServiceLineItemBulkUpdateInput,
  ServiceLineItemNewVersionInput,
} from '@graphql/types';

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

  private createNewVersion = async (
    payload: ServiceLineItemNewVersionInput,
  ) => {
    try {
      const { contractLineItem_NewVersion } =
        await this.transport.graphql.request<
          SERVICE_LINE_CREATE_NEW_VERSION_RESPONSE,
          SERVICE_LINE_CREATE_NEW_VERSION_PAYLOAD
        >(SERVICE_LINE_CREATE_NEW_VERSION_MUTATION, {
          input: {
            ...payload,
          },
        });
      runInAction(() => {
        console.log('🏷️ ----- : ', contractLineItem_NewVersion);
        // serverId = contract_Create.metadata.id;
        //
        // newContract.value.metadata.id = serverId;
        //
        // this.value.set(serverId, newContract);
        // this.value.delete(tempId);
        //
        // this.sync({ action: 'APPEND', ids: [serverId] });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      // if (serverId) {
      //   setTimeout(() => {
      //     runInAction(() => {
      //       this.root.organizations.value.get(organizationId)?.invalidate();
      //       this.value.get(serverId)?.invalidate();
      //
      //       this.root.organizations.sync({
      //         action: 'INVALIDATE',
      //         ids: [organizationId],
      //       });
      //     });
      //   }, 500);
      // }
    }
  };

  private createNewServiceLineItem = async (payload: ServiceLineItemInput) => {
    try {
      const { contractLineItem_Create } = await this.transport.graphql.request<
        SERVICE_LINE_CREATE_RESPONSE,
        SERVICE_LINE_CREATE_PAYLOAD
      >(SERVICE_LINE_CREATE_MUTATION, {
        input: {
          ...payload,
        },
      });
      runInAction(() => {
        // serverId = contract_Create.metadata.id;
        //
        // newContract.value.metadata.id = serverId;
        //
        // this.value.set(serverId, newContract);
        // this.value.delete(tempId);
        //
        // this.sync({ action: 'APPEND', ids: [serverId] });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      // if (serverId) {
      //   setTimeout(() => {
      //     runInAction(() => {
      //       this.root.organizations.value.get(organizationId)?.invalidate();
      //       this.value.get(serverId)?.invalidate();
      //
      //       this.root.organizations.sync({
      //         action: 'INVALIDATE',
      //         ids: [organizationId],
      //       });
      //     });
      //   }, 500);
      // }
    }
  };

  private isServiceLineItemInput(
    payload: ServiceLineItemNewVersionInput | ServiceLineItemInput,
  ): payload is ServiceLineItemInput {
    return (payload as ServiceLineItemInput).contractId !== undefined;
  }

  private isServiceLineItemNewVersionInput(
    payload: ServiceLineItemNewVersionInput | ServiceLineItemInput,
  ): payload is ServiceLineItemNewVersionInput {
    return (payload as ServiceLineItemNewVersionInput).id !== undefined;
  }
  create = async (
    payload: ServiceLineItemNewVersionInput | ServiceLineItemInput,
  ) => {
    if (this.isServiceLineItemInput(payload)) {
      await this.createNewServiceLineItem(payload);
    } else if (this.isServiceLineItemNewVersionInput(payload)) {
      await this.createNewVersion(payload);
    }
  };
}

type SERVICE_LINE_CREATE_PAYLOAD = {
  input: ServiceLineItemUpdateInput;
};
type SERVICE_LINE_CREATE_RESPONSE = {
  contractLineItem_Create: any;
};
const SERVICE_LINE_CREATE_MUTATION = gql`
  mutation contractLineItemCreate($input: ServiceLineItemInput!) {
    contractLineItem_Create(input: $input)
  }
`;

type SERVICE_LINE_CREATE_NEW_VERSION_PAYLOAD = {
  input: ServiceLineItemNewVersionInput;
};
type SERVICE_LINE_CREATE_NEW_VERSION_RESPONSE = {
  contractLineItem_NewVersion: any;
};
const SERVICE_LINE_CREATE_NEW_VERSION_MUTATION = gql`
  mutation contractLineItemCreateNewVersion(
    $input: ServiceLineItemNewVersionInput!
  ) {
    contractLineItem_NewVersion(input: $input)
  }
`;
