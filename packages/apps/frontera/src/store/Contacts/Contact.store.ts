import type { RootStore } from '@store/root';

import set from 'lodash/set';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';

import { Contact, DataSource, ContactUpdateInput } from '@graphql/types';

import { ContactService } from './Contact.service';

export class ContactStore implements Store<Contact> {
  value: Contact;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<Contact>();
  update = makeAutoSyncable.update<Contact>();
  private service: ContactService;
  organizationId: string = '';

  constructor(public root: RootStore, public transport: Transport) {
    this.value = getDefaultValue();

    makeAutoSyncable(this, {
      channelName: 'Contact',
      mutator: this.save,
      getId: (d) => d?.id,
    });
    makeAutoObservable(this);
    this.service = ContactService.getInstance(transport);
  }

  async invalidate() {}

  set id(id: string) {
    this.value.id = id;
  }
  get id() {
    return this.value.id;
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const type = diff?.op;
    const path = diff?.path;

    match(path)
      .with(['phoneNumbers', 0, ...P.array()], () => {
        if (type === 'add') {
          this.addPhoneNumber();
        }
        if (type === 'update') {
          this.updatePhoneNumber();
        }
      })
      .with(['socials', ...P.array()], ([_, index]) => {
        if (type === 'add') {
          this.addSocial();
        }
        if (type === 'update') {
          this.updateSocial(index as number);
        }
      })
      .with(['jobRoles', 0, ...P.array()], () => {
        if (type === 'add') {
          this.addJobRole();
        }
        if (type === 'update') {
          this.updateJobRole();
        }
      })
      .with(['emails', 0, ...P.array()], () => {
        if (type === 'add') {
          this.addEmail();
        }
        if (type === 'update') {
          this.updateEmail();
        }
      })
      .otherwise(() => {
        const payload = makePayload<ContactUpdateInput>(operation);
        this.updateContact(payload);
      });
  }

  async linkOrganization(organizationId: string) {
    runInAction(() => {
      this.organizationId = organizationId;
    });

    try {
      await this.service.linkOrganization({
        input: {
          contactId: this.value.id,
          organizationId,
        },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  private async updateContact(input: ContactUpdateInput) {
    try {
      await this.service.updateContact({
        input: { ...input, id: this.id, patch: true },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async addJobRole() {
    try {
      const { jobRole_Create } = await this.service.addJobRole({
        contactId: this.id,
        input: {
          organizationId: this.organizationId,
          description: this.value.jobRoles[0].description,
        },
      });

      runInAction(() => {
        set(this.value.jobRoles?.[0], 'id', jobRole_Create.id);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async updateJobRole() {
    try {
      await this.service.updateJobRole({
        contactId: this.id,
        input: {
          id: this.value.jobRoles[0].id,
          description: this.value.jobRoles[0].description,
        },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async addEmail() {
    const email = this.value.emails?.[0].email ?? '';

    try {
      const { emailMergeToContact } = await this.service.addContactEmail({
        contactId: this.id,
        input: {
          email,
        },
      });

      runInAction(() => {
        set(this.value.emails?.[0], 'id', emailMergeToContact.id);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async updateEmail() {
    const email = this.value.emails?.[0].email ?? '';

    try {
      await this.service.updateContactEmail({
        input: {
          id: this.value.emails[0].id,
          email,
        },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async removeEmail() {
    const email = this.value.emails?.[0].email ?? '';

    try {
      await this.service.removeContactEmail({
        contactId: this.id,
        email,
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async addPhoneNumber() {
    const phoneNumber = this.value.phoneNumbers?.[0].rawPhoneNumber ?? '';

    try {
      const { phoneNumberMergeToContact } = await this.service.addPhoneNumber({
        contactId: this.id,
        input: {
          phoneNumber,
        },
      });

      runInAction(() => {
        set(this.value.phoneNumbers?.[0], 'id', phoneNumberMergeToContact.id);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async updatePhoneNumber() {
    const phoneNumber = this.value.phoneNumbers?.[0].rawPhoneNumber ?? '';

    try {
      await this.service.updatePhoneNumber({
        contactId: this.id,
        input: {
          id: this.value.phoneNumbers[0].id,
          phoneNumber,
        },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async removePhoneNumber(id: string) {
    try {
      await this.service.removePhoneNumber({
        id,
        contactId: this.id,
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async addSocial() {
    const url = this.value.socials?.[0].url ?? '';

    try {
      const { contact_AddSocial } = await this.service.addSocial({
        contactId: this.id,
        input: {
          url,
        },
      });

      runInAction(() => {
        const serverId = contact_AddSocial.id;
        set(this.value.socials?.[0], 'id', serverId);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async updateSocial(index: number) {
    const social = this.value.socials?.[index];

    try {
      await this.service.updateSocial({
        input: {
          id: social.id,
          url: social.url,
        },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async findEmail() {
    try {
      await this.service.findEmail({
        contactId: this.id,
        organizationId: this.organizationId,
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }
}

const getDefaultValue = (): Contact => ({
  id: crypto.randomUUID(),
  createdAt: '',
  customFields: [],
  emails: [],
  fieldSets: [],
  firstName: '',
  jobRoles: [],
  lastName: '',
  locations: [],
  notes: {
    content: [],
    totalElements: 0,
    totalPages: 0,
  },
  phoneNumbers: [],
  profilePhotoUrl: '',
  notesByTime: [],
  organizations: {
    content: [],
    totalPages: 0,
    totalElements: 0,
    totalAvailable: 0,
  },
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
});