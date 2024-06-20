import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { UserStore } from '@store/Users/User.store';
import { Store, makeAutoSyncable } from '@store/store';

import { Note, DataSource } from '@graphql/types';

export class NoteStore implements Store<Note> {
  value: Note = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<Note>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<Note>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Note',
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

const defaultValue: Note = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'Note',
  content: '',
  createdBy: UserStore.getDefaultValue(),
  includes: [],
  noted: [],
  sourceOfTruth: DataSource.Openline,
  updatedAt: new Date().toISOString(),
  contentType: '',
};
