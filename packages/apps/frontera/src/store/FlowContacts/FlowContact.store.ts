import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';
import { FlowContactsService } from '@store/FlowContacts/__service__';

import { DataSource, FlowContact } from '@graphql/types';

export class FlowContactStore implements Store<FlowContact> {
  value: FlowContact = getDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<FlowContact>();
  update = makeAutoSyncable.update<FlowContact>();
  private service: FlowContactsService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'FlowContact',
      getId: (d: FlowContact) => d?.metadata?.id,
    });
    makeAutoObservable(this);

    this.service = FlowContactsService.getInstance(transport);
  }

  get id() {
    return this.value.metadata?.id;
  }

  get contactId() {
    return this.value.contact?.metadata?.id;
  }

  get contact() {
    return this.root.contacts.value.get(this.value.contact?.metadata?.id);
  }

  setId(id: string) {
    this.value.metadata.id = id;
  }

  async invalidate() {}

  public deleteFlowContact = async () => {
    this.isLoading = true;

    try {
      await this.service.deleteFlowContact({
        id: this.id,
      });

      runInAction(() => {
        const contactStore = this.contact;

        contactStore?.update(
          (c) => {
            c.flows = [];

            return c;
          },
          { mutate: false },
        );
        this.root.ui.toastSuccess(
          `Contact removed from '${this.value.contact.name}'`,
          'unlink-contact-from-sequence-success',
        );
        this.root.contacts.sync({
          action: 'INVALIDATE',
          ids: [this.contactId],
        });

        this.root.flows.invalidate();
      });
    } catch (e) {
      runInAction(() => {
        this.root.ui.toastError(
          `We couldn't remove a contact from a sequence`,
          'unlink-contact-from-sequence-error',
        );
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };
}

const getDefaultValue = (): FlowContact => ({
  metadata: {
    source: DataSource.Openline,
    appSource: DataSource.Openline,
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    sourceOfTruth: DataSource.Openline,
  },
  contact: {
    id: crypto.randomUUID(),
    createdAt: '',
    customFields: [],
    emails: [],
    fieldSets: [],
    firstName: '',
    jobRoles: [],
    lastName: '',
    locations: [],
    phoneNumbers: [],
    profilePhotoUrl: '',
    organizations: {
      content: [],
      totalPages: 0,
      totalElements: 0,
      totalAvailable: 0,
    },
    flows: [],
    socials: [],
    timezone: '',
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
    timelineEvents: [],
    timelineEventsTotalCount: 0,
    updatedAt: '',
    appSource: DataSource.Openline,
    description: '',
    prefix: '',
    name: '',
    owner: null,
    tags: [],
    template: null,
    connectedUsers: [],
    metadata: {
      source: DataSource.Openline,
      appSource: DataSource.Openline,
      id: crypto.randomUUID(),
      created: '',
      lastUpdated: new Date().toISOString(),
      sourceOfTruth: DataSource.Openline,
    },
  },
});
