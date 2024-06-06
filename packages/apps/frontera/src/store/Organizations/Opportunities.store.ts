import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { OpportunityStore } from '@store/Organizations/Opportunity.store.ts';

import { Opportunity, ServiceLineItem } from '@graphql/types';

export class OpportunitiesStore implements GroupStore<ServiceLineItem> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<Opportunity>> = new Map();
  organizationId: string = '';
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<Opportunity>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: `Opportunity:${this.root.session.value.tenant}`,
      getItemId: (item: Opportunity) => item?.id,
      ItemStore: OpportunityStore,
    });
  }
}
