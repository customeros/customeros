import { Page, expect } from '@playwright/test';

import {
  writeTextInLocator,
  createRequestPromise,
  createResponsePromise,
  ensureLocatorIsVisible,
  clickLocatorThatIsVisible,
  clickLocatorsThatAreVisible,
  doubleClickLocatorThatIsVisible,
} from '../../helper';

export class OrganizationPeoplePage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgPeopleAddSomeone = 'button[data-test="org-people-add-someone"]';
  private orgPeopleContactName = 'input[data-test="org-people-contact-name"]';
  private orgPeopleContactTitle = 'input[data-test="org-people-contact-title"]';
  private orgPeopleContactJobRoles =
    'div[data-test="org-people-contact-job-roles"]';
  private jobRoleInfluencer = 'div[role="option"]:has-text("Influencer")';
  private orgPeopleContactClose =
    'button[data-test="org-people-contact-close"]';
  private orgPeopleContactDelete =
    'button[data-test="org-people-contact-delete"]';
  private orgPeopleContactEmail = 'input[data-test="org-people-contact-email"]';
  private orgPeopleContactPhoneNumber =
    'input[data-test="org-people-contact-phone-number"]';
  private orgPeopleContactPersonas =
    'div[data-test="org-people-contact-personas"]';
  private orgPeopleContactSocialLink =
    'input[data-test="org-people-contact-social-link"]';
  private orgPeopleContactTimezone =
    'div[data-test="org-people-contact-timezone"]';

  async addContactEmpty() {
    await clickLocatorsThatAreVisible(this.page, this.orgPeopleAddSomeone);

    const response = await this.page.waitForResponse(
      (response) =>
        response.url().includes('customer-os-api') && //createContact
        response
          .json()
          .then(
            (body) =>
              body.data &&
              body.data?.contact_CreateForOrganization?.id !== undefined,
          )
          .catch(() => false),
    );

    await response.json();
  }

  async addNameToContact() {
    const orgPeopleContactNameInput = this.page.locator(
      this.orgPeopleContactName,
    );

    const requestPromise = createRequestPromise(this.page, 'name', 'John Doe');

    const responsePromise = createResponsePromise(
      this.page,
      'contact_Update?.id',
      undefined,
    );

    await orgPeopleContactNameInput.pressSequentially('John Doe', {
      delay: 500,
    });
    await Promise.all([requestPromise, responsePromise]);
    await expect(orgPeopleContactNameInput).toHaveValue('John Doe');
  }

  async addTitleToContact() {
    const orgPeopleContactTitleInput = this.page.locator(
      this.orgPeopleContactTitle,
    );

    const requestPromise = createRequestPromise(this.page, 'jobTitle', 'CTO');

    const responsePromise = createResponsePromise(
      this.page,
      'jobRole_Create?.id',
      undefined,
    );

    await orgPeopleContactTitleInput.pressSequentially('CTO', { delay: 500 });
    await Promise.all([requestPromise, responsePromise]);
    await expect(orgPeopleContactTitleInput).toHaveValue('CTO');
  }

  async addJobRolesToContact() {
    const orgPeopleContactJobRolesInput = this.page.locator(
      this.orgPeopleContactJobRoles,
    );

    await orgPeopleContactJobRolesInput.click();

    await this.page.waitForSelector('[role="listbox"]', { state: 'visible' });

    const influencerOption = this.page.locator(this.jobRoleInfluencer);

    const requestPromise = createRequestPromise(
      this.page,
      'description',
      'Influencer',
    );

    const responsePromise = createResponsePromise(
      this.page,
      'jobRole_Create?.id',
      undefined,
    );

    await influencerOption.click();
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  private async addDetailsToCustomer() {
    await clickLocatorThatIsVisible(this.page, this.orgPeopleContactTitle);

    let page = await writeTextInLocator(
      this.page,
      this.orgPeopleContactEmail,
      'contact@org.com',
    );

    page = await writeTextInLocator(
      page,
      this.orgPeopleContactPhoneNumber,
      '0741111111',
    );

    page = await writeTextInLocator(
      page,
      this.orgPeopleContactPersonas,
      'testPersonas',
    );

    await page.keyboard.press('Enter');
    await clickLocatorsThatAreVisible(page, this.orgPeopleContactPersonas);
    page = await writeTextInLocator(
      this.page,
      this.orgPeopleContactSocialLink,
      'www.linkedin.com/in/test',
    );

    const requestPromise = createRequestPromise(
      this.page,
      'url',
      'www.linkedin.com/in/test',
    );

    const responsePromise = createResponsePromise(
      this.page,
      'contact_AddSocial?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(page, this.orgPeopleContactTimezone);
    await Promise.all([requestPromise, responsePromise]);

    await doubleClickLocatorThatIsVisible(page, this.orgPeopleContactTimezone);

    const locator = await ensureLocatorIsVisible(
      page,
      this.orgPeopleContactTimezone,
    );

    await locator.pressSequentially('new salem');
    await page.keyboard.press('Enter');
  }

  async createContactFromEmpty() {
    await this.addContactEmpty();
    await this.addNameToContact();
    await this.addTitleToContact();
    await this.addJobRolesToContact();
    await this.addDetailsToCustomer();
  }
}
