import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { rdiffResult } from 'recursive-diff';
import { Transport } from '@infra/transport.ts';

import {
  Tag,
  ContactUpdateInput,
} from '@shared/types/__generated__/graphql.types';

import type { Contact } from '../Contact.dto.ts';

import AddJobRoleDocument from './addJobRole.graphql';
import UpdateContactDocument from './contactUpdate.graphql';
import ContactsByIdDocument from './getContactsById.graphql';
import SearchContactsDocument from './searchContacts.graphql';
import UpdateContactRoleDocument from './updateJobRole.graphql';
import UpdateContactEmailDocument from './emailReplace.graphql';
import AddContactSocialDocument from './addContactSocial.graphql';
import CreateContactMutationDocument from './createContact.graphql';
import DeleteContactMutationDocument from './deleteContact.graphql';
import LinkOrganizationDocument from './linkContactWithOrg.graphql';
import ContactByEmailDocument from './contactExistsByEmail.graphql';
import ArchiveContactMutationDocument from './archiveContact.graphql';
import AddTagsToContactMutationDocument from './addTagsToContact.graphql';
import AddContactPhoneNumberDocument from './addContactPhoneNumber.graphql';
import ContactExistsByLinkedInDocument from './contactExistsByLinkedIn.graphql';
import UpdateContactSocialMutationDocument from './updateContactSocial.graphql';
import CreateContactForOrgMutationDocument from './createContactForOrg.graphql';
import UpdateContactPhoneNumberDocument from './updateContactPhoneNumber.graphql';
import RemoveContactPhoneNumberDocument from './removeContactPhoneNumber.graphql';
import FindWorkContactEmailMutationDocument from './findWorkContactEmail.graphql';
import RemoveTagsFromContactMutationDocument from './removeTagsFromContact.graphql';
import CreateContactBulkByEmailMutationDocument from './createContactBulkByEmail.graphql';
import SetPrimaryEmailForContactMutationDocument from './setPrimaryEmailForContact.graphql';
import {
  AddJobRoleMutation,
  AddJobRoleMutationVariables,
} from './addJobRole.generated';
import CreateContactBulkByLinkedInMutationDocument from './createContactBulkByLinkedIn.graphql';
import {
  SearchContactsQuery,
  SearchContactsQueryVariables,
} from './searchContacts.generated.ts';
import {
  CreateContactMutation,
  CreateContactMutationVariables,
} from './createContact.generated';
import {
  DeleteContactMutation,
  DeleteContactMutationVariables,
} from './deleteContact.generated';
import {
  UpdateContactMutation,
  UpdateContactMutationVariables,
} from './contactUpdate.generated';
import {
  ArchiveContactMutation,
  ArchiveContactMutationVariables,
} from './archiveContact.generated';
import {
  GetContactsByIdsQuery,
  GetContactsByIdsQueryVariables,
} from './getContactsById.generated.ts';
import {
  ContactByEmailQuery,
  ContactByEmailQueryVariables,
} from './contactExistsByEmail.generated.ts';
import {
  UpdateContactRoleMutation,
  UpdateContactRoleMutationVariables,
} from './updateJobRole.generated';
import {
  AddContactSocialMutation,
  AddContactSocialMutationVariables,
} from './addContactSocial.generated';
import {
  AddTagsToContactMutation,
  AddTagsToContactMutationVariables,
} from './addTagsToContact.generated';
import {
  UpdateContactEmailMutation,
  UpdateContactEmailMutationVariables,
} from './emailReplace.generated';
import {
  LinkOrganizationMutation,
  LinkOrganizationMutationVariables,
} from './linkContactWithOrg.generated';
import {
  UpdateContactSocialMutation,
  UpdateContactSocialMutationVariables,
} from './updateContactSocial.generated';
import {
  FindWorkContactEmailMutation,
  FindWorkContactEmailMutationVariables,
} from './findWorkContactEmail.generated';
import {
  RemoveTagFromContactMutation,
  RemoveTagFromContactMutationVariables,
} from './removeTagsFromContact.generated';
import {
  AddContactPhoneNumberMutation,
  AddContactPhoneNumberMutationVariables,
} from './addContactPhoneNumber.generated';
import {
  ContactExistsByLinkedInQuery,
  ContactExistsByLinkedInQueryVariables,
} from './contactExistsByLinkedIn.generated.ts';
import {
  RemoveContactPhoneNumberMutation,
  RemoveContactPhoneNumberMutationVariables,
} from './removeContactPhoneNumber.generated';
import {
  UpdateContactPhoneNumberMutation,
  UpdateContactPhoneNumberMutationVariables,
} from './updateContactPhoneNumber.generated';
import {
  SetPrimaryEmailForContactMutation,
  SetPrimaryEmailForContactMutationVariables,
} from './setPrimaryEmailForContact.generated';
import {
  CreateContactBulkByEmailMutation,
  CreateContactBulkByEmailMutationVariables,
} from './createContactBulkByEmail.generated.ts';
import {
  CreateContactBulkByLinkedInMutation,
  CreateContactBulkByLinkedInMutationVariables,
} from './createContactBulkByLinkedIn.generated.ts';
import {
  CreateContactMutation as CreateContactForOrgMutation,
  CreateContactMutationVariables as CreateContactForOrgMutationVariables,
} from './createContactForOrg.generated';
class ContactService {
  private static instance: ContactService | null = null;
  private transport: Transport = Transport.getInstance();

