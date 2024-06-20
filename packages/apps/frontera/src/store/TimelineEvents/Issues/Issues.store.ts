import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Issue } from '@graphql/types';

import { IssueStore } from './Issue.store';

export class IssuesStore implements GroupStore<Issue> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, IssueStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<Issue>();
  subscribe = makeAutoSyncableGroup.subscribe;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Issues',
      getItemId: (item) => item.id,
      ItemStore: IssueStore,
    });
  }

  getByOrganizationId(id: string): IssueStore[] {
    return this.root.timelineEvents
      .getByOrganizationId(id)
      ?.filter((item) => item.value.__typename === 'Issue') as IssueStore[];
  }
}
