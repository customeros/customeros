import { Store } from '@store/_store';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { action, computed, runInAction } from 'mobx';

import {
  Tag,
  ContactInput,
  Contact as ContactData,
} from '@shared/types/__generated__/graphql.types';

import { Contact, ContactDatum } from './Contact.dto';
import { ContactsView } from './__views__/Contacts.view';
import { ContactService } from './__service__/Contacts.service';
import { FlowContactsView } from './__views__/FlowContacts.view';

export class ContactsStore extends Store<ContactDatum, Contact> {
  private service: ContactService;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Contacts',
      getId: (data) => data?.metadata?.id,
      factory: Contact,
    });
    this.service = ContactService.getInstance(transport);

    new ContactsView(this);
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

  @action
  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;

      const { contacts } = await this.service.getContacts({
        pagination: { limit: 1000, page: 0 },
      });

      const data = contacts.content as ContactData[];
      const totalElements = contacts.totalElements;

      runInAction(() => {
        data.forEach((contact) => {
          if (!contact) return;
          const record = new Contact(this, contact);

          this.value.set(record.id, record);
        });
        this.size = this.value.size;

        if (this.totalElements !== totalElements) {
          this.totalElements = totalElements;
        }
      });

      await this.bootstrapRest();
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isBootstrapped = true;
      });
    }
  }

  @action
  async bootstrapRest() {
    let page = 1;

    while (this.totalElements > this.value.size) {
      try {
        const { contacts } = await this.service.getContacts({
          pagination: { limit: 1000, page },
        });

        const data = contacts.content as ContactData[];

        page++;
        runInAction(() => {
          data.forEach((contact) => {
            if (!contact) return;
            const record = new Contact(this, contact);

            this.value.set(record.id, record);
          });

          this.size = this.value.size;
        });
      } catch (e) {
        runInAction(() => {
          this.error = (e as Error)?.message;
        });
        break;
      }
    }

    runInAction(() => {
      this.isBootstrapped = this.totalElements === this.value.size;
      this.isBootstrapping = false;
    });
  }

  @action
  async create(
    organizationId: string,
    options?: { onSuccess?: (serverId: string) => void },
    input?: ContactInput,
  ) {
    const newContact = new Contact(this, Contact.default());
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
