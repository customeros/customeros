import set from 'lodash/set';
import merge from 'lodash/merge';
import { Entity } from '@store/record';
import { Transport } from '@infra/transport';
import { FlowStore } from '@store/Flows/Flow.store';
import { JobRoleDatum } from '@store/JobRoles/JobRole.dto';
import { countryMap } from '@assets/countries/countriesMap';
import { action, computed, observable, runInAction } from 'mobx';

import { ContactsStore } from './Contacts.store';
import { ContactService } from './__service__/Contacts.service';
import { GetContactsByIdsQuery } from './__service__/getContactsById.generated';

export type ContactDatum = NonNullable<
  GetContactsByIdsQuery['ui_contacts'][number]
>;

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
    this.service = ContactService.getInstance();
  }

  @computed
  get isEnriching(): boolean {
    return (
      this.value.enrichedRequestedAt &&
      !this.value.enrichedAt &&
      !this.value.enrichedFailedAt
    );
  }

  @computed
  get emailEnriching(): boolean {
    return (
      this.value.enrichedEmailRequestedAt &&
      !this.value.enrichedEmailEnrichedAt &&
      !this.value.enrichedEmailFound
    );
  }

  @computed
  get id() {
    return this.value.id;
  }

  set id(id: string) {
    this.value.id = id;
  }

  @computed
  get organizationId() {
    return this.value.primaryOrganizationId;
  }

  @computed
  get organization() {
    return this.store.root.organizations.value.get(
      this.value.primaryOrganizationId || '',
    )?.value;
  }

  @computed
  get hasFlows() {
    return this.value.flows?.length > 0;
  }

  @computed
  get flows(): FlowStore[] | undefined {
    if (!this.value.flows?.length) return undefined;

    return this.value.flows.reduce((acc, flow) => {
      const flowStore = this.store.root.flows?.value.get(flow) as FlowStore;

      if (flowStore) {
        acc.push(flowStore);
      }

      return acc;
    }, [] as FlowStore[]);
  }

  @computed
  get flowsIds(): string[] {
    return this.value.flows ?? [];
  }

  @computed
  get name() {
    return (
      this.value.name || `${this.value.firstName} ${this.value.lastName}`.trim()
    );
  }

  @action
  setName(name: string) {
    this.draft();
    this.value.name = name;
    this.commit({ syncOnly: true });
  }

  @computed
  get emailId() {
    return this.value.emails?.[0]?.id;
  }

  @computed
  get connectedUsers() {
    return this.value.connectedUsers.map(
      (id) => this.store.root.users.value.get(id)?.value,
    );
  }

  @computed
  get jobRoles() {
    return this.value.jobRoleIds.reduce((acc, id) => {
      const record = this.store.root.jobRoles.getById(id);

      if (record) acc.push(record.value);

      return acc;
    }, [] as JobRoleDatum[]);
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

        set(this.value, 'linkedInInternalId', serverId);
      });
    } catch (e) {
      runInAction(() => {});
    } finally {
      options?.onSuccess?.(this.value.linkedInInternalId || '');
    }
  }

  async findEmail() {
    try {
      await this.service.findEmail({
        contactId: this.id,
        organizationId: this.organizationId || '',
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

  @action
  public changeName(name: string) {
    this.draft();
    this.value.name = name;
    this.commit({ syncOnly: true });
  }

  @action
  public addJobRole(jobRoleId: string) {
    this.draft();
    this.value.jobRoleIds.push(jobRoleId);
    this.commit({ syncOnly: true });
  }

  @action
  public addTag(id: string) {
    this.draft();

    if (!this.value.tags) {
      this.value.tags = [];
    }

    const tag = this.store.root.tags.getById(id);

    if (!tag) {
      console.error(`Contact.addTag: Tag with id ${id} not found`);

      return;
    }

    this.value.tags.push(tag.value);
    this.commit({ syncOnly: true });
  }

  @action
  public removeTag(id: string) {
    this.draft();

    if (!this.value.tags) {
      this.value.tags = [];
    }

    this.value.tags = this.value.tags.filter((tag) => tag.metadata.id !== id);
    this.commit({ syncOnly: true });
  }

  // @deprecated
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

  @action
  removeOrganization() {
    this.value.primaryOrganizationId = '';
    this.value.primaryOrganizationName = '';
    this.value.primaryOrganizationJobRoleId = '';
    this.value.primaryOrganizationJobRoleTitle = '';
    this.value.primaryOrganizationJobRoleDescription = '';
    this.value.primaryOrganizationJobRoleStartDate = '';
    this.value.primaryOrganizationJobRoleEndDate = '';
  }

  static default(payload?: Partial<ContactDatum>): ContactDatum {
    return merge(
      {
        id: crypto.randomUUID(),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        firstName: '',
        lastName: '',
        name: '',
        prefix: '',
        hide: false,
        description: '',
        timezone: '',
        profilePhotoUrl: '',
        enrichedAt: '',
        enrichedFailedAt: '',
        enrichedRequestedAt: '',
        enrichedEmailRequestedAt: '',
        enrichedEmailEnrichedAt: '',
        enrichedEmailFound: false,
        linkedInInternalId: '',
        linkedInUrl: '',
        linkedInAlias: '',
        linkedInExternalId: '',
        linkedInFollowerCount: 0,
        primaryOrganizationId: '',
        primaryOrganizationName: '',
        primaryOrganizationJobRoleId: '',
        primaryOrganizationJobRoleTitle: '',
        primaryOrganizationJobRoleDescription: '',
        primaryOrganizationJobRoleStartDate: '',
        primaryOrganizationJobRoleEndDate: '',
        emails: [],
        phones: [],
        tags: [],
        jobRoleIds: [],
        locations: [],
        connectedUsers: [],
        flows: [],
      },
      payload ?? {},
    );
  }
}
