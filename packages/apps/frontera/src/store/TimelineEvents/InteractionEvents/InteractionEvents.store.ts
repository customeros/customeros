import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { InteractionEvent } from '@graphql/types';

import { InteractionEventStore } from './InteractionEvent.store';

export class InteractionEventsStore implements GroupStore<InteractionEvent> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, InteractionEventStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<InteractionEvent>();
  subscribe = makeAutoSyncableGroup.subscribe;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'InteractionEvents',
      getItemId: (item) => item.id,
      ItemStore: InteractionEventStore,
    });
  }
}
