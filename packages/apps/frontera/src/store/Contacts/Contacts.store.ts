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
import { TargetsContactsView } from './__views__/TargetsContacts.view';

export class ContactsStore extends Store<ContactDatum, Contact> {
  private chunkSize = 50;
  private service = ContactService.getInstance();
  @observable accessor searchResults: Map<string, string[]> = new Map();
  @observable accessor cursors: Map<string, number> = new Map();
  @observable accessor availableCounts: Map<string, number> = new Map();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Contacts',
      getId: (data) => data?.id,
      factory: Contact,
    });

    this.hydrate();

    new ContactsView(this);
    new FlowContactsView(this);
    new TargetsContactsView(this);
  }

  canLoadNext(preset: string) {
    const ids = this.searchResults.get(preset);
    const cursor = this.cursors.get(preset) ?? 0;

    if (!ids) return false;

    return ids.length > this.chunkSize * (cursor + 1);
  }

  @computed
  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  archive = (ids: string[]) => {
    ids.forEach((id) => {
      this.softDelete(id);
    });
  };

  public getById(id: string) {
    if (!this.value || typeof id !== 'string') return null;

    if (!this.value.has(id)) {
      setTimeout(() => {
        this.retrieve([id]);
      }, 0);
    }

    return this.value.get(id) as Contact;
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
    const cursor = (
      this.cursors.has(viewDefPrest)
        ? this.cursors.get(viewDefPrest)
        : this.cursors.set(viewDefPrest, 0).get(viewDefPrest)
    ) as number;

    if (!viewDef) {
      console.error(`viewDef with preset=${viewDefPrest} not found`);

      return;
    }

    try {
      runInAction(() => {
        if (cursor > 0) {
          // reset chunk if new search is performed
          this.cursors.set(viewDefPrest, 0);
        }
        this.isLoading = true;
      });

      const payload = viewDef.toSearchPayload();

      const { ui_contacts_search: searchResult } =
        await this.service.searchContacts({ ...payload });

      if (cursor === 0) {
        const ids = (searchResult?.ids ?? []).slice(
          this.chunkSize * cursor,
          this.chunkSize * cursor + this.chunkSize,
        );

        // retrieve first chunk of data after new search is performed
        await this.retrieve(ids);
      }
      runInAction(() => {
        this.isLoading = false;
        this.availableCounts.set(viewDefPrest, searchResult?.totalElements);
        this.totalElements = searchResult?.totalAvailable ?? 0;
        this.searchResults.set(viewDefPrest, searchResult?.ids ?? []);
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    }
  }

  @action
  async retrieve(ids: string[]) {
    ids.forEach((id) => {
      if (this.value.has(id)) {
        return;
      }

      this.value.set(id, new Contact(this, Contact.default({ id })));
    });

    try {
      const { ui_contacts } = await this.service.getContactsByIds({
        ids,
      });

      runInAction(() => {
        ui_contacts.forEach((raw) => {
          if (this.value.has(raw.id)) {
            const record = this.value.get(raw.id);

            if (!record) return;
            record?.draft();
            Object.assign(record?.value, raw);
            record?.commit({ syncOnly: true });
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
        ids.forEach((id) => {
          this.value.delete(id);
        });
        this.error = (err as Error)?.message;
      });
    }
  }

  @action
  public async loadNext(preset: string) {
    let cursor = this.cursors.get(preset) ?? 0;

    runInAction(() => {
      cursor++;
      this.cursors.set(preset, cursor);
    });

    const ids = this.searchResults.get(preset);

    const chunkedIds = (ids ?? []).slice(
      this.chunkSize * cursor,
      this.chunkSize * cursor + this.chunkSize,
    );

    try {
      runInAction(() => {
        this.isLoading = true;
      });

      await this.retrieve(chunkedIds);
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
  public async invalidate(id: string) {
    try {
      const { ui_contacts } = await this.service.getContactsByIds({
        ids: [id],
      });

      if (!ui_contacts) return;

      const data = ui_contacts[0];

      if (!data) return;

      runInAction(() => {
        const record = this.value.get(id);

        if (record) {
          Object.assign(record.value, data);
        } else {
          const record = new Contact(this, data);

          this.value.set(record.id, record);
        }
      });
    } catch (e) {
      console.error('Failed invalidating Contact with ID: ' + id);
    }
  }

  @action
  async create(
    organizationId: string,
    options?: { onSuccess?: (serverId: string) => void },
    input?: ContactInput,
  ) {
    let serverId: string | undefined;

    try {
      const { contact_CreateForOrganization } =
        await this.service.createContactForOrganization({
          organizationId,
          input: input ?? {},
        });

      runInAction(() => {
        serverId = contact_CreateForOrganization.id;

        const newContact = new Contact(
          this.root.contacts,
          Contact.default({
            name: input?.name || '',
            id: serverId,
          }),
        );

        this.value.set(serverId, newContact);
        this.sync({ action: 'APPEND', ids: [serverId] });
        this.totalElements++;
        this.version++;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      serverId && options?.onSuccess?.(serverId);
      await this.root.contacts.invalidate(serverId!);
      await this.root.organizations.invalidate(organizationId);
    }
  }

  @action
  async createWithEmail(
    organizationId: string,
    options?: { onSuccess?: (serverId: string) => void },
    input?: ContactInput,
  ) {
    let serverId: string | undefined;

    try {
      const { contact_CreateForOrganization } =
        await this.service.createContactForOrganization({
          organizationId,
          input: input ?? {},
        });

      runInAction(() => {
        serverId = contact_CreateForOrganization.id;

        const newContact = new Contact(
          this.root.contacts,
          Contact.default({
            emails: [
              {
                email: input?.email?.email || '',
                primary: true,
                id: '',
                emailValidationDetails: {
                  __typename: undefined,
                  verified: false,
                  verifyingCheckAll: false,
                  isValidSyntax: undefined,
                  isRisky: undefined,
                  isFirewalled: undefined,
                  provider: undefined,
                  firewall: undefined,
                  isCatchAll: undefined,
                  canConnectSmtp: undefined,
                  deliverable: undefined,
                  isMailboxFull: undefined,
                  isRoleAccount: undefined,
                  isFreeAccount: undefined,
                  smtpSuccess: undefined,
                },
              },
            ],
            id: serverId,
          }),
        );

        this.value.set(serverId, newContact);
        this.sync({ action: 'APPEND', ids: [serverId] });
        this.totalElements++;
        this.version++;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      serverId && options?.onSuccess?.(serverId);
      await this.root.contacts.invalidate(serverId!);
      await this.root.organizations.invalidate(organizationId);
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

    this.value.set(tempId, newContact);

    let serverId: string | undefined;

    const organization = this.root.organizations.value.get(organizationId);

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

        const newContact = new Contact(
          this,
          Contact.default({ id: serverId, linkedInUrl: socialUrl }),
        );

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
      await this.root.contacts.invalidate(serverId!);
      await this.root.organizations.invalidate(organizationId);
      this.root.contacts.value.get(serverId!)?.commit({ syncOnly: true });
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
      id: socialId,
      firstName: '',
      lastName: '',
      name: '',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      flows: [],
      locations: [],
      emails: [],
      connectedUsers: [],
      tags: [],
      enrichedEmailEnrichedAt: null,
      enrichedEmailFound: null,
      enrichedFailedAt: null,
      enrichedAt: null,
      description: '',
      phones: [],
      prefix: '',
      timezone: '',
      profilePhotoUrl: '',
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
    const organizationId = this.value.get(id)?.organizationId;

    try {
      runInAction(() => {
        if (organizationId) {
          const organization =
            this.root.organizations.value.get(organizationId);

          const foundIdx = organization?.value?.contacts.findIndex(
            (c) => c === id,
          );

          if (foundIdx && foundIdx > -1) {
            organization?.value?.contacts.splice(foundIdx, 1);
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
    const organizationId = this.value.get(id)?.organizationId;

    try {
      runInAction(() => {
        if (organizationId) {
          const organization =
            this.root.organizations.value.get(organizationId);

          const foundIdx = organization?.value?.contacts.findIndex(
            (c) => c === id,
          );

          if (foundIdx && foundIdx > -1) {
            organization?.draft();
            organization?.value?.contacts.splice(foundIdx, 1);
            organization?.commit({ syncOnly: true });
          }
        }
        this.value.delete(id);
        // this.version++;
        this.totalElements--;
      });

      await this.service.archiveContact({ contactId: id });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.sync({ action: 'DELETE', ids: [id] });
        this.root.organizations.invalidate(organizationId || '');
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
