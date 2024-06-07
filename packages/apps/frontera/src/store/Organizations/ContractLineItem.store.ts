import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makePayload } from '@store/util.ts';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';

import {
  DataSource,
  BilledType,
  ServiceLineItem,
  ServiceLineItemCloseInput,
  ServiceLineItemUpdateInput,
} from '@graphql/types';

export class ContractLineItemStore implements Store<ServiceLineItem> {
  value: ServiceLineItem = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<ServiceLineItem>();
  update = makeAutoSyncable.update<ServiceLineItem>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'ContractLineItem',
      mutator: this.save,
      getId: (d: ServiceLineItem) => d?.metadata?.id,
    });
    makeAutoObservable(this);
  }
  get id() {
    return this.value.metadata.id;
  }
  set id(id: string) {
    this.value.metadata.id = id;
  }

  async invalidate() {}

  private async updateServiceLineItem(payload: ServiceLineItemUpdateInput) {
    try {
      this.isLoading = true;

      await this.transport.graphql.request<
        unknown,
        SERVICE_LINE_UPDATE_PAYLOAD
      >(SERVICE_LINE_UPDATE_MUTATION, {
        input: {
          ...payload,
          id: this.id,
        },
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
  }
  private async closeServiceLineItem(payload: ServiceLineItemCloseInput) {
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
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    // const type = diff?.op;
    const path = diff?.path;
    // const value = diff?.val;

    // TODO implement code to handle closing SLI
    match(path).otherwise(() => {
      const payload = makePayload<ServiceLineItemUpdateInput>(operation);
      this.updateServiceLineItem(payload);
    });
  }
}

const defaultValue: ServiceLineItem = {
  closed: false,
  externalLinks: [],
  metadata: {
    id: `new-${crypto.randomUUID()}`,
    appSource: DataSource.Openline,
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
  },
  description: '',
  billingCycle: BilledType.Monthly,
  price: 0,
  quantity: 0,
  comments: '',
  serviceEnded: null,
  parentId: '',
  serviceStarted: new Date().toISOString(),
  tax: {
    salesTax: false,
    vat: false,
    taxRate: 0,
  },
};

type SERVICE_LINE_UPDATE_PAYLOAD = {
  input: ServiceLineItemUpdateInput;
};
type SERVICE_LINE_UPDATE_RESPONSE = {
  contractLineItem_Update: ServiceLineItem;
};
const SERVICE_LINE_UPDATE_MUTATION = gql`
  mutation contractLineItemUpdate($input: ServiceLineItemUpdateInput!) {
    contractLineItem_Update(input: $input) {
      metadata {
        id
      }
    }
  }
`;

type SERVICE_LINE_CLOSE_PAYLOAD = {
  input: ServiceLineItemCloseInput;
};
type SERVICE_LINE_CLOSE_RESPONSE = {
  contractLineItem_Close: any;
};
const SERVICE_LINE_CLOSE_MUTATION = gql`
  mutation contractLineItemCreateNewVersion(
    $input: ServiceLineItemCloseInput!
  ) {
    contractLineItem_Close(input: $input)
  }
`;
