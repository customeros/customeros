import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { when, runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Contact, DataSource, Pagination, ContactInput } from '@graphql/types';

import mock from './mock.json';
import { ContactStore } from './Contact.store';
import { ContactService } from './__service__/Contact.service.ts';

export class ContactsStore implements GroupStore<Contact> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, ContactStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<Contact>();
  totalElements = 0;
  isFullyLoaded = false;
  private service: ContactService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = ContactService.getInstance(transport);

    makeAutoSyncableGroup(this, {
      channelName: 'Contacts',
      getItemId: (item) => item?.id,
      ItemStore: ContactStore,
    });
    makeAutoObservable(this);

    when(
      () =>
        this.isBootstrapped &&
        this.totalElements > 0 &&
        this.totalElements !== this.value.size &&
        !this.root.demoMode,
      async () => {
        await this.bootstrapRest();
      },
    );

    when(
      () => this.isBootstrapped && this.totalElements === this.value.size,
      () => {
        this.isFullyLoaded = true;
        this.isLoading = false;
      },
    );
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray(compute: (arr: ContactStore[]) => ContactStore[]) {
    const arr = this.toArray();

    return compute(arr);
  }

  delete = (ids: string[]) => {
    ids.forEach((id) => {
      this.remove(id);
    });
  };

  archive = (ids: string[]) => {
    ids.forEach((id) => {
      this.softDelete(id);
    });
  };

  async bootstrap() {
    if (this.root.demoMode) {
      this.load(mock.data.contacts.content as unknown as Contact[]);
      this.isBootstrapped = true;
      this.totalElements = mock.data.contacts.totalElements;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;

      const { contacts } = await this.transport.graphql.request<
        CONTACTS_QUERY_RESPONSE,
        CONTACTS_QUERY_PAYLOAD
      >(CONTACTS_QUERY, {
        pagination: { limit: 1000, page: 0 },
      });

      this.load(contacts.content);
      runInAction(() => {
        this.totalElements = contacts.totalElements;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      this.isBootstrapped = true;
    }
  }

  async bootstrapRest() {
    let page = 1;

    while (this.totalElements > this.value.size) {
      try {
        const { contacts } = await this.transport.graphql.request<
          CONTACTS_QUERY_RESPONSE,
          CONTACTS_QUERY_PAYLOAD
        >(CONTACTS_QUERY, {
          pagination: { limit: 1000, page },
        });

        runInAction(() => {
          page++;
          this.load(contacts.content);
        });
      } catch (e) {
        runInAction(() => {
          this.error = (e as Error)?.message;
        });
        break;
      }
    }
  }

  async create(
    organizationId: string,
    options?: { onSuccess?: (serverId: string) => void },
    input?: ContactInput,
  ) {
    const newContact = new ContactStore(this.root, this.transport);
    const tempId = newContact.value.metadata?.id;
    let serverId: string | undefined;

    this.value.set(tempId, newContact);

    if (organizationId) {
      const organization = this.root.organizations.value.get(organizationId);

      organization?.value.contacts.content.unshift(newContact.value);
      organization?.commit({ syncOnly: true });
    }

    try {
      const { contact_CreateForOrganization } =
        await this.transport.graphql.request<
          CREATE_CONTACT_MUTATION_RESPONSE,
          CREATE_CONTACT_MUTATION_PAYLOAD
        >(CREATE_CONTACT_MUTATION, {
          organizationId,
          input: input || {},
        });

      runInAction(() => {
        serverId = contact_CreateForOrganization.id;
        newContact.setId(serverId);

        this.value.set(serverId, newContact);
        this.value.delete(tempId);

        this.sync({ action: 'APPEND', ids: [serverId] });
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      serverId && options?.onSuccess?.(serverId);
      setTimeout(() => {
        if (serverId) {
          this.value.get(serverId)?.invalidate();
          this.root.organizations.value.get(organizationId)?.invalidate();
        }
      }, 1000);
    }
  }

  async createWithSocial({
    socialUrl,
    organizationId,
    options,
  }: {
    socialUrl: string;
    organizationId: string;
    options?: {
      onSuccess?: (serverId: string) => void;
    };
  }) {
    this.isLoading = true;

    const newContact = new ContactStore(this.root, this.transport);
    const tempId = newContact.value.id;
    const socialId = crypto.randomUUID();

    newContact.value.socials = [
      {
        metadata: {
          id: socialId,
          source: DataSource.Openline,
          sourceOfTruth: DataSource.Openline,
          appSource: 'organization',
          created: new Date().toISOString(),
          lastUpdated: new Date().toISOString(),
        },
        id: socialId,
        externalId: '',
        url: socialUrl,
        appSource: 'OPENLINE',
        createdAt: new Date().toISOString(),
        sourceOfTruth: DataSource.Openline,
        source: DataSource.Openline,
        alias: socialUrl,
        followersCount: 0,
        updatedAt: new Date().toISOString(),
      },
    ];

    let serverId: string | undefined;

    this.value.set(tempId, newContact);

    const organization = this.root.organizations.value.get(organizationId);

    if (organization) {
      organization?.value?.contacts.content.unshift(newContact.value);
      organization.commit({ syncOnly: true });
    }

    try {
      const { contact_CreateForOrganization } =
        await this.transport.graphql.request<
          CREATE_CONTACT_MUTATION_RESPONSE,
          CREATE_CONTACT_MUTATION_PAYLOAD
        >(CREATE_CONTACT_MUTATION, {
          organizationId,
          input: {
            socialUrl,
          },
        });

      runInAction(() => {
        serverId = contact_CreateForOrganization.id;
        newContact.value.id = serverId;

        this.value.set(serverId, newContact);
        this.value.delete(tempId);

        this.sync({ action: 'APPEND', ids: [serverId] });
        this.isLoading = false;
      });
      this.root.ui.toastSuccess(
        `Contact created for ${organization?.value?.name}`,
        'create-contract-error',
      );
    } catch (e) {
      this.root.ui.toastError(
        `We couldn't create this contact. Please try again.`,
        'create-contract-error',
      );
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      serverId && options?.onSuccess?.(serverId);
      setTimeout(() => {
        if (serverId) {
          this.value.get(serverId)?.invalidate();
        }
      }, 600);
      setTimeout(() => {
        if (serverId) {
          this.value.get(serverId)?.invalidate();
        }
      }, 2000);
    }
  }

  async createWithoutOrg({
    socialUrl,
    options,
  }: {
    socialUrl: string;
    options?: {
      onSuccess?: (serverId: string) => void;
    };
  }) {
    this.isLoading = true;

    const newContact = new ContactStore(this.root, this.transport);
    const tempId = newContact.value.id;
    const socialId = crypto.randomUUID();
    let serverId: string | undefined = undefined;

    newContact.value.socials = [
      {
        metadata: {
          id: socialId,
          source: DataSource.Openline,
          sourceOfTruth: DataSource.Openline,
          appSource: 'organization',
          created: new Date().toISOString(),
          lastUpdated: new Date().toISOString(),
        },
        id: socialId,
        externalId: '',
        url: socialUrl,
        appSource: 'OPENLINE',
        createdAt: new Date().toISOString(),
        sourceOfTruth: DataSource.Openline,
        source: DataSource.Openline,
        alias: socialUrl,
        followersCount: 0,
        updatedAt: new Date().toISOString(),
      },
    ];

    this.value.set(tempId, newContact);

    try {
      const { contact_Create } = await this.transport.graphql.request<
        CREATE_CONTACT_MUTATION_WITHOUT_ORG_RESPONSE,
        CREATE_CONTACT_MUTATION_WITHOUT_ORG_PAYLOAD
      >(CREATE_CONTACT_MUTATION_WITHOUT_ORG, {
        contactInput: {
          socialUrl,
        },
      });

      runInAction(() => {
        serverId = contact_Create;
        newContact.value.id = serverId;

        this.value.set(serverId, newContact);
        this.value.delete(tempId);

        this.sync({ action: 'APPEND', ids: [serverId] });
        this.isLoading = false;
      });

      this.root.ui.toastSuccess(`Contact created`, 'create-contact-success');
    } catch (e) {
      this.root.ui.toastError(
        `We couldn't create this contact. Please try again.`,
        'create-contact-error',
      );
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      serverId && options?.onSuccess?.(serverId);
      setTimeout(() => {
        if (serverId) {
          this.value.get(serverId)?.invalidate();
        }
      }, 600);
      setTimeout(() => {
        if (serverId) {
          this.value.get(serverId)?.invalidate();
        }
      }, 2000);
    }
  }

  async remove(id: string) {
    try {
      runInAction(() => {
        const organizationId = this.value.get(id)?.organizationId;

        if (organizationId) {
          const organization =
            this.root.organizations.value.get(organizationId);

          const foundIdx = organization?.value?.contacts.content.findIndex(
            (c) => c?.id === id,
          );

          if (foundIdx && foundIdx > -1) {
            organization?.value?.contacts.content.splice(foundIdx, 1);
            organization?.commit({ syncOnly: true });
          }
        }
        this.value.delete(id);
      });

      await this.service.deleteContact({ contactId: id });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.sync({ action: 'DELETE', ids: [id] });
      });
    }
  }

  async softDelete(id: string) {
    try {
      runInAction(() => {
        const organizationId = this.value.get(id)?.organizationId;

        if (organizationId) {
          const organization =
            this.root.organizations.value.get(organizationId);

          const foundIdx = organization?.value?.contacts.content.findIndex(
            (c) => c?.id === id,
          );

          if (foundIdx && foundIdx > -1) {
            organization?.value?.contacts.content.splice(foundIdx, 1);
            organization?.commit({ syncOnly: true });
          }
        }
        this.value.delete(id);
      });

      await this.service.archiveContact({ contactId: id });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.sync({ action: 'DELETE', ids: [id] });
      });
    }
  }

  async deleteFlowContacts(contactIds: string[]) {
    this.isLoading = true;

    try {
      await Promise.all(
        contactIds.map((contactId) => {
          this.value.get(contactId)?.deleteFlowContact();
        }),
      );

      runInAction(() => {
        const contactStores = contactIds.map((e) => this.value.get(e));

        contactStores.forEach((contactStore) => {
          contactStore?.update(
            (c) => {
              c.flows = [];

              return c;
            },
            { mutate: false },
          );
        });

        this.root.ui.toastSuccess(
          `${contactIds.length} contacts removed from their flows`,
          'unlink-contact-from-flow-success',
        );
        this.root.contacts.sync({
          action: 'INVALIDATE',
          ids: contactIds,
        });

        this.root.flows.invalidate();
      });
    } catch (e) {
      runInAction(() => {
        this.root.ui.toastError(
          `We couldn't remove those contacts from their flows`,
          'unlink-contact-from-sequence-error',
        );
      });
    } finally {
      this.isLoading = false;
    }
  }
}

