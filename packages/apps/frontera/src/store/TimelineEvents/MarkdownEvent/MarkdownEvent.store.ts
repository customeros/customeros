import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@infra/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { DataSource } from '@graphql/types';
import { uuidv4 } from '@utils/generateUuid';

import { MarkdownEventType } from './types';

export class MarkdownEventStore implements Store<MarkdownEventType> {
  value: MarkdownEventType = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<MarkdownEventType>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<MarkdownEventType>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'MarkdownEvent',
      mutator: this.save,
      getId: (item) => item?.markdownEventMetadata?.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {}

  async invalidate() {}

  async save() {}

  get id() {
    return this.value.markdownEventMetadata.id;
  }

  set id(id: string) {
    this.value.markdownEventMetadata.id = id;
  }
}

const defaultValue: MarkdownEventType = {
  content: '',
  metadata: {
    id: '',
    created: '',
    lastUpdated: '',
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
    appSource: DataSource.Openline,
  },
  markdownEventMetadata: {
    id: uuidv4(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
    appSource: DataSource.Openline,
  },
};
