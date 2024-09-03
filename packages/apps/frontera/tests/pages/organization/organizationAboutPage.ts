import { Page, expect } from '@playwright/test';

import {
  writeTextInLocator,
  createRequestPromise,
  createResponsePromise,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class OrganizationAboutPage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgAboutName = 'input[data-test="org-about-name"]';
  private orgAboutWww = 'input[data-test="org-about-www"]';
  private orgAaboutDescription = 'textarea[data-test="org-about-description"]';
  private orgAboutTags = 'div[data-test="org-about-tags"]';
  private orgAboutRelationship = 'button[data-test="org-about-relationship"]';
  private relationshipNotAFit = 'div[role="menuitem"]:has-text("Not a fit")';
  // private orgAboutStage = 'div[data-test="org-about-stage"]';
  private orgAboutIndustry = 'div[data-test="org-about-industry"]';
  private industryIndependentPowerAndRenewableElectricityProducers =
    'div[role="option"]:has-text("Independent Power and Renewable Electricity Producers")';
  private orgAboutBusinessType = 'div[data-test="org-about-business-type"]';
  private businessTypeB2c = 'div[role="option"]:has-text("B2C")';
  private orgAboutLastFundingRound =
    'div[data-test="org-about-last-funding-round"]';
  private lastFundingRoundFriendsAndFamily =
    'div[role="option"]:has-text("Friends and Family")';
  private orgAboutNumberOfEmployees =
    'div[data-test="org-about-number-of-employees"]';
  private numberOfEmployees10KPlus =
    'div[role="option"]:has-text("10000+ employees")';
  private orgAboutOrgOwner = 'div[data-test="org-about-org-owner"]';
  private orgAboutSocialLink = 'input[data-test="org-about-social-link"]';
  private orgAboutSocialLinkFilledIn = 'p[data-test="org-about-social-link"]';

  async addWebsiteToOrg() {
    await clickLocatorsThatAreVisible(this.page, this.orgAboutWww);

    const requestPromise = createRequestPromise(
      this.page,
      'website',
      'www.qweasdzxc123ads.com',
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Update?.metadata?.id',
      undefined,
    );

    await this.page
      .locator(this.orgAboutWww)
      .pressSequentially('www.qweasdzxc123ads.com', { delay: 500 });

    await Promise.all([requestPromise, responsePromise]);
  }

  async addRelationshipToOrg() {
    await clickLocatorsThatAreVisible(this.page, this.orgAboutRelationship);

    await this.page.waitForSelector('[role="menuitem"]', { state: 'visible' });

    const requestPromise = createRequestPromise(
      this.page,
      'relationship',
      'NOT_A_FIT',
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Update?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(this.page, this.relationshipNotAFit);
    await Promise.all([requestPromise, responsePromise]);
  }

  async addIndustryToOrg() {
    const orgAboutIndustryInput = this.page.locator(this.orgAboutIndustry);

    await orgAboutIndustryInput.click();

    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const requestPromise = createRequestPromise(
      this.page,
      'industry',
      'Independent Power and Renewable Electricity Producers',
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Update?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(
      this.page,
      this.industryIndependentPowerAndRenewableElectricityProducers,
    );
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  async addBusinessTypeToOrg() {
    const orgAboutBusinessTypeInput = this.page.locator(
      this.orgAboutBusinessType,
    );

    await orgAboutBusinessTypeInput.click();

    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const requestPromise = createRequestPromise(this.page, 'market', 'B2C');

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Update?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(this.page, this.businessTypeB2c);
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  async addLastFundingRoundToOrg() {
    const orgAboutLastFundingRoundInput = this.page.locator(
      this.orgAboutLastFundingRound,
    );

    await orgAboutLastFundingRoundInput.click();
    await this.page.keyboard.press('f');
    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const requestPromise = createRequestPromise(
      this.page,
      'lastFundingRound',
      'FRIENDS_AND_FAMILY',
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Update?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(
      this.page,
      this.lastFundingRoundFriendsAndFamily,
    );
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  async addNumberOfEmployees() {
    const orgAboutNumberOfEmployeesInput = this.page.locator(
      this.orgAboutNumberOfEmployees,
    );

    await orgAboutNumberOfEmployeesInput.click();
    await orgAboutNumberOfEmployeesInput.pressSequentially('10');

    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const requestPromise = createRequestPromise(this.page, 'employees', 10001);

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Update?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(this.page, this.numberOfEmployees10KPlus);
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  async populateAboutFields() {
    await this.addWebsiteToOrg();
    await this.page.fill(
      this.orgAaboutDescription,
      'This org is simply the best, better than all the rest',
    );

    this.page = await writeTextInLocator(
      this.page,
      this.orgAboutTags,
      'testOrgTag',
    );
    await this.page.keyboard.press('Enter');
    await clickLocatorsThatAreVisible(this.page, this.orgAboutTags);

    await this.addRelationshipToOrg();
    await this.addIndustryToOrg();
    await this.addBusinessTypeToOrg();
    await this.addLastFundingRoundToOrg();
    await this.addNumberOfEmployees();
    await this.page.fill(
      this.orgAboutSocialLink,
      'www.linkedin.com/in/qweasdzxc123ads',
    );
    await this.page.keyboard.press('Enter');
  }

  async checkPopulatedAboutFields(organizationId: string, owner: string) {
    await this.page.reload();
    await this.page.waitForSelector(this.orgAboutName, { state: 'visible' });

    await Promise.all([
      expect
        .soft(this.page.locator(this.orgAboutName))
        .toHaveValue(organizationId),
      expect
        .soft(this.page.locator(this.orgAboutWww))
        .toHaveValue('www.qweasdzxc123ads.com'),
      expect
        .soft(this.page.locator(this.orgAaboutDescription))
        .toHaveValue('This org is simply the best, better than all the rest'),
      expect
        .soft(this.page.locator(this.orgAboutTags))
        .toContainText('testOrgTag'),
      expect
        .soft(this.page.locator(this.orgAboutRelationship))
        .toContainText('Not a fit'),
      expect
        .soft(this.page.locator(this.orgAboutIndustry))
        .toContainText('Independent Power and Renewable Electricity Producers'),
      expect
        .soft(this.page.locator(this.orgAboutBusinessType))
        .toContainText('B2C'),
      expect
        .soft(this.page.locator(this.orgAboutLastFundingRound))
        .toContainText('Friends and Family'),
      expect
        .soft(this.page.locator(this.orgAboutNumberOfEmployees))
        .toContainText('10000+ employees'),
      expect
        .soft(this.page.locator(this.orgAboutOrgOwner))
        .toContainText(owner),
      expect
        .soft(this.page.locator(this.orgAboutSocialLinkFilledIn))
        .toContainText('/qweasdzxc123ads'),
    ]);
  }
}