type CONTACTS_QUERY_RESPONSE = {
  contacts: {
    totalPages: number;
    content: Contact[];
    totalElements: number;
  };
};
type CONTACTS_QUERY_PAYLOAD = {
  pagination: Pagination;
};
const CONTACTS_QUERY = gql`
  query contacts($pagination: Pagination!) {
    contacts(pagination: $pagination) {
      content {
        id
        name
        firstName
        lastName
        prefix
        description
        timezone
        metadata {
          id
        }
        flows {
          metadata {
            id
          }
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
        updatedAt
        createdAt
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
        primaryEmail {
          id
          primary
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
        latestOrganizationWithJobRole {
          organization {
            name
            metadata {
              id
            }
          }
          jobRole {
            id
            primary
            jobTitle
            description
            company
            startedAt
            endedAt
          }
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
          countryCodeA2
          countryCodeA3
          region
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
          primary
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
          work
        }

        socials {
          id
          url
          alias
          followersCount
        }
        enrichDetails {
          enrichedAt
          failedAt
          requestedAt
        }
        connectedUsers {
          id
        }
        profilePhotoUrl
      }
      totalPages
      totalElements
    }
  }
`;

type CREATE_CONTACT_MUTATION_RESPONSE = {
  contact_CreateForOrganization: {
    id: string;
  };
};
type CREATE_CONTACT_MUTATION_PAYLOAD = {
  input: ContactInput;
  organizationId: string;
};
const CREATE_CONTACT_MUTATION = gql`
  mutation createContact($input: ContactInput!, $organizationId: ID!) {
    contact_CreateForOrganization(
      input: $input
      organizationId: $organizationId
    ) {
      id
    }
  }
`;

type CREATE_CONTACT_MUTATION_WITHOUT_ORG_RESPONSE = {
  contact_Create: string;
};

type CREATE_CONTACT_MUTATION_WITHOUT_ORG_PAYLOAD = {
  contactInput: ContactInput;
};
const CREATE_CONTACT_MUTATION_WITHOUT_ORG = gql`
  mutation CreateContact($contactInput: ContactInput!) {
    contact_Create(input: $contactInput)
  }
`;
