import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { UserStore } from '@store/Users/User.store';
import { Store, makeAutoSyncable } from '@store/store';

import { LogEntry, DataSource } from '@graphql/types';

export class LogEntryStore implements Store<LogEntry> {
  value: LogEntry = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<LogEntry>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<LogEntry>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'LogEntry',
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

const defaultValue: LogEntry = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'LogEntry',
  externalLinks: [],
  sourceOfTruth: DataSource.Openline,
  updatedAt: new Date().toISOString(),
  tags: [],
  startedAt: new Date().toISOString(),
  content: '',
  contentType: '',
  createdBy: UserStore.getDefaultValue(),
};
