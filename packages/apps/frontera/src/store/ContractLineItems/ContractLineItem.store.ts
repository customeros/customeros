import { isEqual } from 'lodash';
import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { getDiff } from 'recursive-diff';
import { RootStore } from '@store/root.ts';
import { Operation } from '@store/types.ts';
import { Transport } from '@store/transport.ts';
import { Store, makeAutoSyncable } from '@store/store.ts';
import { toJS, runInAction, makeAutoObservable } from 'mobx';
import { ContractLineItemService } from '@store/ContractLineItems/ContractLineItem.service.ts';

import { DateTimeUtils } from '@utils/date.ts';
import {
  DataSource,
  BilledType,
  ServiceLineItem,
  ServiceLineItemUpdateInput,
} from '@graphql/types';

export class ContractLineItemStore implements Store<ServiceLineItem> {
  value: ServiceLineItem = defaultValue;
  tempValue: ServiceLineItem = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<ServiceLineItem>();
  update = makeAutoSyncable.update<ServiceLineItem>();
  private service: ContractLineItemService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'ContractLineItem',
      getId: (d: ServiceLineItem) => d?.metadata?.id,
    });
    makeAutoObservable(this);
    this.service = ContractLineItemService.getInstance(transport);
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
  updateTemp(updater: (prev: ServiceLineItem) => ServiceLineItem) {
    const lhs = toJS(this.tempValue);
    const next = updater(this.tempValue);
    const rhs = toJS(next);
    const diff = getDiff(lhs, rhs, true);

    const operation: Operation = {
      id: this.version,
      diff,
      ref: this.transport.refId,
    };

    this.history.push(operation);
    this.tempValue = next;
  }

  async updateServiceLineItem() {
    const isEqualValue = isEqual(this.value, this.tempValue);
    if (isEqualValue) return;

    try {
      this.isLoading = true;

      await this.transport.graphql.request<
        unknown,
        SERVICE_LINE_UPDATE_PAYLOAD
      >(SERVICE_LINE_UPDATE_MUTATION, {
        input: {
          id: this.id,
          price: this.tempValue.price,
          quantity: this.tempValue.quantity,
          description: this.tempValue.description,
          serviceStarted: this.tempValue.serviceStarted,
          serviceEnded: this.tempValue.serviceEnded,
          tax: {
            taxRate: this.tempValue.tax.taxRate,
          },
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
