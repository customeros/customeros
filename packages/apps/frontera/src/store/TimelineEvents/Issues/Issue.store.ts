import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import { Issue, DataSource } from '@graphql/types';

export class IssueStore implements Store<Issue> {
  value: Issue = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<Issue>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<Issue>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Issue',
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

const defaultValue: Issue & { issueStatus: string } = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'Issue',
  assignedTo: [],
  comments: [],
  externalLinks: [],
  followedBy: [],
  interactionEvents: [],
  sourceOfTruth: DataSource.Openline,
  status: '',
  issueStatus: '',
  updatedAt: new Date().toISOString(),
  description: '',
  priority: '',
  reportedBy: undefined,
  subject: '',
  submittedBy: undefined,
  tags: [],
};
