import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { PageView, DataSource } from '@graphql/types';

export class PageViewStore implements Store<PageView> {
  value: PageView = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<PageView>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<PageView>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'PageView',
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

const defaultValue: PageView = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  source: DataSource.Openline,
  __typename: 'PageView',
  application: '',
  endedAt: new Date().toISOString(),
  startedAt: new Date().toISOString(),
  engagedTime: 0,
  orderInSession: 0,
  pageTitle: '',
  pageUrl: '',
  sessionId: '',
  sourceOfTruth: DataSource.Openline,
};
