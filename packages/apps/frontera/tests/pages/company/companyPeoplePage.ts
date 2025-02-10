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

export class CompanyPeoplePage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgPeopleAddSomeone = 'button[data-test="org-people-add-someone"]';
  private orgPeopleAddByName = 'button[data-test="org-people-add-by-name"]';
  private confirmContactCreation =
    'button[data-test="confirm-contact-creation"]';
  private orgPeopleAddContact =
    'button[data-test="org-people-add-new-contact"]';
  private orgPeopleCollapse = 'button[data-test="org-people-collapse"]';
  private orgPeopleNameInput = 'input[data-test="org-people-name-input"]';
  // private orgPeopleContactName = 'input[data-test="org-people-contact-name"]';

  private orgPeopleContactTitle = 'input[data-test="org-people-contact-title"]';
  private orgPeopleContactEmail = 'p[data-test="add-work-email"]';
  private orgPeopleLinkedInUrl = 'span[data-test="org-people-linkedin"]';
  private orgPeopleLinkedInInput = 'input[data-test="linkedin-url-input"]';
  private orgPeopleConfirmLinkedInUrl = 'button[data-test="add-linkedin-url"]';
  private orgPeopleorgAboutTags = 'div[data-test="org-about-tags"]';
  private orgPeopleContactTags = 'div[data-test="contact-tags"]';

  async addContact(contactCreation: string) {
    await clickLocatorsThatAreVisible(
      this.page,
      contactCreation,
      this.orgPeopleAddByName,
    );

    const contactName = createTinyUUID();

    await this.page.fill(this.orgPeopleNameInput, contactName);
    await clickLocatorThatIsVisible(this.page, this.confirmContactCreation);

    // const contactName = await this.addNameToContact();
    //
    // await clickLocatorsThatAreVisible(this.page, this.orgPeopleAddContact);

    return contactName;
  }

  async addNameToContact() {
    await this.page.waitForTimeout(3000);

    const orgPeopleContactNameInput = await clickLocatorThatIsVisible(
      this.page,
      this.orgPeopleNameInput,
    );

    const contactName = createTinyUUID();

    await orgPeopleContactNameInput.pressSequentially(contactName, {
      delay: 100,
    });

    return contactName;
  }

  async addTitleToContact() {
    const orgPeopleContactTitleInput = this.page.locator(
      this.orgPeopleContactTitle,
    );

    const requestPromise = createRequestPromise(
      this.page,
      'input?.jobTitle',
      'CTO',
    );

    const responsePromise = createResponsePromise(
      this.page,
      'jobRole_Save',
      undefined,
    );

    await this.page.waitForTimeout(500);
    await orgPeopleContactTitleInput.pressSequentially('CTO', { delay: 500 });
    await this.page.keyboard.press('Tab');
    await Promise.all([requestPromise, responsePromise]);
    await expect(orgPeopleContactTitleInput).toHaveValue('CTO');
  }

  private async addDetailsToCustomer() {
    // await clickLocatorThatIsVisible(this.page, this.orgPeopleContactTitle);

    await this.addTitleToContact();

    const emailUsername = createTinyUUID();
    let page = await writeTextInLocator(
      this.page,
      this.orgPeopleContactEmail,
      emailUsername + '@org.com',
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
      this.orgPeopleContactTags,
      'testPersonas',
    );

    await page.keyboard.press('Enter');
  }

  async createContactFromEmpty() {
    const contactName = await this.addContact(this.orgPeopleAddSomeone);

    await clickLocatorThatIsVisible(this.page, this.orgPeopleCollapse);
    await this.addDetailsToCustomer();

    return contactName;
  }
}
