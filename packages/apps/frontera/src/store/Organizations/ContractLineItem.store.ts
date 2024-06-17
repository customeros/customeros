import { Channel } from 'phoenix';
import { match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { toJS, runInAction, makeAutoObservable } from 'mobx';

import { DateTimeUtils } from '@utils/date.ts';
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

  async invalidate() {
    try {
      this.isLoading = true;
      const { serviceLineItem } = await this.transport.graphql.request<
        CONTRACT_LINE_ITEM_QUERY_RESULT,
        { id: string }
      >(CONTRACT_LINE_ITEM_QUERY, { id: this.id });

      this.load(serviceLineItem);
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

  private transformHistoryToServiceLineItemUpdateInput = (
    history: Operation[],
  ): ServiceLineItemUpdateInput => {
    const serviceLineItemUpdate: Partial<ServiceLineItemUpdateInput> = {
      id: this.id,
      price: this.value.price,
      quantity: this.value.quantity,
      description: this.value.description,
      serviceStarted: this.value.serviceStarted,
      serviceEnded: this.value.serviceEnded,
      tax: {
        taxRate: this.value.tax.taxRate,
      },
    };
    console.log('🏷️ ----- toJS(this.history): ', toJS(this.history));

    history.forEach((change) => {
      change.diff.forEach((diffItem) => {
        const { path, val } = diffItem;
        const [fieldName, subField] = path;

        if (subField) {
          if (!serviceLineItemUpdate[fieldName]) {
            serviceLineItemUpdate[fieldName] = {};
          }
          (serviceLineItemUpdate[fieldName] as Record<string, unknown>)[
            subField
          ] = val;
        } else {
          (serviceLineItemUpdate as Record<string, unknown>)[fieldName] = val;
        }
      });
    });

    return serviceLineItemUpdate as ServiceLineItemUpdateInput;
  };

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
        this.invalidate();
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
        this.invalidate();
      });
    }
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const path = diff?.path;

    if (this.history.every((e) => !e.diff?.length)) {
      return;
    }
    console.log('🏷️ -----  : UPDATING SLI ');
    match(path).otherwise(() => {
      const payload = this.transformHistoryToServiceLineItemUpdateInput(
        this.history,
      );

      if (payload?.closed) {
        this.closeServiceLineItem({
          id: this.id,
        });

        return;
      }
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
  serviceStarted: DateTimeUtils.addDays(new Date().toISOString(), 1),
  tax: {
    salesTax: false,
    vat: false,
    taxRate: 0,
  },
};

type CONTRACT_LINE_ITEM_QUERY_RESULT = {
  serviceLineItem: ServiceLineItem;
};

const CONTRACT_LINE_ITEM_QUERY = gql`
  query ContractLineItem($id: ID!) {
    serviceLineItem(id: $id) {
      metadata {
        id
        created
        lastUpdated
        source
        appSource
        sourceOfTruth
      }
      description
      billingCycle
      price
      quantity
      comments
      serviceEnded
      parentId
      serviceStarted
      tax {
        salesTax
        vat
        taxRate
      }
    }
  }
`;

type SERVICE_LINE_UPDATE_PAYLOAD = {
  input: ServiceLineItemUpdateInput;
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

const SERVICE_LINE_CLOSE_MUTATION = gql`
  mutation contractLineItemCreateNewVersion(
    $input: ServiceLineItemCloseInput!
  ) {
    contractLineItem_Close(input: $input)
  }
`;
