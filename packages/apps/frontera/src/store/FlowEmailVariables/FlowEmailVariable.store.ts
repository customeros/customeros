import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import { EmailVariableEntity, EmailVariableEntityType } from '@graphql/types';

export class FlowEmailVariableStore implements Store<EmailVariableEntity> {
  value: EmailVariableEntity = getDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
  }

  get id() {
    return this.value.type;
  }

  setId(id: EmailVariableEntityType) {
    this.value.type = id;
  }

  get variables() {
    return this.value.variables;
  }

  invalidate(): Promise<void> {
    return Promise.resolve(undefined);
  }

  load(data: EmailVariableEntity): Promise<void> {
    this.value = data;

    return Promise.resolve();
  }

  update(): void {}
}

const getDefaultValue = (): EmailVariableEntity => ({
  type: EmailVariableEntityType.Contact,
  variables: [],
});
