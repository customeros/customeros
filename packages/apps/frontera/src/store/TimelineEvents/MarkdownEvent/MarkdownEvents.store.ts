import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@infra/transport';
import { GroupOperation } from '@store/types';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { MarkdownEventType } from './types.ts';
import { MarkdownEventStore } from './MarkdownEvent.store.ts';

export class MarkdownEventsStore implements GroupStore<MarkdownEventType> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, MarkdownEventStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<MarkdownEventType>();
  subscribe = makeAutoSyncableGroup.subscribe;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'MarkdownEvents',
      getItemId: (item) => item.markdownEventMetadata.id,
      ItemStore: MarkdownEventStore,
    });
  }
}
