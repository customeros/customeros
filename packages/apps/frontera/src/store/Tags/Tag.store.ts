import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import { Tag, DataSource } from '@shared/types/__generated__/graphql.types';

export class TagStore implements Store<Tag> {
  value: Tag = defaultValue;
  version: number = 0;
  isLoading = false;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<Tag>();
  update = makeAutoSyncable.update<Tag>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'Tag',
      getId: (item) => item?.id,
    });
  }

  get tagName() {
    return this.value.name;
  }

  set id(id: string) {
    this.value.id = id;
  }

  async bootstrap() {}

  async invalidate() {}
}

const defaultValue: Tag = {
  id: crypto.randomUUID(),
  name: '',
  source: DataSource.Na,
  createdAt: '',
  appSource: '',
  updatedAt: '',
  metadata: {
    id: crypto.randomUUID(),
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
    appSource: 'organization',
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
  },
};