  constructor() {}

  static getInstance(): ContactService {
    if (!ContactService.instance) {
      ContactService.instance = new ContactService();
    }

    return ContactService.instance;
  }

  async getContact(id: string) {
    const { ui_contacts } = await this.getContactsByIds({
      ids: [id],
    });

    return ui_contacts[0];
  }

  async getContactsByIds(payload: GetContactsByIdsQueryVariables) {
    return this.transport.graphql.request<
      GetContactsByIdsQuery,
      GetContactsByIdsQueryVariables
    >(ContactsByIdDocument, payload);
  }

  async searchContacts(payload: SearchContactsQueryVariables) {
    return this.transport.graphql.request<
      SearchContactsQuery,
      SearchContactsQueryVariables
    >(SearchContactsDocument, payload);
  }

  async createContact(payload: CreateContactMutationVariables) {
    return this.transport.graphql.request<
      CreateContactMutation,
      CreateContactMutationVariables
    >(CreateContactMutationDocument, payload);
  }

  async createContactForOrganization(
    payload: CreateContactForOrgMutationVariables,
  ) {
    return this.transport.graphql.request<
      CreateContactForOrgMutation,
      CreateContactForOrgMutationVariables
    >(CreateContactForOrgMutationDocument, payload);
  }

  async linkOrganization(payload: LinkOrganizationMutationVariables) {
    return this.transport.graphql.request<
      LinkOrganizationMutation,
      LinkOrganizationMutationVariables
    >(LinkOrganizationDocument, payload);
  }

  async updateContact(payload: UpdateContactMutationVariables) {
    return this.transport.graphql.request<
      UpdateContactMutation,
      UpdateContactMutationVariables
    >(UpdateContactDocument, payload);
  }

  async addJobRole(payload: AddJobRoleMutationVariables) {
    return this.transport.graphql.request<
      AddJobRoleMutation,
      AddJobRoleMutationVariables
    >(AddJobRoleDocument, payload);
  }

  async updateJobRole(payload: UpdateContactRoleMutationVariables) {
    return this.transport.graphql.request<
      UpdateContactRoleMutation,
      UpdateContactRoleMutationVariables
    >(UpdateContactRoleDocument, payload);
  }

  async updateContactEmail(payload: UpdateContactEmailMutationVariables) {
    return this.transport.graphql.request<
      UpdateContactEmailMutation,
      UpdateContactEmailMutationVariables
    >(UpdateContactEmailDocument, payload);
  }

  async addPhoneNumber(payload: AddContactPhoneNumberMutationVariables) {
    return this.transport.graphql.request<
      AddContactPhoneNumberMutation,
      AddContactPhoneNumberMutationVariables
    >(AddContactPhoneNumberDocument, payload);
  }

