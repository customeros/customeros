import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Analysis } from '@graphql/types';

import { AnalysisStore } from './Analysis.store';

export class AnalysesStore implements GroupStore<Analysis> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, AnalysisStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<Analysis>();
  subscribe = makeAutoSyncableGroup.subscribe;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Analyses',
      getItemId: (item) => item.id,
      ItemStore: AnalysisStore,
    });
  }
}
