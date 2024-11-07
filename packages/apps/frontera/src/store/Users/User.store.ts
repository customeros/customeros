import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { User, DataSource } from '@graphql/types';

export class UserStore implements Store<User> {
  value: User = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<User>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<User>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'User',
      mutator: this.save,
      getId: (user) => user?.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {}

  async invalidate() {}

  async save() {
    //
  }

  get id() {
    return this.value.id;
  }

  get name() {
    return (
      this.value?.name?.trim() ||
      `${this.value.firstName} ${this.value.lastName}`.trim()
    );
  }

  set id(id: string) {
    this.value.id = id;
  }

  static getDefaultValue() {
    return defaultValue;
  }
}

const defaultValue: User = {
  id: '',
  name: '',
  firstName: '',
  lastName: '',
  emails: [],
  appSource: '',
  createdAt: '',
  updatedAt: '',
  bot: false,
  calendars: [],
  roles: [],
  internal: false,
  test: false,
  jobRoles: [],
  hasLinkedInToken: false,
  phoneNumbers: [],
  source: DataSource.Openline,
  sourceOfTruth: DataSource.Openline,
  profilePhotoUrl: '',
  timezone: '',
  mailboxes: [],
};
