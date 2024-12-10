import { set, merge } from 'lodash';
import { Entity } from '@store/record';
import { Transport } from '@store/transport';
import { FlowStore } from '@store/Flows/Flow.store';
import { countryMap } from '@assets/countries/countriesMap';
import { action, computed, observable, runInAction } from 'mobx';

import {
  JobRole,
  DataSource,
  Organization,
  OrganizationWithJobRole,
} from '@shared/types/__generated__/graphql.types';

import { ContactsStore } from './Contacts2.store';
import { ContactService } from './__service__/Contacts.service';
import { ContactQuery } from './__service__/getContact.generated';

export type ContactDatum = NonNullable<ContactQuery['contact']>;

export class Contact extends Entity<ContactDatum> {
  private service: ContactService;
  @observable accessor value: ContactDatum = Contact.default();

  constructor(
    store: ContactsStore,
    data: ContactDatum,
    public transport?: Transport,
  ) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
    this.service = ContactService.getInstance(transport || store.transport);
  }

  @computed
  get isEnriching(): boolean {
    return (
      this.value?.enrichDetails?.requestedAt &&
      !this.value?.enrichDetails?.enrichedAt &&
      !this.value?.enrichDetails?.failedAt
    );
  }

  @computed
  get id() {
    return this.value.metadata.id;
  }

  set id(id: string) {
    this.value.metadata.id = id;
  }

  @computed
  get organizationId() {
    return this.value.organizations.content[0]?.metadata?.id;
  }

  get hasActiveOrganization() {
    const org = this.store.root.organizations.getById(this.organizationId);
  }

  @computed
  get organization() {
    return this.store.root.organizations.value.get(this.organizationId)?.value;
  }

  @computed
  get hasFlows() {
    return this.value.flows?.length > 0;
  }

  @computed
  get flows(): FlowStore[] | undefined {
    if (!this.value.flows?.length) return undefined;

    return this.value.flows.reduce((acc, flow) => {
      const flowStore = this.store.root.flows?.value.get(
        flow.metadata.id,
      ) as FlowStore;

      if (flowStore) {
        acc.push(flowStore);
      }

      return acc;
    }, [] as FlowStore[]);
  }

  @computed
  get flowsIds(): string[] | undefined {
    if (!this.flows?.length) return undefined;

    return this.flows.map((flow) => {
      return flow?.id;
    });
  }

  @computed
  get name() {
    return (
      this.value.name || `${this.value.firstName} ${this.value.lastName}`.trim()
    );
  }

  @computed
  get emailId() {
    return this.value.emails?.[0]?.id;
  }

  @computed
  get connectedUsers() {
    return this.value.connectedUsers.map(
      ({ id }) => this.store.root.users.value.get(id)?.value,
    );
  }

  @computed
  get country() {
    if (!this.value.locations?.[0]?.countryCodeA2) return undefined;

    return countryMap.get(this.value.locations[0].countryCodeA2.toLowerCase());
  }

  @action
  deletePersona(personaId: string) {
    this.value.tags = (this.value?.tags || []).filter(
      (tag) => tag.metadata.id !== personaId,
    );
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

      set(this.value.phoneNumbers?.[0], 'id', phoneNumberMergeToContact.id);
    } catch (e) {
      runInAction(() => {});
    }
  }

  async updatePhoneNumber() {
    const phoneNumber = this.value.phoneNumbers?.[0].rawPhoneNumber ?? '';

    try {
      await this.service.updatePhoneNumber({
        input: {
          id: this.value.phoneNumbers[0].id,
          phoneNumber,
        },
      });
    } catch (e) {
      runInAction(() => {});
    }
  }

  async addSocial(
    url: string,
    options?: { onSuccess?: (serverId: string) => void },
  ) {
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
      runInAction(() => {});
    } finally {
      options?.onSuccess?.(this.value.socials?.[0]?.id);
    }
  }

  async findEmail() {
    try {
      await this.service.findEmail({
        contactId: this.id,
        organizationId: this.organizationId,
      });
    } catch (e) {
      runInAction(() => {});
    }
  }

  async setPrimaryEmail(emailId: string) {
    const email = this.value.emails.find((email) => email.id === emailId);

    try {
      await this.service.setPrimaryEmail({
        contactId: this.id,
        email: email?.email || '',
      });
    } catch (e) {
      runInAction(() => {});
    } finally {
      this.invalidate();
    }
  }

  async removeTagFromContact(tagId: string) {
    try {
      await this.service.removeTagsFromContact({
        input: {
          contactId: this.id,
          tag: {
            id: tagId,
          },
        },
      });
    } catch (e) {
      runInAction(() => {});
    }
  }

  async removeAllTagsFromContact() {
    const tags =
      this.value?.tags?.map((tag) =>
        this.removeTagFromContact(tag.metadata.id),
      ) || [];

    try {
      await Promise.all(tags);

      runInAction(() => {
        this.value.tags = [];
        this.store.root.ui.toastSuccess(
          'All tags were removed',
          'tags-remove-success',
        );
      });
    } catch (e) {
      runInAction(() => {});
    }
  }

  static default(payload?: ContactDatum): ContactDatum {
    return merge(
      {
        id: crypto.randomUUID(),
        createdAt: '',
        customFields: [],
        emails: [],
        firstName: '',
        jobRoles: [],
        lastName: '',
        locations: [],
        latestOrganizationWithJobRole: {
          jobRole: {} as JobRole,
          organization: {} as Organization,
        } as OrganizationWithJobRole,

        phoneNumbers: [],
        profilePhotoUrl: '',
        organizations: {
          content: [],
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
          id: crypto.randomUUID(),
        },
        enrichDetails: {
          enrichedAt: '',
          failedAt: '',
          requestedAt: '',
        },
      },
      payload ?? {},
    );
  }
}
