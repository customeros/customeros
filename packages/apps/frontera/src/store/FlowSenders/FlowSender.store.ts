import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';
import { FlowSendersService } from '@store/FlowSenders/__service__';

import { DataSource, FlowSender } from '@graphql/types';

export class FlowSenderStore implements Store<FlowSender> {
  value: FlowSender = getDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<FlowSender>();
  update = makeAutoSyncable.update<FlowSender>();
  private service: FlowSendersService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'FlowSender',
      getId: (d: FlowSender) => d?.metadata?.id,
    });
    makeAutoObservable(this);

    this.service = FlowSendersService.getInstance(transport);
  }

  get id() {
    return this.value.metadata?.id;
  }

  get userId() {
    return this.value.user?.id;
  }

  get user() {
    return this.value.user?.id
      ? this.root.users.value.get(this.value.user.id)
      : null;
  }

  setId(id: string) {
    this.value.metadata.id = id;
  }

  init(data: FlowSender): FlowSender {
    this.value = data;

    return data;
  }

  async invalidate() {}
}

export const getDefaultValue = (): FlowSender => ({
  metadata: {
    source: DataSource.Openline,
    appSource: DataSource.Openline,
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    sourceOfTruth: DataSource.Openline,
  },
  user: {
    id: '',
    name: '',
    firstName: '',
    lastName: '',
    emails: [],
    appSource: '',
    createdAt: '',
    updatedAt: '',
    bot: false,
    test: false,
    calendars: [],
    roles: [],
    internal: false,
    jobRoles: [],
    phoneNumbers: [],
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
    hasLinkedInToken: false,
    profilePhotoUrl: '',
    timezone: '',
    mailboxes: [],
  },
});
