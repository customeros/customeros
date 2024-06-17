import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { UserStore } from '@store/Users/User.store';
import { Store, makeAutoSyncable } from '@store/store';

import { Action, ActionType, DataSource } from '@graphql/types';

export class ActionStore implements Store<Action> {
  value: Action = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<Action>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<Action>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Action',
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

const defaultValue: Action = {
  id: crypto.randomUUID(),
  actionType: ActionType.Created,
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'Action',
  content: '',
  createdBy: UserStore.getDefaultValue(),
  metadata: '',
};
