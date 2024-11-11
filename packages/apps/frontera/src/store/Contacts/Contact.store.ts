import { set } from 'lodash';
import { RootStore } from '@store/root';
import { Syncable } from '@store/syncable';
import { Transport } from '@store/transport';
import { FlowStore } from '@store/Flows/Flow.store';
import { countryMap } from '@assets/countries/countriesMap';
import { action, override, runInAction, makeObservable } from 'mobx';

import {
  Contact,
  DataSource,
  ContactUpdateInput,
} from '@shared/types/__generated__/graphql.types';

import { ContactService } from './__service__/Contacts.service';

export class ContactStore extends Syncable<Contact> {
  private service: ContactService;

  constructor(
    public root: RootStore,
    public transport: Transport,
    data: Contact,
  ) {
    super(root, transport, data ?? getDefaultValue());
    this.service = ContactService.getInstance(transport);

    makeObservable(this, {
      id: override,
      save: override,
      getId: override,
      setId: override,
      invalidate: action,
      getChannelName: override,
    });
  }

  getChannelName(): string {
    return 'Contacts';
  }

  set id(id: string) {
    this.value.id = id;
    this.value.metadata.id = id;
  }

  get id() {
    return this.value.id ?? this.value.metadata.id;
  }

  get organizationId() {
    return this.value.organizations.content[0]?.metadata?.id;
  }

  get hasFlows() {
    return this.value.flows?.length > 0;
  }

  get flows(): FlowStore[] | undefined {
    if (!this.value.flows?.length) return undefined;

    return this.value.flows.map((flow) => {
      return this.root.flows?.value.get(flow.metadata.id) as FlowStore;
    });
  }

  get flowsIds(): string[] | undefined {
    if (!this.flows?.length) return undefined;

    return this.flows.map((flow) => {
      return flow?.id;
    });
  }

  get name() {
    return (
      this.value.name || `${this.value.firstName} ${this.value.lastName}`.trim()
    );
  }

  get emailId() {
    return this.value.emails?.[0]?.id;
  }

  get connectedUsers() {
    return this.value.connectedUsers.map(
      ({ id }) => this.root.users.value.get(id)?.value,
    );
  }

  get country() {
    if (!this.value.locations?.[0]?.countryCodeA2) return undefined;

    return countryMap.get(this.value.locations[0].countryCodeA2.toLowerCase());
  }

  get organization() {
    return this.root.organizations.value.get(this.organizationId)?.value;
  }

  async getRecentChanges() {}

  setId(id: string) {
    this.value.id = id;
    this.value.metadata.id = id;
  }

  getId() {
    return this.value.id ?? this.value.metadata.id;
  }

  deletePersona(personaId: string) {
    this.value.tags = (this.value?.tags || []).filter(
      (id) => id.id !== personaId,
    );
  }

  get hasActiveOrganization() {
    const org = this.root.organizations.value.get(this.organizationId);

    return org && !org.value.hide;
  }

  getFlowById(flowId: string) {
    return this.value.flows.find((flow) => flow.metadata.id === flowId);
  }

  async invalidate(): Promise<void> {
    try {
      this.isLoading = true;

      const { contact } = await this.service.getContact(this.value.id);

      this.load(contact as Contact);
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async linkOrganization(organizationId: string) {
    try {
      this.isLoading = true;

      const { contact_AddOrganizationById } =
        await this.service.linkOrganization({
          input: {
            contactId: this.value.id,
            organizationId,
          },
        });

      this.load(contact_AddOrganizationById as Contact);
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateContact(input: ContactUpdateInput) {
    try {
      this.isLoading = true;

      const { contact_Update } = await this.service.updateContact({
        input: {
          ...input,
          id: this.value.id,
          patch: true,
        },
      });

      this.load(contact_Update as Contact);
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async updateContactName() {
    const name = this.value.name;

    try {
      await this.service.updateContact({
        input: {
          id: this.value.id,
          name: name,
        },
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
        contactId: this.getId(),
        input: {
          organizationId: this.organizationId,
          description: this.value.jobRoles[0].description,
          jobTitle: this.value.jobRoles[0].jobTitle,
        },
      });

      this.load(jobRole_Create as Contact);
      runInAction(() => {
        set(this.value.jobRoles[0], 'id', jobRole_Create.id);
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
        contactId: this.getId(),
        input: {
          id: this.value.jobRoles[0].id,
          description: this.value.jobRoles[0].description,
          jobTitle: this.value.jobRoles[0].jobTitle,
        },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async updateEmail(
    previousEmail: string,
    index?: number,
    primary: boolean = false,
  ) {
    const email = this.value.emails?.[index ?? 0]?.email ?? '';

    try {
      this.isLoading = true;
      await this.service.updateContactEmail({
        contactId: this.getId(),
        input: {
          email,
          primary: primary,
        },
        previousEmail,
      });
      runInAction(() => {
        this.isLoading = false;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      this.invalidate();
    }
  }

  async updateEmailPrimary(previousEmail: string) {
    const email = this.value.primaryEmail?.email ?? '';

    try {
      await this.service.updateContactEmail({
        contactId: this.getId(),
        input: {
          email,
          primary: true,
        },
        previousEmail,
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      this.invalidate();
    }
  }

  async addPhoneNumber() {
    const phoneNumber = this.value.phoneNumbers?.[0].rawPhoneNumber ?? '';

    try {
      const { phoneNumberMergeToContact } = await this.service.addPhoneNumber({
        contactId: this.getId(),
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
        contactId: this.getId(),
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async addSocial(
    url: string,
    options?: { onSuccess?: (serverId: string) => void },
  ) {
    try {
      const { contact_AddSocial } = await this.service.addSocial({
        contactId: this.getId(),
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
    } finally {
      options?.onSuccess?.(this.value.socials?.[0]?.id);
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

  async findEmail(isLoading?: (isLoading: boolean) => void) {
    this.isLoading = true;
    isLoading?.(this.isLoading);

    try {
      await this.service.findEmail({
        contactId: this.getId(),
        organizationId: this.organizationId,
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      this.isLoading = false;
    }
  }

  async setPrimaryEmail(emailId: string) {
    const email = this.value.emails.find((email) => email.id === emailId);

    try {
      await this.service.setPrimaryEmail({
        contactId: this.getId(),
        email: email?.email || '',
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      this.isLoading = false;
      this.invalidate();
    }
  }

  async addTagToContact(tagId: string, tagName: string) {
    try {
      await this.service.addTagsToContact({
        input: {
          contactId: this.getId(),
          tag: {
            id: tagId,
            name: tagName,
          },
        },
      });
      runInAction(() => {
        this.root.ui.toastSuccess('Tag has been added', 'tags-added-success');
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async removeTagFromContact(tagId: string) {
    try {
      await this.service.removeTagsFromContact({
        input: {
          contactId: this.getId(),
          tag: {
            id: tagId,
          },
        },
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    }
  }

  async removeAllTagsFromContact() {
    const tags =
      this.value?.tags?.map((tag) => this.removeTagFromContact(tag.id)) || [];

    try {
      await Promise.all(tags);

      runInAction(() => {
        this.value.tags = [];
        this.root.ui.toastSuccess(
          'All tags were removed',
          'tags-remove-success',
        );
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
  enrichDetails: {
    enrichedAt: '',
    failedAt: '',
    requestedAt: '',
  },
});
