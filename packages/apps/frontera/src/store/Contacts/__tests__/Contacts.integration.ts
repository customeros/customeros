import { it, expect, describe } from 'vitest';

import { Transport } from '../../transport';
import { ContactService } from '../../Contacts/__service__/Contacts.service';
import { OrganizationsService } from '../../Organizations/__service__/Organizations.service';

const transport = new Transport();
const organizationsService = OrganizationsService.getInstance(transport);
const contactService = ContactService.getInstance(transport);

describe('ContactsService - Integration Tests', () => {
  it('create contact', async () => {
    const contact_social_url = 'IT_' + crypto.randomUUID();

    const { contact_Create } = await contactService.createContact({
      contactInput: { socialUrl: contact_social_url },
    });
    const { contact } = await contactService.getContact(contact_Create);

    expect(contact?.socials.length).toBe(1);
    expect(contact?.socials[0].url).toBe(contact_social_url);
    expect(contact?.createdAt).not.toBeNull();
  });

  it('create contact for organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });
    const contact_social_url = 'IT_' + crypto.randomUUID();

    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: organization_Save.metadata.id,
        input: { socialUrl: contact_social_url },
      });
    const { contact } = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect(contact?.socials.length).toBe(1);
    expect(contact?.socials[0].url).toBe(contact_social_url);
    expect(contact?.organizations.content.length).toBe(1);
    expect(contact?.organizations.content[0].name).toBe(organization_name);
  });

  it('update contact', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });
    const contact_social_url = 'IT_' + crypto.randomUUID();
    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: organization_Save.metadata.id,
        input: { socialUrl: contact_social_url },
      });

    const contact_name = 'IT_' + crypto.randomUUID();
    const contact_description = 'IT_' + crypto.randomUUID();
    const contact_firstName = 'IT_' + crypto.randomUUID();
    const contact_lastName = 'IT_' + crypto.randomUUID();
    const contact_prefix = 'Mr.';
    const contact_profilePhotoUrl = 'https://example.com';
    const contact_timezone = 'America/North_Dakota/New_Salem';
    const contact_username = 'zzzzzz';

    await contactService.updateContact({
      input: {
        id: contact_CreateForOrganization.id,
        name: contact_name,
        description: contact_description,
        firstName: contact_firstName,
        lastName: contact_lastName,
        prefix: contact_prefix,
        patch: true,
        profilePhotoUrl: contact_profilePhotoUrl,
        timezone: contact_timezone,
        username: contact_username,
      },
    });

    const { contact } = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect.soft(contact?.name).toBe(contact_name);
    expect.soft(contact?.description).toBe(contact_description);
    expect.soft(contact?.firstName).toBe(contact_firstName);
    expect.soft(contact?.lastName).toBe(contact_lastName);
    expect.soft(contact?.prefix).toBe(contact_prefix);
    expect.soft(contact?.profilePhotoUrl).toBe(contact_profilePhotoUrl);
    expect.soft(contact?.timezone).toBe(contact_timezone);
    expect(contact?.updatedAt).not.toBeNull();
    //TODO: expect.soft(contact?.username).toBe(username) when username is implemented;
  });

  it('gets contacts', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });

    const contact_first_social_url = 'IT_' + crypto.randomUUID();

    const firstContact = await contactService.createContactForOrganization({
      organizationId: organization_Save.metadata.id,
      input: { socialUrl: contact_first_social_url },
    });

    const contact_second_social_url = 'IT_' + crypto.randomUUID();

    const secondContact = await contactService.createContactForOrganization({
      organizationId: organization_Save.metadata.id,
      input: { socialUrl: contact_second_social_url },
    });

    const contact_third_social_url = 'IT_' + crypto.randomUUID();

    const thirdContact = await contactService.createContactForOrganization({
      organizationId: organization_Save.metadata.id,
      input: { socialUrl: contact_third_social_url },
    });

    const { contacts } = await contactService.getContacts({
      pagination: { limit: 1000, page: 0 },
    });
    const contactIds = contacts?.content?.map((contact) => contact.id);

    expect(contactIds).toBeDefined();
    expect(contactIds).toContain(firstContact.contact_CreateForOrganization.id);
    expect(contactIds).toContain(
      secondContact.contact_CreateForOrganization.id,
    );
    expect(contactIds).toContain(thirdContact.contact_CreateForOrganization.id);
  });

  it('links contact to organization', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });

    const contact_social_url = 'IT_' + crypto.randomUUID();

    const { contact_Create } = await contactService.createContact({
      contactInput: { socialUrl: contact_social_url },
    });

    const contactBeforeLink = await contactService.getContact(contact_Create);

    expect(contactBeforeLink.contact?.organizations.content.length).toBe(0);

    await contactService.linkOrganization({
      input: {
        organizationId: organization_Save.metadata.id,
        contactId: contactBeforeLink.contact!.metadata.id,
      },
    });

    const contactAfterLink = await contactService.getContact(contact_Create);

    expect(contactAfterLink.contact?.organizations.content.length).toBe(1);
  });

  it('adds job roles to contact', async () => {
    const organization_name = 'IT_' + crypto.randomUUID();
    const { organization_Save } = await organizationsService.saveOrganization({
      input: { name: organization_name },
    });
    const contact_social_url = 'IT_' + crypto.randomUUID();

    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: organization_Save.metadata.id,
        input: { socialUrl: contact_social_url },
      });

    const contactBeforeFirstJobRole = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect(contactBeforeFirstJobRole.contact?.jobRoles.length).toBe(1);
    expect(contactBeforeFirstJobRole.contact?.jobRoles[0].jobTitle).toBeNull();
    expect(contactBeforeFirstJobRole.contact?.jobRoles[0].description).toBe('');
    expect(contactBeforeFirstJobRole.contact?.jobRoles[0].primary).toBe(false);
    expect(contactBeforeFirstJobRole.contact?.jobRoles[0].company).toBeNull();
    expect(contactBeforeFirstJobRole.contact?.jobRoles[0].startedAt).toBeNull();
    expect(contactBeforeFirstJobRole.contact?.jobRoles[0].endedAt).toBeNull();

    const jobRoleOneDescription = 'IT_' + crypto.randomUUID();
    const jobRoleOneTitle = 'IT_' + crypto.randomUUID();
    const jobRoleOneStartedAt = new Date().toISOString();

    await contactService.addJobRole({
      contactId: contactBeforeFirstJobRole.contact!.metadata.id,
      input: {
        description: jobRoleOneDescription,
        jobTitle: jobRoleOneTitle,
        organizationId: organization_Save.metadata.id,
        startedAt: jobRoleOneStartedAt,
      },
    });

    const contactAfterFirstJobRole = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect(contactAfterFirstJobRole.contact?.jobRoles.length).toBe(2);
    expect(contactAfterFirstJobRole.contact?.jobRoles[1].jobTitle).toBe(
      jobRoleOneTitle,
    );
    expect(contactAfterFirstJobRole.contact?.jobRoles[1].jobTitle).toBe(
      jobRoleOneTitle,
    );
    expect(contactAfterFirstJobRole.contact?.jobRoles[1].description).toBe(
      jobRoleOneDescription,
    );
    expect(contactAfterFirstJobRole.contact?.jobRoles[1].primary).toBe(false);
    expect(contactAfterFirstJobRole.contact?.jobRoles[1].company).toBeNull();
    expect(
      new Date(
        contactAfterFirstJobRole.contact?.jobRoles[1].startedAt,
      ).getTime(),
    ).toBe(new Date(jobRoleOneStartedAt).getTime());
    expect(contactAfterFirstJobRole.contact?.jobRoles[0].endedAt).toBeNull();
  });
});
