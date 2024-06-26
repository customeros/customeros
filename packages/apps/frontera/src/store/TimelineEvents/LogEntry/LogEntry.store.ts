import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
import { UserStore } from '@store/Users/User.store';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';

import { LogEntry, DataSource } from '@graphql/types';

import { LogEntriesService } from './__service__/LogEntries.service';

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
  private service: LogEntriesService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = LogEntriesService.getInstance(transport);

    makeAutoSyncable(this, {
      channelName: 'LogEntry',
      mutator: this.save,
      getId: (item) => item?.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {}
  async invalidate() {
    try {
      const { logEntry } = await this.service.getLogEntry(this.value.id);

      runInAction(() => {
        this.load(logEntry as LogEntry);
      });
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
