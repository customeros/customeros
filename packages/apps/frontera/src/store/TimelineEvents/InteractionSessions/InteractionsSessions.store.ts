import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { InteractionSession } from '@graphql/types';

import { InteractionSessionStore } from './InteractionSession.store';

export class InteractionSessionsStore
  implements GroupStore<InteractionSession>
{
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, InteractionSessionStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<InteractionSession>();
  subscribe = makeAutoSyncableGroup.subscribe;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'InteractionSessions',
      getItemId: (item) => item.id,
      ItemStore: InteractionSessionStore,
    });
  }
}
