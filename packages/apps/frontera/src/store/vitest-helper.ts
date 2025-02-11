// import { expect } from 'vitest';
import { OrganizationRepository } from '@infra/repositories/organization';
import { ContactService } from '@store/Contacts/__service__/Contacts.service.ts';
import { contactsTestState } from '@store/Contacts/__tests__/contactsTestState.ts';
import { organizationsTestState } from '@store/Organizations/__tests__/organizationsTestState.ts';

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

    const {
      contact_CreateForOrganization: { id: contactId },
    } = await contactService.createContactForOrganization({
      input: input?.input || {
        firstName: first_name,
        lastName: last_name,
        email: {
          email,
          primary: true,
        },
        socialUrl: contact_social_url,
      },
      organizationId,
    });

    trackContact(contactId);

    return {
      contactId: contactId,
      firstName: first_name,
      lastName: last_name,
      email,
      organizationId,
    };
  }
}

export const trackOrganization = (organizationId: string) => {
  organizationsTestState.createdOrganizationIds.add(organizationId);
};

export const trackContact = (contactId: string) => {
  contactsTestState.createdContactsIds.add(contactId);
};
