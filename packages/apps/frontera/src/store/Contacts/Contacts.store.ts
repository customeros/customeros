import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Store } from '@store/store.ts';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { when, runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import {
  Tag,
  Contact,
  Pagination,
  ContactInput,
  Organization,
} from '@graphql/types';

import mock from './mock.json';
import { ContactStore } from './Contact.store';
import { ContactService } from './Contact.service';

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
        this.isBootstrapped && this.totalElements > 0 && !this.root.demoMode,
      async () => {
        await this.bootstrapRest();
      },
    );
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray(compute: (arr: Store<Contact>[]) => Contact[]) {
    const arr = this.toArray();

    return compute(arr);
  }

  updateTags = (ids: string[], tags: Tag[]) => {
    ids.forEach((id) => {
      this.value.get(id)?.update((contact) => {
        contact.tags = [...(contact.tags ?? []), ...tags];

        return contact;
      });
    });
  };

  archive = (ids: string[]) => {
    ids.forEach((id) => {
      this.remove(id);
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
        this.isBootstrapped = true;
        this.totalElements = contacts.totalElements;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
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
    organizationId?: string,
    options?: { onSuccess?: (serverId: string) => void },
  ) {
    const newContact = new ContactStore(this.root, this.transport);
    const tempId = newContact.value.id;
    let serverId: string | undefined;

    this.value.set(tempId, newContact);
    if (organizationId) {
      const organization = this.root.organizations.value.get(organizationId);
      organization?.update(
        (v: Organization) => {
          v.contacts.content.push(newContact.value);

          return v;
        },
        { mutate: false },
      );
    }

    try {
      const { contact_Create } = await this.transport.graphql.request<
        CREATE_CONTACT_MUTATION_RESPONSE,
        CREATE_CONTACT_MUTATION_PAYLOAD
      >(CREATE_CONTACT_MUTATION, {
        input: {},
      });

      runInAction(() => {
        serverId = contact_Create.id;
        newContact.value.id = serverId;

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
        }
      }, 1000);
    }
  }

  async remove(id: string) {
    try {
      runInAction(() => {
        const organizationId = this.value.get(id)?.organizationId;

        if (organizationId) {
          const organization =
            this.root.organizations.value.get(organizationId);

          organization?.update(
            (v: Organization) => {
              v.contacts.content = v.contacts.content.filter(
                (c) => c.id !== id,
              );

              return v;
            },
            { mutate: false },
          );
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
        tags {
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
        }
        locations {
          id
          address
          locality
          postalCode
          country
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
            isReachable
            isValidSyntax
            canConnectSmtp
            acceptsMail
            hasFullInbox
            isCatchAll
            isDeliverable
            validated
            isDisabled
          }
        }
        socials {
          id
          url
        }
        profilePhotoUrl
      }
      totalPages
      totalElements
    }
  }
`;

type CREATE_CONTACT_MUTATION_RESPONSE = {
  contact_Create: {
    id: string;
  };
};
type CREATE_CONTACT_MUTATION_PAYLOAD = {
  input: ContactInput;
};
const CREATE_CONTACT_MUTATION = gql`
  mutation createContact($input: ContactInput!) {
    contact_Create(input: $input) {
      id
    }
  }
`;
