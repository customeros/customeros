import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { Meeting, DataSource, MeetingStatus } from '@graphql/types';

export class MeetingStore implements Store<Meeting> {
  value: Meeting = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<Meeting>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<Meeting>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Meeting',
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

const defaultValue: Meeting = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'Meeting',
  attendedBy: [],
  createdBy: [],
  events: [],
  externalSystem: [],
  includes: [],
  note: [],
  sourceOfTruth: DataSource.Openline,
  status: MeetingStatus.Undefined,
  updatedAt: new Date().toISOString(),
  agenda: '',
  agendaContentType: '',
  conferenceUrl: '',
  endedAt: '',
  meetingExternalUrl: '',
  name: '',
  recording: undefined,
  startedAt: '',
};
