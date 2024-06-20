import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { Order, DataSource } from '@graphql/types';

export class OrderStore implements Store<Order> {
  value: Order = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<Order>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<Order>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Order',
      mutator: this.save,
      getId: (item) => item?.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {}
  async invalidate() {}
  async save() {}

  get id() {
    return this.value.id;
  }
  set id(id: string) {
    this.value.id = id;
  }
}

const defaultValue: Order = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'Order',
  sourceOfTruth: DataSource.Openline,
  cancelledAt: null,
  confirmedAt: null,
  fulfilledAt: null,
  paidAt: null,
};
