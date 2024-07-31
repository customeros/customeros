import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { Analysis, DataSource } from '@graphql/types';

export class AnalysisStore implements Store<Analysis> {
  value: Analysis = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<Analysis>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<Analysis>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Analysis',
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

  static getDefaultValue(): Analysis {
    return defaultValue;
  }
}

const defaultValue: Analysis = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'Analysis',
  content: '',
  describes: [],
  sourceOfTruth: DataSource.Openline,
  analysisType: '',
  contentType: '',
};
