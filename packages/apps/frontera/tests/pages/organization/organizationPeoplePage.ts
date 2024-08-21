import { Page, expect } from '@playwright/test';

import { clickLocatorsThatAreVisible } from '../../helper';

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
    'button[data-test="org-people-contact-timezone"]';

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
          const parsedData = JSON.parse(postData);

          return parsedData.variables.input?.name === 'John Doe';
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

        // console.log('addNameToContact Request data:', responseBody); // Log the request data

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

    // await this.page.waitForResponse('**/customer-os-api');
    // await this.page.keyboard.press('Enter');
    // await this.page.keyboard.press('Tab');
  }

  async addJobRolesToContact() {
    const orgPeopleContactJobRolesInput = this.page.locator(
      this.orgPeopleContactJobRoles,
    );

    // await this.page.reload();
    await orgPeopleContactJobRolesInput.click();

    await this.page.waitForSelector('[role="listbox"]', { state: 'visible' });

    const influencerOption = this.page.locator(this.jobRoleInfluencer);

    const requestPromise = this.page.waitForRequest((request) => {
      if (
        !request.url().includes('.clarity.') &&
        !request.url().includes('heapanalytics')
      )
        if (
          request.method() === 'POST' &&
          request.url().includes('customer-os-api')
          // request.url().includes('mid.customeros.ai/customer-os-api')
        ) {
          const postData = request.postData();

          if (postData) {
            const parsedData = JSON.parse(postData);

            return parsedData.variables.input?.description === 'Influencer';
          }
        }

      return false;
    });

    await influencerOption.click();
    await this.page.waitForTimeout(500);

    const responsePromise = this.page.waitForResponse(async (response) => {
      if (
        response.request().method() === 'POST' &&
        response.url().includes('customer-os-api')
      ) {
        const responseBody = await response.json();

        return responseBody.data?.jobRole_Update?.id !== undefined;
      }

      return false;
    });

    await Promise.all([requestPromise, responsePromise]);
  }

  private async ensureJobRoleIsCreated() {
    const requestPromise = this.page.waitForRequest((request) => {
      if (
        request.method() === 'POST' &&
        request.url().includes('customer-os-api')
      ) {
        const postData = request.postData();

        if (postData) {
          const parsedData = JSON.parse(postData);

          return (
            parsedData.operationName === 'addContactSocial' &&
            typeof parsedData.variables.input === 'object' &&
            Object.keys(parsedData.variables.input).length === 0
          );
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

    await Promise.all([requestPromise, responsePromise]);
  }

  async createContactFromEmpty() {
    await this.addContactEmpty();
    await this.addNameToContact();
    await this.addTitleToContact();
    // await this.ensureJobRoleIsCreated();
    await this.addJobRolesToContact();
  }
}
