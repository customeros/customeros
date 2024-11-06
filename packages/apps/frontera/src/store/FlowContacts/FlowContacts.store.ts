import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
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

  public deleteFlowContacts = async (ids: string[]) => {
    if (!ids.length) return;
    this.isLoading = true;

    const flowContacts = ids.map(
      (id) => this.value.get(id) as FlowContactStore,
    );

    const contactStores = flowContacts.map((fc) => fc?.contact);
    const flowStores = contactStores.flatMap((cs) => cs?.flows);

    try {
      await this.service.deleteFlowContacts({
        id: ids,
      });

      runInAction(() => {
        contactStores.forEach((c) => {
          c?.update(
            (c) => {
              c.flows = [];

              return c;
            },
            { mutate: false },
          );
        });
        flowStores.forEach((c) => {
          c?.update(
            (c) => {
              c.contacts = c.contacts.filter(
                (e) => !ids.includes(e.metadata.id),
              );

              return c;
            },
            { mutate: false },
          );
        });

        this.root.contacts.sync({
          action: 'INVALIDATE',
          ids: flowContacts.map((e) => e.contactId),
        });

        const flowsIds = flowStores
          .flatMap((c) => c?.id)
          .filter((id): id is string => typeof id === 'string');

        this.root.flows.sync({
          action: 'INVALIDATE',
          ids: flowsIds,
        });
        this.root.ui.toastSuccess(
          `Contacts removed from flows`,
          'unlink-contact-from-flow-success',
        );
      });
    } catch (e) {
      runInAction(() => {
        this.root.ui.toastError(
          `We couldn't remove a contact from a flow`,
          'unlink-contact-from-flow-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };
}
