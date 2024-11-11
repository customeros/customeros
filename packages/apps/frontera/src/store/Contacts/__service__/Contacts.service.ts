import type { Transport } from '@store/transport.ts';

import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { rdiffResult } from 'recursive-diff';

import {
  Tag,
  ContactUpdateInput,
} from '@shared/types/__generated__/graphql.types';

import { ContactStore } from '../Contact.store';
import AddJobRoleDocument from './addJobRole.graphql';
import ContactQueryDocument from './getContact.graphql';
import ContactsQueryDocument from './getContacts.graphql';
import UpdateContactDocument from './contactUpdate.graphql';
import UpdateContactRoleDocument from './updateJobRole.graphql';
import UpdateContactEmailDocument from './emailReplace.graphql';
import AddContactSocialDocument from './addContactSocial.graphql';
import CreateContactMutationDocument from './createContact.graphql';
import DeleteContactMutationDocument from './deleteContact.graphql';
import LinkOrganizationDocument from './linkContactWithOrg.graphql';
import ArchiveContactMutationDocument from './archiveContact.graphql';
import AddTagsToContactMutationDocument from './addTagsToContact.graphql';
import AddContactPhoneNumberDocument from './addContactPhoneNumber.graphql';
import { ContactQuery, ContactQueryVariables } from './getContact.generated';
import UpdateContactSocialMutationDocument from './updateContactSocial.graphql';
import { ContactsQuery, ContactsQueryVariables } from './getContacts.generated';
import CreateContactForOrgMutationDocument from './createContactForOrg.graphql';
import UpdateContactPhoneNumberDocument from './updateContactPhoneNumber.graphql';
import RemoveContactPhoneNumberDocument from './removeContactPhoneNumber.graphql';
import FindWorkContactEmailMutationDocument from './findWorkContactEmail.graphql';
import RemoveTagsFromContactMutationDocument from './removeTagsFromContact.graphql';
import SetPrimaryEmailForContactMutationDocument from './setPrimaryEmailForContact.graphql';
import {
  AddJobRoleMutation,
  AddJobRoleMutationVariables,
} from './addJobRole.generated';
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
  CreateContactMutation as CreateContactForOrgMutation,
  CreateContactMutationVariables as CreateContactForOrgMutationVariables,
} from './createContactForOrg.generated';
class ContactService {
  private static instance: ContactService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): ContactService {
    if (!ContactService.instance) {
      ContactService.instance = new ContactService(transport);
    }

    return ContactService.instance;
  }

  async getContact(contactId: string) {
    return this.transport.graphql.request<ContactQuery, ContactQueryVariables>(
      ContactQueryDocument,
      { id: contactId },
    );
  }

  async getContacts(payload: ContactsQueryVariables) {
    return this.transport.graphql.request<
      ContactsQuery,
      ContactsQueryVariables
    >(ContactsQueryDocument, payload);
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

  public async mutateOperation(operation: Operation, store: ContactStore) {
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
      .with(['organizations', 'content', ...P.array()], () => {
        if (type === 'add') {
          this.linkOrganization({
            input: {
              contactId: contactId!,
              organizationId: value.metadata.id,
            },
          });
        }
      })
      .with(['phoneNumbers', 0, ...P.array()], () => {
        if (type === 'add') {
          this.addPhoneNumber({
            contactId: contactId!,
            input: {
              phoneNumber: value.rawPhoneNumber,
            },
          });
        }

        if (type === 'update') {
          this.updatePhoneNumber({
            input: {
              id: store.value.phoneNumbers[0].id,
              phoneNumber: store.value.phoneNumbers[0].rawPhoneNumber || '',
            },
          });
        }
      })
      .with(['socials', ...P.array()], ([_]) => {
        if (type === 'add') {
          this.addSocial({
            contactId: contactId!,
            input: {
              url: value.url,
            },
          });
        }

        if (type === 'update') {
          this.updateSocial({
            input: {
              id: store.value.socials[0].id,
              url: store.value.socials[0].url,
            },
          });
        }
      })
      .with(['jobRoles', 0, ...P.array()], () => {
        if (type === 'add') {
          this.addJobRole({
            contactId: contactId!,
            input: {
              description: store.value.jobRoles[0].description,
              jobTitle: store.value.jobRoles[0].jobTitle,
            },
          });
        }

        if (type === 'update') {
          this.updateJobRole({
            contactId: contactId!,
            input: {
              id: store.value.jobRoles[0].id,
              description: store.value.jobRoles[0].description,
              jobTitle: store.value.jobRoles[0].jobTitle,
            },
          });
        }
      })
      .with(['emails', 0, ...P.array()], () => {
        if (type === 'update') {
          const findIndex = store.value.emails.findIndex(
            (email) => email === value,
          );

          this.updateContactEmail({
            contactId: contactId!,
            input: {
              email: value.email,
              primary: store.value.emails[findIndex].primary || false,
            },
            previousEmail: oldValue as string,
          });
        }

        if (type === 'add') {
          this.updateContactEmail({
            contactId: contactId!,
            input: {
              email: value.email,
              primary: true,
            },
            previousEmail: '',
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
                id: value.id,
                name: value.name,
              },
            },
          });
        }

        if (type === 'delete') {
          if (typeof oldValue === 'object') {
            this.removeTagsFromContact({
              input: { contactId: contactId!, tag: { id: oldValue.id } },
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
                    id: tag.id,
                    name: tag.name,
                  },
                },
              });
            });
          }

          if (oldValue) {
            this.removeTagsFromContact({
              input: { contactId: contactId!, tag: { id: oldValue.id } },
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