  async updatePhoneNumber(payload: UpdateContactPhoneNumberMutationVariables) {
    return this.transport.graphql.request<
      UpdateContactPhoneNumberMutation,
      UpdateContactPhoneNumberMutationVariables
    >(UpdateContactPhoneNumberDocument, payload);
  }

  async removePhoneNumber(payload: RemoveContactPhoneNumberMutationVariables) {
    return this.transport.graphql.request<
      RemoveContactPhoneNumberMutation,
      RemoveContactPhoneNumberMutationVariables
    >(RemoveContactPhoneNumberDocument, payload);
  }

  async addSocial(payload: AddContactSocialMutationVariables) {
    return this.transport.graphql.request<
      AddContactSocialMutation,
      AddContactSocialMutationVariables
    >(AddContactSocialDocument, payload);
  }

  async updateSocial(payload: UpdateContactSocialMutationVariables) {
    return this.transport.graphql.request<
      UpdateContactSocialMutation,
      UpdateContactSocialMutationVariables
    >(UpdateContactSocialMutationDocument, payload);
  }

  async findEmail(payload: FindWorkContactEmailMutationVariables) {
    return this.transport.graphql.request<
      FindWorkContactEmailMutation,
      FindWorkContactEmailMutationVariables
    >(FindWorkContactEmailMutationDocument, payload);
  }

  async setPrimaryEmail(payload: SetPrimaryEmailForContactMutationVariables) {
    return this.transport.graphql.request<
      SetPrimaryEmailForContactMutation,
      SetPrimaryEmailForContactMutationVariables
    >(SetPrimaryEmailForContactMutationDocument, payload);
  }

  async deleteContact(payload: DeleteContactMutationVariables) {
    return this.transport.graphql.request<
      DeleteContactMutation,
      DeleteContactMutationVariables
    >(DeleteContactMutationDocument, payload);
  }

  async archiveContact(payload: ArchiveContactMutationVariables) {
    return this.transport.graphql.request<
      ArchiveContactMutation,
      ArchiveContactMutationVariables
    >(ArchiveContactMutationDocument, payload);
  }

  async addTagsToContact(payload: AddTagsToContactMutationVariables) {
    return this.transport.graphql.request<
      AddTagsToContactMutation,
      AddTagsToContactMutationVariables
    >(AddTagsToContactMutationDocument, payload);
  }

  async removeTagsFromContact(payload: RemoveTagFromContactMutationVariables) {
    return this.transport.graphql.request<
      RemoveTagFromContactMutation,
      RemoveTagFromContactMutationVariables
    >(RemoveTagsFromContactMutationDocument, payload);
  }

  async createContactBulkByEmail(
    payload: CreateContactBulkByEmailMutationVariables,
  ) {
    return this.transport.graphql.request<
      CreateContactBulkByEmailMutation,
      CreateContactBulkByEmailMutationVariables
    >(CreateContactBulkByEmailMutationDocument, payload);
  }

  async createContactBulkByLinkedIn(
    payload: CreateContactBulkByLinkedInMutationVariables,
  ) {
    return this.transport.graphql.request<
      CreateContactBulkByLinkedInMutation,
      CreateContactBulkByLinkedInMutationVariables
    >(CreateContactBulkByLinkedInMutationDocument, payload);
  }

  async contactExistsByEmail(payload: ContactByEmailQueryVariables) {
    return this.transport.graphql.request<
      ContactByEmailQuery,
      ContactByEmailQueryVariables
    >(ContactByEmailDocument, payload);
  }

  async contactExistsByLinkedIn(
    payload: ContactExistsByLinkedInQueryVariables,
  ) {
    return this.transport.graphql.request<
      ContactExistsByLinkedInQuery,
      ContactExistsByLinkedInQueryVariables
    >(ContactExistsByLinkedInDocument, payload);
  }

