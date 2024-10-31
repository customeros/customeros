import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { runInAction, makeAutoObservable } from 'mobx';
import { makeAutoSyncableGroup } from '@store/group-store';
import { FlowContactsService } from '@store/FlowContacts/__service__';

import { DataSource, FlowContact, FlowParticipantStatus } from '@graphql/types';

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
    return this.root.contacts.value.get(this.value.contact.metadata.id);
  }

  setId(id: string) {
    this.value.metadata.id = id;
  }

  async invalidate() {}

  public removeFlowContact = async () => {
    return this.service.deleteFlowContact({
      id: this.id,
    });
  };

  public deleteFlowContact = async () => {
    this.isLoading = true;

    const contactStore = this.contact;
    const flowName = this.contact?.flow?.value.name;
    const flowId = this.contact?.flow?.value.metadata.id ?? '';

    try {
      await this.removeFlowContact();
      runInAction(() => {
        contactStore?.update(
          (c) => {
            c.flows = [];

            return c;
          },
          { mutate: false },
        );

        this.root.ui.toastSuccess(
          `Contact removed from '${flowName}'`,
          'unlink-contact-from-flow-success',
        );
        this.root.contacts.sync({
          action: 'INVALIDATE',
          ids: [this.contactId],
        });

        this.root.flows.sync({
          action: 'INVALIDATE',
          ids: [flowId],
        });
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

const getDefaultValue = (): FlowContact => ({
  metadata: {
    source: DataSource.Openline,
    appSource: DataSource.Openline,
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    sourceOfTruth: DataSource.Openline,
  },
  status: FlowParticipantStatus.Scheduled,
  scheduledAction: '',
  scheduledAt: new Date().toISOString(),
  contact: {
    id: crypto.randomUUID(),
    createdAt: '',
    customFields: [],
    emails: [],
    firstName: '',
    jobRoles: [],
    lastName: '',
    locations: [],
    phoneNumbers: [],
    profilePhotoUrl: '',
    enrichDetails: {},
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
    timelineEvents: [],
    timelineEventsTotalCount: 0,
    updatedAt: '',
    appSource: DataSource.Openline,
    description: '',
    prefix: '',
    name: '',
    owner: null,
    tags: [],
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
