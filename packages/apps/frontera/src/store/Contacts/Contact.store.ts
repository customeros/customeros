import type { RootStore } from '@store/root';

import set from 'lodash/set';
import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';
import { rdiffResult } from 'recursive-diff';
import { Store, makeAutoSyncable } from '@store/store';
import { runInAction, makeAutoObservable } from 'mobx';
import { countryMap } from '@assets/countries/countriesMap';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store.ts';

import { Tag, Contact, DataSource, ContactUpdateInput } from '@graphql/types';

import { ContactService } from './Contact.service';

interface ContractStore {
  get name(): string;
}

export class ContactStore implements Store<Contact>, ContractStore {
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

  constructor(public root: RootStore, public transport: Transport) {
    this.value = getDefaultValue();

    makeAutoSyncable(this, {
      channelName: 'Contact',
      mutator: this.save,
      getId: (d) => d?.metadata.id,
    });
    makeAutoObservable(this);
    this.service = ContactService.getInstance(transport);
  }

  async invalidate() {
    try {
      this.isLoading = true;

      const { contact } = await this.transport.graphql.request<
        CONTACT_QUERY_RESULT,
        { id: string }
      >(CONTACT_QUERY, { id: this.id });

      this.load(contact);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  set id(id: string) {
    this.value.id = id;
  }

  get id() {
    return this.value.id;
  }

  get organizationId() {
    return this.value.organizations.content[0]?.metadata?.id;
  }

  get sequence(): FlowSequenceStore | undefined {
    if (!this.value.sequences?.length) return undefined;

    return this.root.flowSequences?.value.get(
      this.value.sequences[0]?.metadata.id,
    ) as FlowSequenceStore;
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

  setId(id: string) {
    this.value.metadata.id = id;
  }

  getId() {
    return this.value.metadata.id;
  }

  deletePersona(personaId: string) {
    this.value.tags = (this.value?.tags || []).filter(
      (id) => id.id !== personaId,
    );
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const type = diff?.op;
    const path = diff?.path;
    const value = diff?.val;
    const oldValue = (diff as rdiffResult & { oldVal: unknown })?.oldVal;

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
          this.addSocial(value.url);
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
      .with(['tags', ...P.array()], () => {
        if (type === 'add') {
          this.addTagToContact(value.id, value.name);
        }

        if (type === 'delete') {
          if (typeof oldValue === 'object') {
            this.removeTagFromContact(oldValue.id);
          }
        }

        // if tag with index different that last one is deleted it comes as an update, bulk creation updates also come as updates
        if (type === 'update') {
          if (!oldValue) {
            (value as Array<Tag>)?.forEach((tag: Tag) => {
              this.addTagToContact(tag.id, tag.name);
            });
          }

          if (oldValue) {
            this.removeTagFromContact(oldValue);
          }
        }
      })
      .otherwise(() => {
        const payload = makePayload<ContactUpdateInput>(operation);

        this.updateContact(payload);
      });
  }

  async linkOrganization(organizationId: string) {
    try {
      await this.service.linkOrganization({
        input: {
          contactId: this.getId(),
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
        input: { ...input, id: this.getId(), patch: true },
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

  async addEmail() {
    const email = this.value.emails?.[0]?.email ?? '';

    try {
      const { emailMergeToContact } = await this.service.addContactEmail({
        contactId: this.getId(),
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
        contactId: this.getId(),
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

  async addSocial(url: string) {
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
        contactId: this.getId(),
        organizationId: this.organizationId,
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
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
type CONTACT_QUERY_RESULT = {
  contact: Contact;
};
const CONTACT_QUERY = gql`
  query contact($id: ID!) {
    contact(id: $id) {
      firstName
      lastName
      name

      prefix
      description
      timezone
      metadata {
        id
      }
      tags {
        metadata {
          id
          source
          sourceOfTruth
          appSource
          created
          lastUpdated
        }
        id
        name
        source
        updatedAt
        createdAt
        appSource
      }
      organizations(pagination: { limit: 2, page: 0 }) {
        content {
          metadata {
            id
          }
          id
          name
        }
        totalElements
        totalAvailable
      }
      tags {
        id
        name
      }
      jobRoles {
        id
        primary
        jobTitle
        description
        company
        startedAt
        endedAt
      }

      locations {
        id
        address
        locality
        postalCode
        country
        region
        countryCodeA2
        countryCodeA3
      }

      phoneNumbers {
        id
        e164
        rawPhoneNumber
        label
        primary
      }
      emails {
        id
        email
        emailValidationDetails {
          verified
          verifyingCheckAll
          isValidSyntax
          isRisky
          isFirewalled
          provider
          firewall
          isCatchAll
          canConnectSmtp
          deliverable
          isMailboxFull
          isRoleAccount
          isFreeAccount
          smtpSuccess
        }
      }
      socials {
        id
        url
        alias
        followersCount
      }
      connectedUsers {
        id
      }
      profilePhotoUrl
    }
  }
`;

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
  phoneNumbers: [],
  profilePhotoUrl: '',
  organizations: {
    content: [],
    totalPages: 0,
    totalElements: 0,
    totalAvailable: 0,
  },
  sequences: [],
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
});
