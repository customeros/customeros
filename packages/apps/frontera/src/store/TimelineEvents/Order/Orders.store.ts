import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Order } from '@graphql/types';

import { OrderStore } from './Order.store';

export class OrdersStore implements GroupStore<Order> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, OrderStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<Order>();
  subscribe = makeAutoSyncableGroup.subscribe;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Orders',
      getItemId: (item) => item.id,
      ItemStore: OrderStore,
    });
  }
}
