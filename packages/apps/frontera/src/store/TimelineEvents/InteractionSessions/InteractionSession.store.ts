import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { DataSource, InteractionSession } from '@graphql/types';

export class InteractionSessionStore implements Store<InteractionSession> {
  value: InteractionSession = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<InteractionSession>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<InteractionSession>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'InteractionSession',
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

const defaultValue: InteractionSession = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'InteractionSession',
  attendedBy: [],
  describedBy: [],
  events: [],
  includes: [],
  name: '',
  sourceOfTruth: DataSource.Openline,
  status: '',
  updatedAt: new Date().toISOString(),
  channel: '',
  channelData: '',
  sessionIdentifier: '',
  type: '',
  // deprecated fields
  startedAt: new Date().toISOString(),
};
