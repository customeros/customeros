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
  DataSource,
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

    const newContact = new ContactStore(
      this.root,
      this.transport,
      ContactStore.getDefaultValue(),
    );

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
        await this.service.createContactForOrganization({
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

    const newContact = new ContactStore(
      this.root,
      this.transport,
      ContactStore.getDefaultValue(),
    );
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
      const { contact_Create } = await this.service.createContact({
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

  updateTags = (ids: string[], tags: Tag[]) => {
    const tagIdsToUpdate = new Set(tags?.map((tag) => tag.id));

    const shouldRemoveTags = ids.every((id) => {
      const contact = this.value.get(id);

      if (!contact) return false;

      const contactIdsTags = new Set(
        (contact.value.tags ?? []).map((tag) => tag.id),
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
          (t) => !tagIdsToUpdate.has(t.id),
        );
      } else {
        const existingIds = new Set(contact.value.tags?.map((t) => t.id) ?? []);
        const newTags = tags.filter((t) => !existingIds.has(t.id));

        if (!Array.isArray(contact.value.tags)) {
          contact.value.tags = [];
        }

        contact.value.tags = [...(contact.value.tags ?? []), ...newTags];

        contact.commit();
      }
    });
  };
}
