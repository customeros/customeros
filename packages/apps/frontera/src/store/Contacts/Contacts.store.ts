import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { SyncableGroup } from '@store/syncable-group';
import {
  when,
  action,
  override,
  computed,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import {
  Tag,
  Contact,
  ContactInput,
} from '@shared/types/__generated__/graphql.types';

import mock from './mock.json';
import { ContactStore } from './Contact.store';
import { ContactService } from './__service__/Contacts.service';

export class ContactsStore extends SyncableGroup<Contact, ContactStore> {
  totalElements = 0;
  private service: ContactService;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, ContactStore);
    this.service = ContactService.getInstance(transport);

    makeObservable(this, {
      totalElements: observable,
      create: action.bound,
      channelName: override,
      isFullyLoaded: computed,
      archive: action.bound,
      delete: action.bound,
    });

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
        this.isFullyLoaded && (this.isLoading = false);
        this.isLoading = false;
      },
    );
  }

  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }

  get channelName() {
    return 'Contacts';
  }

  get persisterKey() {
    return 'Contacts';
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
      this.delete([id]);
    });
  };

  archive = (ids: string[]) => {
    ids.forEach((id) => {
      this.softDelete(id);
    });
  };

  async bootstrap() {
    if (this.root.demoMode) {
      this.load(mock.data.contacts.content as unknown as Contact[], {
        getId: (data) => data.metadata.id,
      });
      this.isBootstrapped = true;
      this.totalElements = mock.data.contacts.totalElements;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;

      const { contacts } = await this.service.getContacts({
        pagination: { limit: 1000, page: 0 },
      });

      this.load(contacts.content as Contact[], {
        getId: (data) => data.metadata.id,
      });
      runInAction(() => {
        this.totalElements = contacts.totalElements;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      this.isLoading = false;
      this.isBootstrapped = true;
    }
  }

  async bootstrapRest() {
    let page = 1;

    while (this.totalElements > this.value.size) {
      try {
        const { contacts } = await this.service.getContacts({
          pagination: { limit: 1000, page },
        });

        runInAction(() => {
          page++;
          this.load(contacts.content as Contact[], {
            getId: (data) => data.metadata.id,
          });
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
    const newContact = new ContactStore(
      this.root,
      this.transport,
      ContactStore.getDefaultValue() as Contact,
    );
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
        await this.service.createContactForOrganization({
          organizationId,
          input: input ?? {},
        });

      runInAction(() => {
        serverId = contact_CreateForOrganization.id;
        newContact.setId(serverId);
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

  async createBulkByEmail({
    emails,
    options,
    flowId,
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
          flowId,
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
      this.isLoading = false;
      setTimeout(() => {
        this.isBootstrapped = false;

        this.bootstrap();
      }, 300);
    }
  }

  async createBulkByLinkedIn({
    linkedInUrls,
    options,
    flowId,
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
          flowId,
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
      this.isLoading = false;

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
