import { randomUUID } from 'crypto';
import { Page, expect } from '@playwright/test';

import {
  createTinyUUID,
  writeTextInLocator,
  createRequestPromise,
  createResponsePromise,
  clickLocatorThatIsVisible,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class OrganizationPeoplePage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgPeopleAddSomeone = 'button[data-test="org-people-add-someone"]';
  private orgPeopleAddByName = 'button[data-test="org-people-add-by-name"]';
  private orgPeopleAddContact = 'button[data-test="org-people-add-contact"]';
  private orgPeopleCollapse = 'button[data-test="org-people-collapse"]';
  private orgPeopleContactName = 'input[data-test="org-people-contact-name"]';
  private orgPeopleContactTitle = 'input[data-test="org-people-contact-title"]';
  private orgPeopleContactEmail = 'p[data-test="add-work-email"]';
  private orgPeopleLinkedInUrl = 'span[data-test="org-people-linkedin"]';
  private orgPeopleLinkedInInput = 'input[data-test="linkedin-url-input"]';
  private orgPeopleConfirmLinkedInUrl = 'button[data-test="add-linkedin-url"]';
  private orgPeopleorgAboutTags = 'div[data-test="org-about-tags"]';

  async addContact(contactCreation: string) {
    await clickLocatorsThatAreVisible(
      this.page,
      contactCreation,
      this.orgPeopleAddByName,
    );

    const createContactResponsePromise = createResponsePromise(
      this.page,
      'contact_CreateForOrganization?.id',
      undefined,
    );

    const contactResponsePromise = createResponsePromise(
      this.page,
      'ui_contacts',
      undefined,
    );
    const organizationResponsePromise = createResponsePromise(
      this.page,
      'ui_organizations',
      undefined,
    );

    await clickLocatorsThatAreVisible(this.page, this.orgPeopleAddContact);

    await Promise.all([
      createContactResponsePromise,
      contactResponsePromise,
      organizationResponsePromise,
    ]);
  }

  async addNameToContact() {
    await this.page.waitForTimeout(3000);
    await clickLocatorThatIsVisible(this.page, this.orgPeopleCollapse);

    const orgPeopleContactNameInput = this.page.locator(
      this.orgPeopleContactName,
    );

    const contactName = createTinyUUID();

    const requestPromise = createRequestPromise(this.page, 'name', contactName);

    const responsePromise = createResponsePromise(
      this.page,
      'contact_Update.id',
      undefined,
    );

    await orgPeopleContactNameInput.pressSequentially(contactName, {
      delay: 100,
    });
    await orgPeopleContactNameInput.press('Tab');

    const [_, response] = await Promise.all([requestPromise, responsePromise]);

    await expect(orgPeopleContactNameInput).toHaveValue(contactName);

    const responseBody = await response.json();
    const contactId = responseBody.data?.contact_Update?.id;

    return { contactName, contactId };
  }

  async addTitleToContact() {
    const orgPeopleContactTitleInput = this.page.locator(
      this.orgPeopleContactTitle,
    );

    const requestPromise = createRequestPromise(this.page, 'jobTitle', 'CTO');

    const responsePromise = createResponsePromise(
      this.page,
      'jobRole_Update?.id',
      undefined,
    );

    await orgPeopleContactTitleInput.pressSequentially('CTO', { delay: 500 });
    await orgPeopleContactTitleInput.press('Tab');
    await Promise.all([requestPromise, responsePromise]);
    await expect(orgPeopleContactTitleInput).toHaveValue('CTO');
  }

  private async addDetailsToCustomer() {
    await clickLocatorThatIsVisible(this.page, this.orgPeopleContactTitle);

    let page = await writeTextInLocator(
      this.page,
      this.orgPeopleContactEmail,
      'contact@org.com',
    );

    await page.keyboard.press('Enter');

    const contactLinkedInProfile = 'www.linkedin.com/in/' + randomUUID();

    await clickLocatorThatIsVisible(page, this.orgPeopleLinkedInUrl);
    page = await writeTextInLocator(
      page,
      this.orgPeopleLinkedInInput,
      contactLinkedInProfile,
    );
    await clickLocatorThatIsVisible(page, this.orgPeopleConfirmLinkedInUrl);

    page = await writeTextInLocator(
      page,
      this.orgPeopleorgAboutTags,
      'testPersonas',
    );

    await page.keyboard.press('Enter');
  }

  async createContactFromEmpty() {
    await this.addContact(this.orgPeopleAddSomeone);

    const { contactName, contactId } = await this.addNameToContact();

    await this.addTitleToContact();
    await this.addDetailsToCustomer();

    return { contactName, contactId };
  }
}
