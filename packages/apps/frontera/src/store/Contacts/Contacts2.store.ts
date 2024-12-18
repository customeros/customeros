import { Store } from '@store/_store';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { action, computed, observable, runInAction } from 'mobx';

import {
  Tag,
  ContactInput,
  SortingDirection,
} from '@shared/types/__generated__/graphql.types';

import { Contact, ContactDatum } from './Contact.dto';
import { ContactsView } from './__views__/Contacts.view';
import { ContactService } from './__service__/Contacts.service';
import { FlowContactsView } from './__views__/FlowContacts.view';

export class ContactsStore extends Store<ContactDatum, Contact> {
  private chunkSize = 50;
  private service: ContactService;
  @observable accessor searchedIds: string[] = [];
  @observable accessor chunk = 0;
  @observable accessor availableCounts: Map<string, number> = new Map();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Contacts',
      getId: (data) => data?.id,
      factory: Contact,
    });
    this.service = ContactService.getInstance();

    new ContactsaView(this);
    new FlowContactsView(this);
  }

  @computed
  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  delete = (ids: string[]) => {
    ids.forEach((id) => {
      this.delete([id]);
    });
  };

  archive = (ids: string[]) => {
    ids.forEach((id) => {
      this.softDelete(id);
    });
  };

  @computed
  get canLoadNext() {
    return this.searchedIds.length > this.chunkSize * (this.chunk + 1);
  }

  @action
  async getAllData() {
    runInAction(() => {
      this.isBootstrapping = true;
    });

    try {
      const { ui_contacts_search } = await this.service.searchContacts({
        limit: this.chunkSize,
        sort: {
          by: 'CONTACTS_CREATED_AT',
          caseSensitive: false,
          direction: SortingDirection.Desc,
        },
      });

      await this.retrieve(ui_contacts_search.ids);

      // const totalElements = ui_organizations_search.totalElements;

      runInAction(() => {
        this.size = this.value.size;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
        this.isBootstrapped = true;
        this.isBootstrapping = false;
      });
    }
  }

  @action
  async search(viewDefPrest: string) {
    const viewDef = this.root.tableViewDefs.getById(viewDefPrest);

    if (!viewDef) {
      console.error(`viewDef with preset=${viewDefPrest} not found`);

      return;
    }

    try {
      runInAction(() => {
        if (this.chunk > 0) {
          // reset chunk if new search is performed
          this.chunk = 0;
        }
        this.isLoading = true;
      });

      const payload = viewDef.toSearchPayload();

      const { ui_contacts_search: searchResult } =
        await this.service.searchContacts({ ...payload });

      if (this.chunk === 0) {
        const ids = (searchResult?.ids ?? []).slice(
          this.chunkSize * this.chunk,
          this.chunkSize * this.chunk + this.chunkSize,
        );

        // retrieve first chunk of data after new search is performed
        await this.retrieve(ids);
      }

      runInAction(() => {
        this.isLoading = false;
        this.availableCounts.set(viewDefPrest, searchResult?.totalElements);
        this.searchedIds = searchResult?.ids ?? [];
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    }
  }

  @action
  async retrieve(ids: string[]) {
    try {
      const { ui_contacts } = await this.service.getContactsByIds({
        ids,
      });

      runInAction(() => {
        ui_contacts.forEach((raw) => {
          if (this.value.has(raw.id)) {
            Object.assign(raw, this.value.get(raw.id)?.value);
          } else {
            const record = new Contact(this, raw);

            this.value.set(record.id, record);
          }
        });

        this.size = this.value.size;
        this.version++;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    }
  }

  @action
  public async loadNext() {
    if (!this.canLoadNext) return;

    runInAction(() => {
      this.chunk++;
    });

    const ids = this.searchedIds.slice(
      this.chunkSize * this.chunk,
      this.chunkSize * this.chunk + this.chunkSize,
    );

    try {
      runInAction(() => {
        this.isLoading = true;
      });

      await this.retrieve(ids);
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

  @action
  async create(
    organizationId: string,
    options?: { onSuccess?: (serverId: string) => void },
    input?: ContactInput,
  ) {
    const newContact = new Contact(this, Contact.default());
    const tempId = newContact.id;
    let serverId: string | undefined;

    this.value.set(tempId, newContact);

    if (organizationId) {
      const organization = this.root.organizations.value.get(organizationId);

      organization?.value.contacts.content.unshift(newContact.value);
      organization?.commit({ syncOnly: true });
    }

    try {
      const { contact_CreateForOrganization } =
        await this.service.createContactForOrganization({
          organizationId,
          input: input ?? {},
        });

      runInAction(() => {
        serverId = contact_CreateForOrganization.id;
        newContact.id = serverId;
        newContact.commit({ syncOnly: true });

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

  @action
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

    const newContact = new Contact(this, Contact.default());

    const tempId = newContact.id;
    const socialId = crypto.randomUUID();

    (newContact.value = {
      __typename: 'Contact',
      metadata: {
        id: socialId,
      },
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      firstName: null,
      lastName: null,
      name: null,
      prefix: null,
      flows: [],
      enrichDetails: {},
      jobRoles: [],
      locations: [],
      phoneNumbers: [],
      emails: [],
      socials: [],
      connectedUsers: [],
      tags: [],
      timezone: null,
      organizations: {
        content: [],
        totalAvailable: 0,
        totalElements: 0,
      },
      profilePhotoUrl: null,
      description: null,
      latestOrganizationWithJobRole: {
        jobRole: {
          id: '',
          primary: false,
          jobTitle: '',
          description: '',
          company: '',
          startedAt: null,
          endedAt: null,
        },
        organization: {
          metadata: {
            id: '',
          },
          name: '',
        },
      },
    }),
      this.value.set(tempId, newContact);

    let serverId: string | undefined;

    const organization = this.root.organizations.value.get(organizationId);

    if (organization) {
      organization?.value?.contacts.content.unshift(newContact.value);
      organization.commit({ syncOnly: true });
    }

    try {
      const { contact_CreateForOrganization } =
        await this.service.createContactForOrganization({
          organizationId,
          input: {
            socialUrl,
          },
        });

      runInAction(() => {
        serverId = contact_CreateForOrganization.id;
        newContact.id = serverId;
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
      }, 2000);
    }
  }

  @action
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

    const newContact = new Contact(this, Contact.default());
    const tempId = newContact.id;
    const socialId = crypto.randomUUID();
    let serverId: string | undefined = undefined;

    (newContact.value = {
      __typename: 'Contact',
      metadata: {
        id: socialId,
      },
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      firstName: null,
      lastName: null,
      name: null,
      prefix: null,
      flows: [],
      enrichDetails: {},
      jobRoles: [],
      locations: [],
      phoneNumbers: [],
      emails: [],
      socials: [],
      connectedUsers: [],
      tags: [],
      timezone: null,
      organizations: {
        content: [],
        totalAvailable: 0,
        totalElements: 0,
      },
      profilePhotoUrl: null,
      description: null,
      latestOrganizationWithJobRole: {
        jobRole: {
          id: '',
          primary: false,
          jobTitle: '',
          description: '',
          company: '',
          startedAt: null,
          endedAt: null,
        },
        organization: {
          metadata: {
            id: '',
          },
          name: '',
        },
      },
    }),
      this.value.set(tempId, newContact);

    try {
      const { contact_Create } = await this.service.createContact({
        contactInput: {
          socialUrl,
        },
      });

      runInAction(() => {
        serverId = contact_Create;
        newContact.id = serverId;
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
      }, 2000);
    }
  }

  @action
  async createBulkByEmail({
    emails,
    options,
  }: {
    flowId?: string;
    emails: string[];
    options?: {
      onSuccess?: () => void;
      onError?: (err: string) => void;
    };
  }) {
    this.isLoading = true;

    try {
      const { contact_CreateBulkByEmail } =
        await this.service.createContactBulkByEmail({
          emails,
        });

      runInAction(() => {
        this.sync({ action: 'APPEND', ids: contact_CreateBulkByEmail });
        options?.onSuccess?.();
        this.isLoading = false;
      });

      this.root.ui.toastSuccess(`Contacts created`, 'create-contact-success');
    } catch (e) {
      this.root.ui.toastError(
        `We couldn't create this contact. Please try again.`,
        'create-contact-error',
      );
      runInAction(() => {
        this.error = (e as Error)?.message;
        options?.onError?.(this.error);
      });
    } finally {
      setTimeout(() => {
        this.isBootstrapped = false;
        this.bootstrap();
      }, 300);
    }
  }

  @action
  async createBulkByLinkedIn({
    linkedInUrls,
    options,
  }: {
    flowId?: string;
    linkedInUrls: string[];
    options?: {
      onSuccess?: () => void;
      onError?: (err: string) => void;
    };
  }) {
    this.isLoading = true;

    try {
      const { contact_CreateBulkByLinkedIn } =
        await this.service.createContactBulkByLinkedIn({
          linkedInUrls,
        });

      runInAction(() => {
        this.sync({ action: 'APPEND', ids: contact_CreateBulkByLinkedIn });
        options?.onSuccess?.();
        this.isLoading = false;
      });

      this.root.ui.toastSuccess(`Contacts created`, 'create-contact-success');
    } catch (e) {
      this.root.ui.toastError(
        `We couldn't create this contact. Please try again.`,
        'create-contact-error',
      );
      runInAction(() => {
        this.error = (e as Error)?.message;
        options?.onError?.(this.error);
      });
    } finally {
      setTimeout(() => {
        this.isBootstrapped = false;
        this.bootstrap();
      }, 300);
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
            (c) => c?.metadata.id === id,
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
            (c) => c?.metadata.id === id,
          );

          if (foundIdx && foundIdx > -1) {
            organization?.value?.contacts.content.splice(foundIdx, 1);
            organization?.commit({ syncOnly: true });
          }
        }
        this.value.delete(id);
      });

      await this.service.archiveContact({ contactId: id });
      this.totalElements--;
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

  updateTags = (ids: string[], tags: Tag[]) => {
    const tagIdsToUpdate = new Set(tags?.map((tag) => tag.metadata.id));

    const shouldRemoveTags = ids.every((id) => {
      const contact = this.value.get(id);

      if (!contact) return false;

      const contactIdsTags = new Set(
        (contact.value.tags ?? []).map((tag) => tag.metadata.id),
      );

      return Array.from(tagIdsToUpdate).every((tagId) =>
        contactIdsTags.has(tagId),
      );
    });

    ids.forEach((id) => {
      const contact = this.value.get(id);

      if (!contact) return;

      if (shouldRemoveTags) {
        contact.value.tags = contact.value.tags?.filter(
          (t) => !tagIdsToUpdate.has(t.metadata.id),
        );
      } else {
        const existingIds = new Set(
          contact.value.tags?.map((t) => t.metadata.id) ?? [],
        );
        const newTags = tags.filter((t) => !existingIds.has(t.metadata.id));

        if (!Array.isArray(contact.value.tags)) {
          contact.value.tags = [];
        }

        contact.value.tags = [...(contact.value.tags ?? []), ...newTags];

        contact.commit();
      }
    });
  };
}
