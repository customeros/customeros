import { it, expect, describe } from 'vitest';
import { VitestHelper, generateEmail } from '@store/vitest-helper.ts';
import { OrganizationRepository } from '@infra/repositories/organization/organization.repository.ts';

import {
  EmailLabel,
  EntityType,
  PhoneNumberLabel,
  SortingDirection,
} from '@graphql/types';

import { ContactService } from '../../Contacts/__service__/Contacts.service';
import { JobRolesService } from '../../JobRoles/__service__/JobRoles.service.ts';

const organizationsRepository = OrganizationRepository.getInstance();
const contactService = ContactService.getInstance();
const jobRolesService = JobRolesService.getInstance();

describe('ContactsService - Integration Tests', () => {
  it('create contact', async () => {
    const contact_social_url =
      'https://www.linkedin.com/in/Vitest-' + crypto.randomUUID();

    const { contactId } = await VitestHelper.createContactForTest(
      contactService,
      {
        input: {
          socialUrl: contact_social_url,
        },
      },
    );
    const contact = await contactService.getContact(contactId);

    expect(contact?.id).toBe(contactId);
    expect(contact?.createdAt).not.toBeNull();
    expect(contact?.createdAt).not.toBe('');
    expect(contact?.updatedAt).not.toBeNull();
    expect(contact?.firstName).toBe('');
    expect(contact?.lastName).toBe('');
    expect(contact?.name).toBe('');
    expect(contact?.prefix).toBe('');
    expect(contact?.description).toBe('');
    expect(contact?.timezone).toBe('');
    expect(contact?.profilePhotoUrl).toBe('');
    expect(contact?.enrichedAt).toBeNull();
    expect(contact?.enrichedFailedAt).toBeNull();
    // expect(contact?.enrichDetails.requestedAt).not.toBeNull(); asynchronous call so it generates false positive
    expect(contact?.enrichedEmailRequestedAt).toBeNull();
    expect(contact?.enrichedEmailEnrichedAt).toBeNull();
    expect(contact?.enrichedEmailFound).toBeNull();
    expect(contact?.linkedInInternalId).not.toBeNull;
    expect(Array.isArray(contact?.linkedInUrl)).toBe(false);
    expect(contact?.linkedInUrl).toBe(contact_social_url);
    expect(contact?.linkedInAlias).toBe('');
    expect(contact?.linkedInExternalId).toBe('');
    expect(contact?.linkedInFollowerCount).toBe(0);
    expect(contact?.primaryOrganizationId).toBeNull;
    expect(contact?.primaryOrganizationName).toBeNull;
    expect(contact?.primaryOrganizationJobRoleId).toBeNull;
    expect(contact?.primaryOrganizationJobRoleTitle).toBeNull;
    expect(contact?.primaryOrganizationJobRoleDescription).toBeNull;
    expect(contact?.primaryOrganizationJobRoleStartDate).toBeNull;
    expect(contact?.primaryOrganizationJobRoleEndDate).toBeNull;
    expect(contact?.emails.length).toBe(0);
    expect(contact?.phones.length).toBe(0);
    expect(contact?.tags).toEqual([]);
    expect(contact?.flows.length).toBe(0);
    expect(contact?.emails.find((e) => e.primary)).toBeUndefined;
    expect(contact?.locations.length).toBe(0);
    expect(contact?.emails.length).toBe(0);
    expect(contact?.connectedUsers.length).toBe(0);
  });

  it('creates contact by Linkedin for organization', async () => {
    const { organizationId, organizationName } =
      await VitestHelper.createOrganizationForTest(organizationsRepository);

    const { contactId } =
      await VitestHelper.createContactForOrganizationForTest(
        contactService,
        organizationId,
        {
          socialUrl:
            'https://www.linkedin.com/in/vitest-' + crypto.randomUUID(),
        },
      );

    const contact = await contactService.getContact(contactId);

    expect(contact?.primaryOrganizationName).toBe(organizationName);
    expect(contact?.linkedInUrl).toBeDefined();
    expect(contact?.emails.length).toBe(0);
  });

  it('creates contact by Email for organization', async () => {
    const { organizationId, organizationName } =
      await VitestHelper.createOrganizationForTest(organizationsRepository);

    const { contactId } =
      await VitestHelper.createContactForOrganizationForTest(
        contactService,
        organizationId,
        {
          email: `Vitest-${crypto.randomUUID()}@example.com`,
        },
      );

    const contact = await contactService.getContact(contactId);

    expect(contact?.primaryOrganizationName).toBe(organizationName);
    expect(contact?.emails.length).toBe(1);
    expect(contact?.linkedInUrl).toBeNull();
  });

  it('creates contact by Name for organization', async () => {
    const { organizationId, organizationName } =
      await VitestHelper.createOrganizationForTest(organizationsRepository);

    const { contactId } =
      await VitestHelper.createContactForOrganizationForTest(
        contactService,
        organizationId,
        {
          name: 'Vitest-' + crypto.randomUUID(),
        },
      );

    const contact = await contactService.getContact(contactId);

    expect(contact?.primaryOrganizationName).toBe(organizationName);
    expect(contact?.name).toBeDefined();
    expect(contact?.emails.length).toBe(0);
    expect(contact?.linkedInUrl).toBeNull();
  });

  it('update contact created in an organization', async () => {
    const { organizationId } = await VitestHelper.createOrganizationForTest(
      organizationsRepository,
    );
    const contact_social_url =
      'https://www.linkedin.com/in/vitest-' + crypto.randomUUID();

    const { contactId } =
      await VitestHelper.createContactForOrganizationForTest(
        contactService,
        organizationId,
        {
          socialUrl: contact_social_url,
        },
      );
    const contact_name = 'Vitest_' + crypto.randomUUID();
    const contact_description = 'Vitest_' + crypto.randomUUID();
    const contact_prefix = 'Mr.';
    const contact_profilePhotoUrl = 'https://example.com';
    const contact_timezone = 'America/North_Dakota/New_Salem';
    const contact_username = 'zzzzzz';

    await contactService.updateContact({
      input: {
        id: contactId,
        name: contact_name,
        description: contact_description,
        prefix: contact_prefix,
        patch: true,
        profilePhotoUrl: contact_profilePhotoUrl,
        timezone: contact_timezone,
        username: contact_username,
      },
    });

    const contact = await contactService.getContact(contactId);

    expect.soft(contact?.name).toBe(contact_name);
    expect.soft(contact?.description).toBe(contact_description);
    expect.soft(contact?.prefix).toBe(contact_prefix);
    expect.soft(contact?.profilePhotoUrl).toBe(contact_profilePhotoUrl);
    expect.soft(contact?.timezone).toBe(contact_timezone);
    expect(contact?.updatedAt).not.toBeNull();
    //TODO: expect.soft(contact?.username).toBe(username) when username is implemented;
  });

  it('gets contacts', async () => {
    const { organizationId } = await VitestHelper.createOrganizationForTest(
      organizationsRepository,
    );

    const contact_first_social_url =
      'https://www.linkedin.com/in/vitest-' + crypto.randomUUID();

    const firstContact = await VitestHelper.createContactForOrganizationForTest(
      contactService,
      organizationId,
      {
        socialUrl: contact_first_social_url,
      },
    );
    const firstContactId = firstContact.contactId;

    const contact_second_social_url =
      'https://www.linkedin.com/in/vitest-' + crypto.randomUUID();

    const secondContact =
      await VitestHelper.createContactForOrganizationForTest(
        contactService,
        organizationId,
        {
          socialUrl: contact_second_social_url,
        },
      );

    const secondContactId: string = secondContact.contactId;
    const contact_third_social_url =
      'https://www.linkedin.com/in/vitest-' + crypto.randomUUID();

    const thirdContact = await VitestHelper.createContactForOrganizationForTest(
      contactService,
      organizationId,
      {
        socialUrl: contact_third_social_url,
      },
    );

    const thirdContactId: string = thirdContact.contactId;
    const { ui_contacts } = await contactService.getContactsByIds({
      ids: [firstContactId, secondContactId, thirdContactId],
    });
    const contactIds = ui_contacts?.map((contact) => contact.id);

    expect(contactIds).toBeDefined();
    expect(contactIds).toContain(firstContactId);
    expect(contactIds).toContain(secondContactId);
    expect(contactIds).toContain(thirdContactId);
  });

  it('adds job roles to contact', async () => {
    const { organizationId, organizationName } =
      await VitestHelper.createOrganizationForTest(organizationsRepository);
    const contact_social_url =
      'https://www.linkedin.com/in/vitest-' + crypto.randomUUID();

    const contactId = await VitestHelper.createContactForOrganizationForTest(
      contactService,
      organizationId,
      {
        socialUrl: contact_social_url,
      },
    );

    const contactBeforeFirstJobRole = await contactService.getContact(
      contactId.contactId,
    );

    expect(contactBeforeFirstJobRole.primaryOrganizationJobRoleTitle).toBe('');
    expect(
      contactBeforeFirstJobRole.primaryOrganizationJobRoleDescription,
    ).toBe('');
    expect(
      contactBeforeFirstJobRole.primaryOrganizationJobRoleStartDate,
    ).toBeNull();
    expect(
      contactBeforeFirstJobRole.primaryOrganizationJobRoleEndDate,
    ).toBeNull();
    expect(contactBeforeFirstJobRole.primaryOrganizationName).toBe(
      organizationName,
    );

    const contactAfterCreation = await contactService.getContact(
      contactId.contactId,
    );
    const jobRoleBeforeUpdate = await jobRolesService.getJobRoles({
      ids: contactAfterCreation.jobRoleIds[0],
    });

    expect(jobRoleBeforeUpdate.jobRoles[0].description).toBe('');
    expect(jobRoleBeforeUpdate.jobRoles[0].jobTitle).toBeNull();
    expect(jobRoleBeforeUpdate.jobRoles[0].primary).toBe(true);
    expect(jobRoleBeforeUpdate.jobRoles[0].startedAt).toBe(null);
    expect(jobRoleBeforeUpdate.jobRoles[0].contact?.metadata.id).toBe(
      contactBeforeFirstJobRole.id,
    );

    const jobRoleCreateOneDescription = 'Vitest_' + crypto.randomUUID();
    const jobRoleCreateOneTitle = 'Vitest_' + crypto.randomUUID();
    const jobRoleCreateOneStartedAt = new Date().toISOString();

    await jobRolesService.saveJobRoles({
      input: {
        description: jobRoleCreateOneDescription,
        jobTitle: jobRoleCreateOneTitle,
        organizationId: organizationId,
        startedAt: jobRoleCreateOneStartedAt,
        contactId: contactBeforeFirstJobRole.id,
      },
    });

    const contactAfterFirstJobRole = await contactService.getContact(
      contactId.contactId,
    );

    expect(contactAfterFirstJobRole.jobRoleIds.length).toBe(1);

    const jobRole = await jobRolesService.getJobRoles({
      ids: contactAfterFirstJobRole.jobRoleIds[0],
    });

    expect(jobRole.jobRoles[0].description).toBe(jobRoleCreateOneDescription);
    expect(jobRole.jobRoles[0].jobTitle).toBe(jobRoleCreateOneTitle);
    expect(jobRole.jobRoles[0].primary).toBe(true);
    expect(new Date(jobRole.jobRoles[0].startedAt).toISOString()).toBe(
      new Date(jobRoleCreateOneStartedAt).toISOString(),
    );
    expect(jobRole.jobRoles[0].contact?.metadata.id).toBe(
      contactBeforeFirstJobRole.id,
    );
  });

  it('links contact to organization', async () => {
    const { organizationId, organizationName } =
      await VitestHelper.createOrganizationForTest(organizationsRepository);
    const { contactId } = await VitestHelper.createContactForTest(
      contactService,
    );

    await organizationsRepository.getOrganization(organizationId);

    const contactBeforeLink = await contactService.getContact(contactId);

    expect(contactBeforeLink.primaryOrganizationName).toBeNull;

    await contactService.linkOrganization({
      input: {
        contactId: contactId,
        organizationId: organizationId,
      },
    });
    expect(
      (await contactService.getContact(contactId)).primaryOrganizationId,
    ).toBe(organizationId);
    expect(
      (await contactService.getContact(contactId)).primaryOrganizationName,
    ).toBe(organizationName);
  });

  it('updates contact email', async () => {
    const { contactId } = await VitestHelper.createContactForTest(
      contactService,
    );
    const emailOne = 'Vitest_' + crypto.randomUUID() + '@example.com';

    await contactService.updateContactEmail({
      contactId: contactId,
      previousEmail: '',
      input: {
        email: emailOne,
        appSource: 'Vitest_test',
        label: EmailLabel.Personal,
      },
    });

    const contactAfterEmailUpdate = await contactService.getContact(contactId);

    expect(contactAfterEmailUpdate.emails.length).toBe(2);
    expect(contactAfterEmailUpdate.emails[0].id).not.toBeNull();
    expect(contactAfterEmailUpdate.emails[0].email).toBe(emailOne);
    expect(contactAfterEmailUpdate.emails[0].primary).toBe(false);
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.verified,
    ).toBe(false);
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails
        .verifyingCheckAll,
    ).toBe(false);
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.isValidSyntax,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.isRisky,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.isFirewalled,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.provider,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.firewall,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.isCatchAll,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.canConnectSmtp,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.deliverable,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.isMailboxFull,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.isRoleAccount,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.isFreeAccount,
    ).toBeNull();
    expect(
      contactAfterEmailUpdate.emails[0].emailValidationDetails.smtpSuccess,
    ).toBeNull();

    const emailTwo = 'Vitest_' + crypto.randomUUID() + '@example.com';

    await contactService.updateContactEmail({
      contactId: contactId,
      previousEmail: '',
      input: {
        email: emailTwo,
        appSource: 'Vitest_test',
        label: EmailLabel.Work,
      },
    });

    const contactAfterSecondEmailUpdate = await contactService.getContact(
      contactId,
    );

    expect(contactAfterSecondEmailUpdate?.emails.length).toBe(3);

    let emailOneEntry = contactAfterSecondEmailUpdate?.emails.find(
      (email) => email.email === emailOne,
    );
    let emailTwoEntry = contactAfterSecondEmailUpdate?.emails.find(
      (email) => email.email === emailTwo,
    );

    expect(emailOneEntry).toBeDefined();
    expect(emailTwoEntry).toBeDefined();

    expect(emailOneEntry?.id).not.toBeNull();
    expect(emailOneEntry?.primary).toBe(false);

    expect(emailTwoEntry?.id).not.toBeNull();
    expect(emailTwoEntry?.primary).toBe(false);

    await contactService.setPrimaryEmail({
      contactId: contactId,
      email: emailTwo,
    });

    const contactAfterSecondEmailSetPrimary = await contactService.getContact(
      contactId,
    );

    expect(contactAfterSecondEmailSetPrimary?.emails.length).toBe(3);

    emailOneEntry = contactAfterSecondEmailSetPrimary?.emails.find(
      (email) => email.email === emailOne,
    );
    emailTwoEntry = contactAfterSecondEmailSetPrimary?.emails.find(
      (email) => email.email === emailTwo,
    );

    expect(emailOneEntry).toBeDefined();
    expect(emailTwoEntry).toBeDefined();

    expect(emailOneEntry?.id).not.toBeNull();
    expect(emailOneEntry?.primary, 'The first email is still primary').toBe(
      false,
    );

    expect(emailTwoEntry?.id).not.toBeNull();
    expect(emailTwoEntry?.primary, 'The first email is still primary').toBe(
      true,
    );
  });

  it('checks successful response from better contact', async () => {
    const { organizationId } = await VitestHelper.createOrganizationForTest(
      organizationsRepository,
    );
    const contact_social_url =
      'https://www.linkedin.com/in/vitest-' + crypto.randomUUID();

    const contactId = await VitestHelper.createContactForOrganizationForTest(
      contactService,
      organizationId,
      {
        socialUrl: contact_social_url,
      },
    );

    await contactService.getContact(contactId.contactId);

    const { contact_FindWorkEmail } = await contactService.findEmail({
      contactId: contactId.contactId,
      organizationId: organizationId,
    });

    expect(
      contact_FindWorkEmail.accepted,
      'The findEmail mutation returned error',
    ).toBe(true);
  });

  it('does CRUD ops for phone number', async () => {
    const { organizationId } = await VitestHelper.createOrganizationForTest(
      organizationsRepository,
    );
    const contact_social_url =
      'https://www.linkedin.com/in/vitest-' + crypto.randomUUID();

    const contactId = await VitestHelper.createContactForOrganizationForTest(
      contactService,
      organizationId,
      {
        socialUrl: contact_social_url,
      },
    );

    //also here the phonenumber is not used anymore in our app until we reimplement it in BE

    // const contactBeforeAddingPhoneNumber = await contactService.getContact(
    //   contact_CreateForOrganization.id,
    // );

    // expect(
    //   contactBeforeAddingPhoneNumber?.phoneNumbers.length,
    //   'The contact has phone number before adding',
    // ).toBe(0);

    const expectedCreateFirstPhoneNumber = Array.from(
      crypto.getRandomValues(new Uint32Array(1)),
    )[0]
      .toString()
      .slice(-8)
      .padStart(8, '0');
    const expectedFirstPhoneNumberLabel = PhoneNumberLabel.Mobile;

    const firstAddedNumber = await contactService.addPhoneNumber({
      contactId: contactId.contactId,
      input: {
        phoneNumber: expectedCreateFirstPhoneNumber,
        label: expectedFirstPhoneNumberLabel,
        primary: false,
        countryCodeA2: 'US',
      },
    });

    expect(firstAddedNumber.phoneNumberMergeToContact.id).not.toBeNull();
    expect(firstAddedNumber.phoneNumberMergeToContact.rawPhoneNumber).toBe(
      expectedCreateFirstPhoneNumber,
    );

    // const contactAfterAddingFirstPhoneNumber = await contactService.getContact(
    //   contact_CreateForOrganization.id,
    // );

    // expect(
    //   contactAfterAddingFirstPhoneNumber.contact?.phoneNumbers.length,
    //   "The contact doesn't have exactly 1 phone number",
    // ).toBe(1);
    // expect(contactAfterAddingFirstPhoneNumber.contact?.phoneNumbers[0].id).toBe(
    //   firstAddedNumber.phoneNumberMergeToContact.id,
    // );
    // expect(
    //   contactAfterAddingFirstPhoneNumber.contact?.phoneNumbers[0].label,
    // ).toBe(expectedFirstPhoneNumberLabel);
    // expect(
    //   contactAfterAddingFirstPhoneNumber.contact?.phoneNumbers[0]
    //     .rawPhoneNumber,
    // ).toBe(expectedCreateFirstPhoneNumber);
    // expect(
    //   contactAfterAddingFirstPhoneNumber.contact?.phoneNumbers[0].e164,
    // ).toBeNull();
    // expect(
    //   contactAfterAddingFirstPhoneNumber.contact?.phoneNumbers[0].primary,
    // ).toBe(false);

    // const expectedUpdateFirstPhoneNumber = Array.from(
    //   crypto.getRandomValues(new Uint32Array(1)),
    // )[0]
    //   .toString()
    //   .slice(-8)
    //   .padStart(8, '0');
    // const { phoneNumber_Update } = await contactService.updatePhoneNumber({
    //   input: {
    //     id: firstAddedNumber.phoneNumberMergeToContact.id,
    //     phoneNumber: expectedUpdateFirstPhoneNumber,
    //     countryCodeA2: 'RO',
    //   },

    // expect(phoneNumber_Update.id).not.toBeNull();

    // const contactAfterUpdatingFirstPhoneNumber =
    //   await contactService.getContact(contact_CreateForOrganization.id);

    // expect(
    //   contactAfterUpdatingFirstPhoneNumber.contact?.phoneNumbers.length,
    //   "The contact doesn't have exactly 1 phone number",
    // ).toBe(1);
    // expect(
    //   contactAfterUpdatingFirstPhoneNumber.contact?.phoneNumbers[0].id,
    // ).toBe(firstAddedNumber.phoneNumberMergeToContact.id);
    // expect(
    //   contactAfterUpdatingFirstPhoneNumber.contact?.phoneNumbers[0].label,
    // ).toBe(expectedFirstPhoneNumberLabel);
    // expect(
    //   contactAfterUpdatingFirstPhoneNumber.contact?.phoneNumbers[0]
    //     .rawPhoneNumber,
    // ).toBe(expectedUpdateFirstPhoneNumber);
    // expect(
    //   contactAfterUpdatingFirstPhoneNumber.contact?.phoneNumbers[0].e164,
    // ).toBeNull();
    // expect(
    //   contactAfterUpdatingFirstPhoneNumber.contact?.phoneNumbers[0].primary,
    // ).toBe(false);

    // const expectedCreateSecondPhoneNumber = Array.from(
    //   crypto.getRandomValues(new Uint32Array(1)),
    // )[0]
    //   .toString()
    //   .slice(-8)
    //   .padStart(8, '0');
    // const expectedSecondPhoneNumberLabel = PhoneNumberLabel.Mobile;

    // const secondAddedNumber = await contactService.addPhoneNumber({
    //   contactId: contact_CreateForOrganization.id,
    //   input: {
    //     phoneNumber: expectedCreateSecondPhoneNumber,
    //     label: expectedSecondPhoneNumberLabel,
    //     primary: false,
    //     countryCodeA2: 'US',
    //   },
    // });

    // expect(secondAddedNumber.phoneNumberMergeToContact.id).not.toBeNull();
    // expect(secondAddedNumber.phoneNumberMergeToContact.rawPhoneNumber).toBe(
    //   expectedCreateSecondPhoneNumber,
    // );

    // const contactAfterAddingSecondPhoneNumber = await contactService.getContact(
    //   contact_CreateForOrganization.id,
    // );

    // expect(
    //   contactAfterAddingSecondPhoneNumber.contact?.phoneNumbers.length,
    //   "The contact doesn't have exactly 2 phone numbers",
    // ).toBe(2);

    // await contactService.removePhoneNumber({
    //   contactId: contact_CreateForOrganization.id,
    //   id: firstAddedNumber.phoneNumberMergeToContact.id,
    // });

    // const contactAfterRemovingFirstPhoneNumber =
    //   await contactService.getContact(contact_CreateForOrganization.id);

    // expect(
    //   contactAfterRemovingFirstPhoneNumber.contact?.phoneNumbers.length,
    //   "The contact doesn't have exactly 1 phone numbers",
    // ).toBe(1);

    // expect(
    //   contactAfterRemovingFirstPhoneNumber.contact?.phoneNumbers[0]
    //     .rawPhoneNumber,
    // ).toBe(expectedCreateSecondPhoneNumber);
  });

  // it('creates and updates social for contact', async () => {
  // const { id } = await VitestHelper.createOrganizationForTest(
  //   organizationsRepository,
  // );

  // const { contact_CreateForOrganization } =
  //   await contactService.createContactForOrganization({
  //     organizationId: id,
  //     input: {},
  //   });
  // const contactCreated = await contactService.getContact(
  //   contact_CreateForOrganization.id,
  // );

  // expect(contactCreated.contact?.socials.length).toBe(0);
  // expect(contactCreated.contact?.organizations.content.length).toBe(1);
  // expect(contactCreated.contact?.organizations.content[0].name).toBe(name);

  // const contact_added_social_url =
  //   'https://www.linkedin.com/in/Vitest_' + crypto.randomUUID();

  // await contactService.addSocial({
  //   contactId: contactCreated.contact!.metadata.id,
  //   input: { url: contact_added_social_url },
  // });

  // const contactAddedSocial = await contactService.getContact(
  //   contact_CreateForOrganization.id,
  // );

  // expect(contactAddedSocial.contact?.socials.length).toBe(1);
  // expect(contactAddedSocial.contact?.socials[0].url).toBe(
  //   contact_added_social_url,
  // );

  // const contact_updated_social_url =
  //   'https://www.linkedin.com/in/Vitest_' + crypto.randomUUID();

  // await contactService.updateSocial({
  //   input: {
  //     id: contactAddedSocial.contact!.socials[0].id,
  //     url: contact_updated_social_url,
  //   },
  // });

  // const contactUpdatedSocial = await contactService.getContact(
  //   contact_CreateForOrganization.id,
  // );

  // expect(contactUpdatedSocial.contact?.socials.length).toBe(1);
  // expect(contactUpdatedSocial.contact?.socials[0].url).toBe(
  //   contact_updated_social_url,
  // );
  // });

  // it('archives contact', async () => {
  //   const contact_social_url =
  //     'https://www.linkedin.com/in/Vitest_' + crypto.randomUUID();

  //   const { contact_Create } = await contactService.createContact({
  //     contactInput: { socialUrl: contact_social_url },
  //   });
  //   const contacts_before_archiving = await contactService.getContact({
  //     pagination: { limit: 100, page: 0 },
  //   });

  //   let hasId = (id: string): boolean => {
  //     return (
  //       contacts_before_archiving?.contacts?.content?.some(
  //         (cont) => cont.id === id,
  //       ) ?? false
  //     );
  //   };
  //   let contactExistsInDashboard = hasId(contact_Create);

  //   expect(contactExistsInDashboard).toBe(true);

  //   await contactService.archiveContact({
  //     contactId: contact_Create,
  //   });

  //   const retrieved_contacts = await contactService.getContacts({
  //     pagination: { limit: 100, page: 0 },
  //   });

  //   hasId = (id: string): boolean => {
  //     return (
  //       retrieved_contacts?.contacts?.content?.some((cont) => cont.id === id) ??
  //       false
  //     );
  //   };
  //   contactExistsInDashboard = hasId(contact_Create);

  //   expect(contactExistsInDashboard).toBe(false);
  // });

  it('creates and updates tags for contact', async () => {
    const { organizationId, organizationName } =
      await VitestHelper.createOrganizationForTest(organizationsRepository);

    const contact_social_url =
      'https://www.linkedin.com/in/vitest-' + crypto.randomUUID();

    const contactId = await VitestHelper.createContactForOrganizationForTest(
      contactService,
      organizationId,
      {
        socialUrl: contact_social_url,
      },
    );
    const contactCreated = await contactService.getContact(contactId.contactId);

    expect(contactCreated?.tags).toEqual([]);
    expect(contactCreated.primaryOrganizationName).toBe(organizationName);

    const firstTagName = 'Vitest_' + crypto.randomUUID();

    await contactService.addTagsToContact({
      input: {
        contactId: contactCreated.id,
        tag: { name: firstTagName, entityType: EntityType.Contact },
      },
    });

    const contactAddedFirstTag = await contactService.getContact(
      contactId.contactId,
    );

    expect(contactAddedFirstTag.tags?.length).toBe(1);
    expect(contactAddedFirstTag.tags?.[0]?.name).toBe(firstTagName);

    const secondTagName = 'Vitest_' + crypto.randomUUID();

    await contactService.addTagsToContact({
      input: {
        contactId: contactCreated.id,
        tag: { name: secondTagName },
      },
    });

    const contactAddedSecondTag = await contactService.getContact(
      contactId.contactId,
    );

    expect(contactAddedSecondTag?.tags?.length).toBe(2);

    await contactService.removeTagsFromContact({
      input: {
        contactId: contactCreated.id,
        tag: { name: firstTagName },
      },
    });

    const contactRemovedFirstTag = await contactService.getContact(
      contactId.contactId,
    );

    expect(contactRemovedFirstTag?.tags?.length).toBe(1);
    expect(contactRemovedFirstTag?.tags?.[0]?.name).toBe(secondTagName);
  });

  it('creates a contact by LinkedinIn, with no organization', async () => {
    const linkedIn_urls = [
      'https://www.linkedin.com/in/Vitest-' + crypto.randomUUID(),
    ];

    // const contactId =
    await VitestHelper.createContactBulkByLinkedInForTest(
      contactService,
      linkedIn_urls,
    );

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
        linkedIn_urls.includes(contact.linkedInUrl),
    );

    expect(createContacts).toHaveLength(linkedIn_urls.length);
  });

  it('creates multiple contacts by LinkedinIn, with no organization', async () => {
    const linkedIn_urls = [
      'https://www.linkedin.com/in/Vitest-' + crypto.randomUUID(),
      'https://www.linkedin.com/in/Vitest-' + crypto.randomUUID(),
    ];

    // const contactId =
    await VitestHelper.createContactBulkByLinkedInForTest(
      contactService,
      linkedIn_urls,
    );

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
        linkedIn_urls.includes(contact.linkedInUrl),
    );

    expect(createContacts).toHaveLength(linkedIn_urls.length);
  });

  it('creates a contact by Email, with no organization', async () => {
    const emails = [generateEmail()];

    await VitestHelper.createContactBulkByEmailForTest(contactService, emails);

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

    expect(createContacts).toHaveLength(emails.length);
  });

  it('creates multiple contacts by Email, with no organization', async () => {
    const emails = [generateEmail(), generateEmail()];

    await VitestHelper.createContactBulkByEmailForTest(contactService, emails);

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

    expect(createContacts).toHaveLength(emails.length);
  });
});
