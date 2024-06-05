import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { DataSource, BilledType, ServiceLineItem } from '@graphql/types';

export class OpportunityStore implements Store<ServiceLineItem> {
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
      channelName: '',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
    makeAutoObservable(this);
  }

  async invalidate() {}

  set id(id: string) {
    this.value.metadata.id = id;
  }

  private async save() {
    // const payload: PAYLOAD = {
    //   input: {
    //     ...omit(this.value, 'metadata', 'owner'),
    //     contractId: this.value.metadata.id,
    //   },
    // };
    try {
      this.isLoading = true;
      // await this.transport.graphql.request(UPDATE_CONTRACT_DEF, payload);
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }
}

const defaultValue: ServiceLineItem = {
  closed: false,
  externalLinks: [],
  metadata: {
    id: '',
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
