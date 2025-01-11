import { it, expect, describe } from 'vitest';
import { VitestHelper } from '@store/vitest-helper.ts';

import { EmailLabel, EntityType, PhoneNumberLabel } from '@graphql/types';

import { ContactService } from '../../Contacts/__service__/Contacts.service';
import { OrganizationsService } from '../../Organizations/__service__/Organizations.service';

const organizationsService = OrganizationsService.getInstance();
const contactService = ContactService.getInstance();

describe('ContactsService - Integration Tests', () => {
  it('create contact', async () => {
    const contact_social_url =
      'https://www.linkedin.com/in/IT_' + crypto.randomUUID();

    const { contact_Create } = await contactService.createContact({
      contactInput: { socialUrl: contact_social_url },
    });
    const contact = await contactService.getContact(contact_Create);

    expect(contact?.id).toBe(contact_Create);
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

  it('creates contact for organization', async () => {
    const { id, name } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );
    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: id,
        input: {},
      });
    const contact = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect(contact?.primaryOrganizationName).toBe(name);
  });

  it('update contact', async () => {
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );
    const contact_social_url = 'IT_' + crypto.randomUUID();
    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: id,
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

    const contact = await contactService.getContact(
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
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const contact_first_social_url = 'IT_' + crypto.randomUUID();

    const firstContact = await contactService.createContactForOrganization({
      organizationId: id,
      input: { socialUrl: contact_first_social_url },
    });
    const firstContactId = firstContact.contact_CreateForOrganization.id;

    const contact_second_social_url = 'IT_' + crypto.randomUUID();

    const secondContact = await contactService.createContactForOrganization({
      organizationId: id,
      input: { socialUrl: contact_second_social_url },
    });

    const secondContactId: string =
      secondContact.contact_CreateForOrganization.id;
    const contact_third_social_url = 'IT_' + crypto.randomUUID();

    const thirdContact = await contactService.createContactForOrganization({
      organizationId: id,
      input: { socialUrl: contact_third_social_url },
    });

    const thirdContactId: string =
      thirdContact.contact_CreateForOrganization.id;
    const { ui_contacts } = await contactService.getContactsByIds({
      ids: [firstContactId, secondContactId, thirdContactId],
    });
    const contactIds = ui_contacts?.map((contact) => contact.id);

    expect(contactIds).toBeDefined();
    expect(contactIds).toContain(firstContact.contact_CreateForOrganization.id);
    expect(contactIds).toContain(
      secondContact.contact_CreateForOrganization.id,
    );
    expect(contactIds).toContain(thirdContact.contact_CreateForOrganization.id);
  });

  it('adds job roles to contact', async () => {
    const { id, name } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );
    const contact_social_url = 'IT_' + crypto.randomUUID();

    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: id,
        input: { socialUrl: contact_social_url },
      });

    const contactBeforeFirstJobRole = await contactService.getContact(
      contact_CreateForOrganization.id,
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
    expect(contactBeforeFirstJobRole.primaryOrganizationName).toBe(name);

    const jobRoleCreateOneDescription = 'IT_' + crypto.randomUUID();
    const jobRoleCreateOneTitle = 'IT_' + crypto.randomUUID();
    const jobRoleCreateOneStartedAt = new Date().toISOString();

    const { jobRole_Create } = await contactService.addJobRole({
      contactId: contactBeforeFirstJobRole.id,
      input: {
        description: jobRoleCreateOneDescription,
        jobTitle: jobRoleCreateOneTitle,
        organizationId: id,
        startedAt: jobRoleCreateOneStartedAt,
      },
    });

    //this needs to be reimplemented because now we are working with only primary job role

    // const contactAfterFirstJobRole = await contactService.getContact(
    //   contact_CreateForOrganization.id,
    // );

    // expect(contactAfterFirstJobRole.contact?.jobRoles.length).toBe(2);
    // expect(contactAfterFirstJobRole.contact?.jobRoles[1].jobTitle).toBe(
    //   jobRoleCreateOneTitle,
    // );
    // expect(contactAfterFirstJobRole.contact?.jobRoles[1].description).toBe(
    //   jobRoleCreateOneDescription,
    // );
    // expect(contactAfterFirstJobRole.contact?.jobRoles[1].primary).toBe(false);
    // expect(contactAfterFirstJobRole.contact?.jobRoles[1].company).toBeNull();
    // expect(
    //   new Date(
    //     contactAfterFirstJobRole.contact?.jobRoles[1].startedAt,
    //   ).getTime(),
    // ).toBe(new Date(jobRoleCreateOneStartedAt).getTime());
    // expect(contactAfterFirstJobRole.contact?.jobRoles[0].endedAt).toBeNull();

    const jobRoleUpdateDescription = 'IT_' + crypto.randomUUID();
    const jobRoleUpdateTitle = 'IT_' + crypto.randomUUID();
    const jobRoleUpdateStartedAt = new Date().toISOString();
    const jobRoleUpdateCompany = new Date().toISOString();

    await new Promise((resolve) => setTimeout(resolve, 1000)); // 1 second delay to have different endedAt

    const jobRoleUpdateEndedAt = new Date().toISOString();

    await contactService.updateJobRole({
      contactId: contactBeforeFirstJobRole.id,
      input: {
        id: jobRole_Create.id,
        description: jobRoleUpdateDescription,
        jobTitle: jobRoleUpdateTitle,
        organizationId: id,
        startedAt: jobRoleUpdateStartedAt,
        company: jobRoleUpdateCompany,
        endedAt: jobRoleUpdateEndedAt,
        primary: true,
      },
    });

    //Same here needs to be reimplemented

    // const contactAfterUpdateJobRole = await contactService.getContact(
    //   contact_CreateForOrganization.id,
    // );

    // expect(contactAfterUpdateJobRole.contact?.jobRoles.length).toBe(2);
    // expect(contactAfterUpdateJobRole.contact?.jobRoles[1].jobTitle).toBe(
    //   jobRoleUpdateTitle,
    // );
    // expect(contactAfterUpdateJobRole.contact?.jobRoles[1].description).toBe(
    //   jobRoleUpdateDescription,
    // );
    // expect(contactAfterUpdateJobRole.contact?.jobRoles[1].primary).toBe(true);
    // expect(contactAfterUpdateJobRole.contact?.jobRoles[1].company).toBe(
    //   jobRoleUpdateCompany,
    // );
    // expect(
    //   new Date(
    //     contactAfterUpdateJobRole.contact?.jobRoles[1].startedAt,
    //   ).getTime(),
    // ).toBe(new Date(jobRoleUpdateStartedAt).getTime());
    // expect(
    //   new Date(
    //     contactAfterUpdateJobRole.contact?.jobRoles[1].endedAt,
    //   ).getTime(),
    // ).toBe(new Date(jobRoleUpdateEndedAt).getTime());
  });

  it('links contact to organization', async () => {
    const { id, name } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );
    const { contact_Create } = await contactService.createContact({
      contactInput: {},
    });

    await organizationsService.getOrganization(id);

    const contactBeforeLink = await contactService.getContact(contact_Create);

    expect(contactBeforeLink.primaryOrganizationName).toBeNull;

    await contactService.linkOrganization({
      input: {
        contactId: contact_Create,
        organizationId: id,
      },
    });
    expect(
      (await contactService.getContact(contact_Create)).primaryOrganizationId,
    ).toBe(id);
    expect(
      (await contactService.getContact(contact_Create)).primaryOrganizationName,
    ).toBe(name);
  });

  it('updates contact email', async () => {
    const { contact_Create } = await contactService.createContact({
      contactInput: {},
    });
    const emailOne = 'IT_' + crypto.randomUUID() + '@example.com';

    await contactService.updateContactEmail({
      contactId: contact_Create,
      previousEmail: '',
      input: {
        email: emailOne,
        appSource: 'IT_test',
        label: EmailLabel.Personal,
      },
    });

    const contactAfterEmailUpdate = await contactService.getContact(
      contact_Create,
    );

    expect(contactAfterEmailUpdate.emails.length).toBe(1);
    expect(contactAfterEmailUpdate.emails[0].id).not.toBeNull();
    expect(contactAfterEmailUpdate.emails[0].email).toBe(emailOne);
    expect(contactAfterEmailUpdate.emails[0].primary).toBe(true);
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

    const emailTwo = 'IT_' + crypto.randomUUID() + '@example.com';

    await contactService.updateContactEmail({
      contactId: contact_Create,
      previousEmail: '',
      input: {
        email: emailTwo,
        appSource: 'IT_test',
        label: EmailLabel.Work,
      },
    });

    const contactAfterSecondEmailUpdate = await contactService.getContact(
      contact_Create,
    );

    expect(contactAfterSecondEmailUpdate?.emails.length).toBe(2);

    let emailOneEntry = contactAfterSecondEmailUpdate?.emails.find(
      (email) => email.email === emailOne,
    );
    let emailTwoEntry = contactAfterSecondEmailUpdate?.emails.find(
      (email) => email.email === emailTwo,
    );

    expect(emailOneEntry).toBeDefined();
    expect(emailTwoEntry).toBeDefined();

    expect(emailOneEntry?.id).not.toBeNull();
    expect(emailOneEntry?.primary).toBe(true);

    expect(emailTwoEntry?.id).not.toBeNull();
    expect(emailTwoEntry?.primary).toBe(false);

    await contactService.setPrimaryEmail({
      contactId: contact_Create,
      email: emailTwo,
    });

    const contactAfterSecondEmailSetPrimary = await contactService.getContact(
      contact_Create,
    );

    expect(contactAfterSecondEmailSetPrimary?.emails.length).toBe(2);

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
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );
    const contact_social_url = 'IT_' + crypto.randomUUID();

    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: id,
        input: { socialUrl: contact_social_url },
      });

    await contactService.getContact(contact_CreateForOrganization.id);

    const { contact_FindWorkEmail } = await contactService.findEmail({
      contactId: contact_CreateForOrganization.id,
      organizationId: id,
    });

    expect(
      contact_FindWorkEmail.accepted,
      'The findEmail mutation returned error',
    ).toBe(true);
  });

  it('does CRUD ops for phone number', async () => {
    const { id } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );
    const contact_social_url = 'IT_' + crypto.randomUUID();

    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: id,
        input: { socialUrl: contact_social_url },
      });

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
      contactId: contact_CreateForOrganization.id,
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
  //   organizationsService,
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
  //   'https://www.linkedin.com/in/IT_' + crypto.randomUUID();

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
  //   'https://www.linkedin.com/in/IT_' + crypto.randomUUID();

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
  //     'https://www.linkedin.com/in/IT_' + crypto.randomUUID();

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
    const { id, name } = await VitestHelper.createOrganizationForTest(
      organizationsService,
    );

    const { contact_CreateForOrganization } =
      await contactService.createContactForOrganization({
        organizationId: id,
        input: {},
      });
    const contactCreated = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect(contactCreated?.tags).toEqual([]);
    expect(contactCreated.primaryOrganizationName).toBe(name);

    const firstTagName = crypto.randomUUID();

    await contactService.addTagsToContact({
      input: {
        contactId: contactCreated.id,
        tag: { name: firstTagName, entityType: EntityType.Contact },
      },
    });

    const contactAddedFirstTag = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect(contactAddedFirstTag.tags?.length).toBe(1);
    expect(contactAddedFirstTag.tags?.[0]?.name).toBe(firstTagName);

    const secondTagName = crypto.randomUUID();

    await contactService.addTagsToContact({
      input: {
        contactId: contactCreated.id,
        tag: { name: secondTagName },
      },
    });

    const contactAddedSecondTag = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect(contactAddedSecondTag?.tags?.length).toBe(2);

    await contactService.removeTagsFromContact({
      input: {
        contactId: contactCreated.id,
        tag: { name: firstTagName },
      },
    });

    const contactRemovedFirstTag = await contactService.getContact(
      contact_CreateForOrganization.id,
    );

    expect(contactRemovedFirstTag?.tags?.length).toBe(1);
    expect(contactRemovedFirstTag?.tags?.[0]?.name).toBe(secondTagName);
  });
});
