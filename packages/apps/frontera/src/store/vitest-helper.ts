// import { expect } from 'vitest';
import { OrganizationRepository } from '@infra/repositories/organization';
import { ContactService } from '@store/Contacts/__service__/Contacts.service.ts';
import { contactsTestState } from '@store/Contacts/__tests__/contactsTestState.ts';
import { organizationsTestState } from '@store/Organizations/__tests__/organizationsTestState.ts';

import { SortingDirection } from '@graphql/types';

interface ContactCreateInput {
  name?: string;
  socialUrl?: string;
  email?: {
    email: string;
    primary: boolean;
  };
}

export class VitestHelper {
  static async createOrganizationForTest(
    organizationRepository: OrganizationRepository,
    input?: {
      input: { name?: string; ownerId?: string; organizationId?: string };
    },
  ) {
    const organization_name = 'vitest-' + crypto.randomUUID();
    const organization_domain = 'vitest-' + crypto.randomUUID() + '.com';
    const { organization_Save } = await organizationRepository.saveOrganization(
      input || {
        input: { name: organization_name, domains: [organization_domain] },
      },
    );
    const { metadata } = organization_Save;

    trackOrganization(metadata.id);

    return {
      organizationId: metadata.id,
      organizationName: organization_name,
    };
  }

  static async createContactForTest(
    contactService: ContactService,
    input?: {
      input: {
        id?: string;
        lastName?: string;
        firstName?: string;
        socialUrl?: string;
      };
    },
  ) {
    const first_name = 'vitest-' + crypto.randomUUID();
    const last_name = 'vitest-' + crypto.randomUUID();
    const contact_social_url =
      'https://www.linkedin.com/in/Vitest-' + crypto.randomUUID();
    const email = `${first_name}@${crypto.randomUUID()}.com`;

    const { contact_Create: contactId } = await contactService.createContact({
      contactInput: input?.input || {
        firstName: first_name,
        lastName: last_name,
        email: {
          email,
          primary: true,
        },
        socialUrl: contact_social_url,
      },
    });

    trackContact(contactId);

    return {
      contactId: contactId,
      firstName: first_name,
      lastName: last_name,
      email,
    };
  }

  static async createContactForOrganizationForTest(
    contactService: ContactService,
    organizationId: string,
    input: {
      name?: string;
      email?: string;
      socialUrl?: string;
    },
  ) {
    const contactInput: ContactCreateInput = {};

    if (input.socialUrl) {
      // Create by LinkedIn
      contactInput.socialUrl = input.socialUrl;
    } else if (input.email) {
      // Create by Email
      contactInput.email = {
        email: input.email,
        primary: true,
      };
    } else if (input.name) {
      // Create by Name
      contactInput.name = input.name;
    }

    const {
      contact_CreateForOrganization: { id: contactId },
    } = await contactService.createContactForOrganization({
      input: contactInput,
      organizationId,
    });

    trackContact(contactId);

    return {
      contactId,
      ...input,
      organizationId,
    };
  }

  static async createContactBulkByLinkedInForTest(
    contactService: ContactService,
    linkedInUrls: string[],
  ) {
    const { contact_CreateBulkByLinkedIn: contactIds } =
      await contactService.createContactBulkByLinkedIn({
        linkedInUrls,
      });

    const { ui_contacts_search } = await contactService.searchContacts({
      where: {
        AND: [], // Empty AND array for no specific filters
      },
      sort: {
        by: 'CONTACTS_UPDATED_AT',
        direction: SortingDirection.Desc,
      },
    });
    const createContacts = (
      await contactService.getContactsByIds({ ids: ui_contacts_search.ids })
    ).ui_contacts.filter(
      (contact) =>
        contact.linkedInUrl !== null &&
        contact.linkedInUrl !== undefined &&
        linkedInUrls.includes(contact.linkedInUrl),
    );

    createContacts.forEach((createContacts) => trackContact(createContacts.id));

    return {
      contactIds,
      linkedInUrls,
    };
  }

  static async createContactBulkByEmailForTest(
    contactService: ContactService,
    emails: string[],
  ) {
    const { contact_CreateBulkByEmail: contactIds } =
      await contactService.createContactBulkByEmail({
        emails,
      });

    const { ui_contacts_search } = await contactService.searchContacts({
      where: {
        AND: [],
      },
      sort: {
        by: 'CONTACTS_UPDATED_AT',
        direction: SortingDirection.Desc,
      },
    });

    const createContacts = (
      await contactService.getContactsByIds({ ids: ui_contacts_search.ids })
    ).ui_contacts.filter((contact) =>
      contact.emails.some(
        (emailObj) =>
          emailObj.email !== null &&
          emailObj.email !== undefined &&
          emails.includes(emailObj.email),
      ),
    );

    createContacts.forEach((createContacts) => trackContact(createContacts.id));

    return {
      contactIds,
      emails,
    };
  }
}

export const trackOrganization = (organizationId: string) => {
  organizationsTestState.createdOrganizationIds.add(organizationId);
};

export const trackContact = (contactId: string) => {
  contactsTestState.createdContactsIds.add(contactId);
};

export const generateEmail = (): string => {
  const letters = 'abcdefghijklmnopqrstuvwxyz';
  const firstPart = Array.from({ length: 5 }, () =>
    letters.charAt(Math.floor(Math.random() * letters.length)),
  ).join('');
  const secondPart = Array.from({ length: 4 }, () =>
    letters.charAt(Math.floor(Math.random() * letters.length)),
  ).join('');

  return `${firstPart}.${secondPart}@msn.com`;
};
