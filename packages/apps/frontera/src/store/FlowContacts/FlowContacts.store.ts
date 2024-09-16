import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { FlowContactsService } from '@store/FlowContacts/__service__';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';
import { FlowContactStore } from '@store/FlowContacts/FlowContact.store.ts';

import { FlowContact } from '@graphql/types';

export class FlowContactsStore implements GroupStore<FlowContact> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<FlowContact>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<FlowContact>();
  totalElements = 0;
  private service: FlowContactsService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'FlowContacts',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: FlowContactStore,
    });
    this.service = FlowContactsService.getInstance(transport);
  }

  getByContactId(contactId: string): Store<FlowContact> | undefined {
    return Array.from(this.value.values()).find(
      (flowContactStore) =>
        flowContactStore.value.contact?.metadata?.id === contactId,
    );
  }

  getByContactIds(contactIds: string[]): Store<FlowContact>[] {
    const contactIdSet = new Set(contactIds);

    return Array.from(this.value.values()).filter(
      (flowContactStore) =>
        flowContactStore.value.contact?.metadata?.id &&
        contactIdSet.has(flowContactStore.value.contact.metadata.id),
    );
  }
}
