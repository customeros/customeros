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

import { DateTimeUtils } from '@utils/date.ts';
import {
  ServiceLineItem,
  ServiceLineItemInput,
  ServiceLineItemUpdateInput,
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
  toArray() {
    return Array.from(this.value.values());
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
    payload:
      | (ServiceLineItemNewVersionInput & { contractId: string })
      | ServiceLineItemInput,
  ) => {
    const newContractLineItem = new ContractLineItemStore(
      this.root,
      this.transport,
    );
    const tempId = `new-${crypto.randomUUID()}`;

    if (!payload?.id) {
      if (payload) {
        merge(newContractLineItem.value, {
          ...payload,
          metadata: { id: tempId },
        });
      }

      this.value.set(tempId, newContractLineItem);
      this.root.contracts.value.get(payload.contractId)?.update(
        (prev) => ({
          ...prev,
          contractLineItems: [
            ...(prev?.contractLineItems ?? []),
            newContractLineItem?.value,
          ],
        }),
        { mutate: false },
      );

      // await this.createNewServiceLineItem(payload);
    } else if (this.isServiceLineItemNewVersionInput(payload) && payload.id) {
      const prevVersions = this.toArray()
        .filter(
          (e) =>
            e.value.parentId === payload.id ||
            e.value.metadata.id === payload.id,
        )
        .sort(
          (a, b) =>
            new Date(a?.value?.serviceStarted) -
            new Date(b.value.serviceStarted),
        );
      const prevVersion = prevVersions[prevVersions.length - 1]?.value;

      merge(newContractLineItem.value, {
        ...prevVersion,
        ...payload,
        serviceStarted: DateTimeUtils.addDays(
          prevVersion?.serviceStarted ?? new Date().toISOString(),
          1,
        ),
        parentId: payload.id,
        metadata: { id: tempId },
      });
      this.value.set(tempId, newContractLineItem);

      this.root.contracts.value.get(payload.contractId)?.update(
        (prev) => ({
          ...prev,
          contractLineItems: [
            ...(prev?.contractLineItems ?? []),
            newContractLineItem?.value,
          ],
        }),
        { mutate: false },
      );

      // await this.createNewVersion(payload);
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
