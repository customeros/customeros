import { Page, expect } from '@playwright/test';

import {
  writeTextInLocator,
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
    // await this.page.waitForResponse('**/customer-os-api');
  }

  async addNameToContact() {
    const orgPeopleContactNameInput = this.page.locator(
      this.orgPeopleContactName,
    );

    const requestPromise = this.page.waitForRequest((request) => {
      if (
        request.method() === 'POST' &&
        request.url().includes('customer-os-api')
      ) {
        const postData = request.postData();

        if (postData) {
          try {
            const parsedData = JSON.parse(postData);

            if (parsedData?.variables?.input?.name) {
              return parsedData.variables.input.name === 'John Doe';
            }
          } catch (e) {
            console.warn('Failed to parse request postData:', e);

            return false;
          }

          return false;
        }
      }

      return false;
    });

    const responsePromise = this.page.waitForResponse(async (response) => {
      if (
        response.request().method() === 'POST' &&
        response.url().includes('customer-os-api')
      ) {
        const responseBody = await response.json();

        return responseBody.data?.contact_Update?.id !== undefined;
      }

      return false;
    });

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

    const requestPromise = this.page.waitForRequest((request) => {
      if (
        request.method() === 'POST' &&
        request.url().includes('customer-os-api')
      ) {
        const postData = request.postData();

        if (postData) {
          const parsedData = JSON.parse(postData);

          return parsedData.variables.input?.jobTitle === 'CTO';
        }
      }

      return false;
    });

    const responsePromise = this.page.waitForResponse(async (response) => {
      if (
        response.request().method() === 'POST' &&
        response.url().includes('customer-os-api')
      ) {
        const responseBody = await response.json();

        return responseBody.data?.jobRole_Create?.id !== undefined;
      }

      return false;
    });

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

    const requestPromise = this.page.waitForRequest((request) => {
      if (
        request.method() === 'POST' &&
        request.url().includes('customer-os-api')
      ) {
        const postData = request.postData();

        if (postData) {
          const parsedData = JSON.parse(postData);

          return parsedData.variables.input?.description === 'Influencer';
        }
      }

      return false;
    });

    const responsePromise = this.page.waitForResponse(async (response) => {
      if (
        response.request().method() === 'POST' &&
        response.url().includes('customer-os-api')
      ) {
        const responseBody = await response.json();

        return responseBody.data?.jobRole_Create?.id !== undefined;
      }

      return false;
    });

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

    const requestPromise = page.waitForRequest((request) => {
      if (
        request.method() === 'POST' &&
        request.url().includes('customer-os-api')
      ) {
        const postData = request.postData();

        if (postData) {
          const parsedData = JSON.parse(postData);

          return parsedData.variables.input?.url === 'www.linkedin.com/in/test';
        }
      }

      return false;
    });

    const responsePromise = page.waitForResponse(async (response) => {
      if (
        response.request().method() === 'POST' &&
        response.url().includes('customer-os-api')
      ) {
        const responseBody = await response.json();

        return responseBody.data?.contact_AddSocial?.id !== undefined;
      }

      return false;
    });

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