  public async mutateOperation(operation: Operation, store: Contact) {
    const diff = operation.diff?.[0];
    const type = diff?.op;
    const path = diff?.path;
    const value = diff?.val;
    const oldValue = (diff as rdiffResult & { oldVal: unknown })?.oldVal;

    const contactId = operation.entityId;

    if (!operation.diff.length) {
      return;
    }

    if (!contactId) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }
    match(path)
      .with(
        [
          P.union('primaryOrganizationId', 'primaryOrganizationName'),
          ...P.array(),
        ],
        async () => {
          await this.linkOrganization({
            input: {
              contactId,
              organizationId: value,
            },
          });
          store.store.root.ui.toastSuccess(
            "Contact's company was changed",
            'org-linked',
          );
        },
      )
      .with(['linkedInUrl', ...P.array()], async () => {
        await this.addSocial({
          contactId: contactId!,
          input: {
            url: value,
          },
        });
      })
      .with(['primaryOrganizationJobRoleTitle', ...P.array()], () => {
        if (store.value.primaryOrganizationJobRoleId?.length === 0) {
          this.addJobRole({
            contactId: contactId!,
            input: {
              description: store.value.primaryOrganizationJobRoleDescription,
              jobTitle: store.value.primaryOrganizationJobRoleTitle,
            },
          });
        } else {
          this.updateJobRole({
            contactId: contactId!,
            input: {
              id: store.value.primaryOrganizationJobRoleId || '',
              description: store.value.primaryOrganizationJobRoleDescription,
              jobTitle: store.value.primaryOrganizationJobRoleTitle,
            },
          });
        }
      })
      .with([...P.array(), 'primary'], () => {
        if (type === 'update') {
          this.setPrimaryEmail({
            contactId: contactId!,
            email:
              store.value.emails.find((email) => email.primary)?.email || '',
          });
        }
      })
      .with(['emails', ...P.array()], async () => {
        if (type === 'add') {
          await this.updateContactEmail({
            contactId: contactId!,
            input: {
              email: value.email,
            },
            previousEmail: '',
          });
        }

        if (type === 'update') {
          const findIndex = store.value.emails.findIndex(
            (email) => email.email === value,
          );

          this.updateContactEmail({
            contactId: contactId!,
            input: {
              email: value,
              primary: store.value.emails[findIndex]?.primary || false,
            },
            previousEmail: oldValue as string,
          });
        }

        if (type === 'delete') {
          this.updateContactEmail({
            contactId: contactId!,
            input: {
              email: '',
            },
            previousEmail: oldValue.email,
          });
        }
      })
      .with(['tags', ...P.array()], () => {
        if (type === 'add') {
          this.addTagsToContact({
            input: {
              contactId: contactId!,
              tag: {
                name: value.name,
              },
            },
          });
        }

        if (
          operation.diff.length > 1 &&
          operation.diff[operation.diff.length - 1].op === 'delete'
        ) {
          const pathIsId = match(operation.diff[0]?.path)
            .with([...P.array(), 'id'], () => true)
            .otherwise(() => false);

          const toBeDeletedName = !pathIsId
            ? (operation.diff[0] as rdiffResult & { oldVal: string })?.oldVal
            : store.store.root.tags.getById(
                (operation.diff[0] as rdiffResult & { oldVal: string })?.oldVal,
              )?.tagName;

          this.removeTagsFromContact({
            input: {
              contactId,
              tag: {
                name: toBeDeletedName,
              },
            },
          });

          return;
        }

        if (type === 'delete') {
          if (typeof oldValue === 'object') {
            this.removeTagsFromContact({
              input: {
                contactId: contactId!,
                tag: { name: oldValue.name },
              },
            });
          }
        }

        // if tag with index different that last one is deleted it comes as an update, bulk creation updates also come as updates
        if (type === 'update') {
          if (!oldValue) {
            (value as Array<Tag>)?.forEach((tag: Tag) => {
              this.addTagsToContact({
                input: {
                  contactId: contactId!,
                  tag: {
                    name: tag.name,
                  },
                },
              });
            });
          }

          if (oldValue) {
            this.removeTagsFromContact({
              input: { contactId: contactId!, tag: { name: oldValue.name } },
            });
          }
        }
      })
      .otherwise(() => {
        const payload = makePayload<ContactUpdateInput>(operation);

        this.updateContact({ input: { ...payload, id: contactId } });
      });
  }
}

export { ContactService };
