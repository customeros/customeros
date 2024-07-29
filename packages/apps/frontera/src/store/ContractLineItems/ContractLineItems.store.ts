import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Store } from '@store/store.ts';
import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { GroupOperation } from '@store/types.ts';
import { runInAction, makeAutoObservable } from 'mobx';
import { ContractStore } from '@store/Contracts/Contract.store.ts';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store.ts';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store.ts';

import { DateTimeUtils } from '@utils/date.ts';
import {
  ServiceLineItem,
  ServiceLineItemInput,
  ServiceLineItemCloseInput,
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

  createNewVersion = async (payload: ServiceLineItem, contractId: string) => {
    const newCli = new ContractLineItemStore(this.root, this.transport);
    const tempId = payload.metadata.id;
    if (payload) {
      merge(newCli.value, payload);
    }
    let serverId = '';
    const formatPayload: ServiceLineItemNewVersionInput = {
      tax: {
        taxRate: payload.tax.taxRate,
      },
      id: payload.parentId,
      price: payload.price,
      quantity: payload.quantity,
      description: payload.description,
      serviceStarted: payload.serviceStarted,
    };

    try {
      const { contractLineItem_NewVersion } =
        await this.transport.graphql.request<
          SERVICE_LINE_CREATE_NEW_VERSION_RESPONSE,
          SERVICE_LINE_CREATE_NEW_VERSION_PAYLOAD
        >(SERVICE_LINE_CREATE_NEW_VERSION_MUTATION, {
          input: {
            ...formatPayload,
          },
        });
      runInAction(() => {
        serverId = contractLineItem_NewVersion.metadata.id;
        newCli.value.metadata.id = serverId;
        const contract = this.root.contracts.value.get(contractId)?.value;

        if (contract) {
          const filteredContractLineItems =
            contract.contractLineItems?.filter(
              (e) => e.metadata.id !== tempId,
            ) ?? [];

          contract.contractLineItems = [
            ...filteredContractLineItems,
            newCli.value,
          ];
        }
        this.value.set(serverId, newCli);
        this.value.delete(tempId);

        this.sync({
          action: 'APPEND',
          ids: [serverId],
        });
      });
    } catch (err) {
      runInAction(() => {
        this.root.contractLineItems.value.delete(tempId);
        this.root.ui.toastError(
          `Failed to create '${payload.description}' service line item`,
          'failed-to-create-service-line-item',
        );
        this.error = (err as Error).message;
      });
    } finally {
      if (serverId) {
        setTimeout(() => {
          runInAction(() => {
            this.root.contractLineItems.sync({
              action: 'INVALIDATE',
              ids: [serverId],
            });
          });
        }, 1000);
      }
    }
  };
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
    // TODO clean up needed
    const newContractLineItem = new ContractLineItemStore(
      this.root,
      this.transport,
    );
    const tempId = `new-${crypto.randomUUID()}`;
    const contractStore = this.root.contracts.value.get(
      payload.contractId,
    ) as ContractStore;
    if (!(payload as ServiceLineItemNewVersionInput)?.id) {
      if (payload) {
        merge(newContractLineItem.tempValue, {
          ...payload,
          metadata: { id: tempId },
        });
      }

      this.value.set(tempId, newContractLineItem);
      const cli = contractStore?.tempValue;
      if (cli) {
        cli.contractLineItems = [
          ...(cli.contractLineItems ?? []),
          newContractLineItem?.tempValue,
        ];
      }
    } else if (this.isServiceLineItemNewVersionInput(payload) && payload.id) {
      const prevVersions = this.toArray()
        .filter(
          (e) =>
            (e as ContractLineItemStore)?.tempValue.parentId === payload.id ||
            (e as ContractLineItemStore)?.tempValue.metadata.id === payload.id,
        )
        .sort(
          (a, b) =>
            new Date(
              (a as ContractLineItemStore)?.tempValue?.serviceStarted ?? 0,
            ).getTime() -
            new Date(
              (b as ContractLineItemStore)?.tempValue?.serviceStarted ?? 0,
            ).getTime(),
        );
      const prevVersion = (
        prevVersions[prevVersions.length - 1] as ContractLineItemStore
      )?.tempValue;

      const serviceStarted =
        !prevVersion?.serviceStarted ||
        DateTimeUtils.isPast(prevVersion.serviceStarted)
          ? new Date().toISOString()
          : prevVersion.serviceStarted;
      merge(newContractLineItem.tempValue, {
        ...prevVersion,
        ...payload,
        serviceStarted,
        parentId: payload.id,
        metadata: { id: tempId },
      });
      this.value.set(tempId, newContractLineItem);

      const cli = contractStore?.tempValue;
      if (cli) {
        cli.contractLineItems = [
          ...(cli.contractLineItems ?? []),
          newContractLineItem?.tempValue,
        ];
      }
    }
  };

  closeServiceLineItem = async (payload: ServiceLineItemCloseInput) => {
    try {
      this.isLoading = true;

      await this.transport.graphql.request<unknown, SERVICE_LINE_CLOSE_PAYLOAD>(
        SERVICE_LINE_CLOSE_MUTATION,
        {
          input: {
            ...payload,
          },
        },
      );
      runInAction(() => {
        this.root.contractLineItems.value.get(payload.id)?.update(
          (prev) => ({
            ...prev,
            serviceEnded: new Date().toISOString(),
            closed: true,
          }),
          { mutate: false },
        );
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };

  createNewServiceLineItem = async (
    payload: ServiceLineItem,
    contractId: string,
  ) => {
    const newCli = new ContractLineItemStore(this.root, this.transport);
    const tempId = payload.metadata.id;
    if (payload) {
      merge(newCli.value, payload);
    }
    let serverId = '';
    try {
      const { contractLineItem_Create } = await this.transport.graphql.request<
        SERVICE_LINE_CREATE_RESPONSE,
        SERVICE_LINE_CREATE_PAYLOAD
      >(SERVICE_LINE_CREATE_MUTATION, {
        input: {
          tax: {
            taxRate: payload.tax.taxRate,
          },
          contractId,
          billingCycle: payload.billingCycle,
          price: payload.price,
          quantity: payload.quantity,
          serviceEnded: payload.serviceEnded,
          description: payload.description,
          serviceStarted: payload.serviceStarted,
        },
      });

      runInAction(() => {
        serverId = contractLineItem_Create.metadata.id;
        newCli.value.metadata.id = serverId;
        const contract = this.root.contracts.value.get(contractId)?.value;

        if (contract) {
          const filteredContractLineItems =
            contract.contractLineItems?.filter(
              (e) => e.metadata.id !== tempId,
            ) ?? [];

          contract.contractLineItems = [
            ...filteredContractLineItems,
            newCli.value,
          ];
        }
        this.value.set(serverId, newCli);
        this.value.delete(tempId);

        this.sync({ action: 'APPEND', ids: [serverId] });
      });
    } catch (err) {
      runInAction(() => {
        this.root.contractLineItems.value.delete(tempId);
        this.root.ui.toastError(
          `Failed to create '${payload.description}' service line item`,
          'failed-to-create-service-line-item',
        );
        this.error = (err as Error).message;
      });
    }
  };
}

type SERVICE_LINE_CREATE_PAYLOAD = {
  input: ServiceLineItemInput;
};
type SERVICE_LINE_CREATE_RESPONSE = {
  contractLineItem_Create: ServiceLineItem;
};
const SERVICE_LINE_CREATE_MUTATION = gql`
  mutation contractLineItemCreate($input: ServiceLineItemInput!) {
    contractLineItem_Create(input: $input) {
      metadata {
        id
      }
    }
  }
`;

type SERVICE_LINE_CREATE_NEW_VERSION_PAYLOAD = {
  input: ServiceLineItemNewVersionInput;
};
type SERVICE_LINE_CREATE_NEW_VERSION_RESPONSE = {
  contractLineItem_NewVersion: ServiceLineItem;
};
const SERVICE_LINE_CREATE_NEW_VERSION_MUTATION = gql`
  mutation contractLineItemCreateNewVersion(
    $input: ServiceLineItemNewVersionInput!
  ) {
    contractLineItem_NewVersion(input: $input) {
      metadata {
        id
      }
    }
  }
`;
type SERVICE_LINE_CLOSE_PAYLOAD = {
  input: ServiceLineItemCloseInput;
};

const SERVICE_LINE_CLOSE_MUTATION = gql`
  mutation contractLineItemClose($input: ServiceLineItemCloseInput!) {
    contractLineItem_Close(input: $input)
  }
`;
