import { Page, expect } from '@playwright/test';

import {
  writeTextInLocator,
  createRequestPromise,
  createResponsePromise,
  clickLocatorThatIsVisible,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class CompanyAboutPage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgAboutName = 'input[data-test="org-about-name"]';
  private orgAboutDomainEmpty = 'span[data-test="org-about-domain-empty"]';
  private orgAboutDomainFilled = 'span[data-test="org-about-domain-filled"]';
  private orgAboutDescription = 'div[data-test="org-about-description"]';
  private orgAboutTags = 'div[data-test="org-about-tags"]';
  private orgAboutRelationship = 'button[data-test="org-about-relationship"]';
  private relationshipNotAFit = 'div[role="menuitem"]:has-text("Not a fit")';
  private orgAboutStage = 'div[data-test="org-about-stage"]';
  private orgAboutIndustry = 'span[data-test="org-about-industry"]';
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
  private orgAboutOwner = 'div[data-test="org-about-org-owner"]';
  private orgAboutOwnerCustomerosFeTesting =
    'div[role="option"]:has-text("customeros.fe.testing")';
  private orgAboutSocialLink = 'div[data-test="org-about-social-link"]';
  // private readonly socialLinkDataTest =
  //   this.orgAboutSocialLink.match(/data-test="([^"]+)"/)?.[1] ??
  //   'org-about-social-link';
  private orgAboutSocialLinkFilledIn = 'p[data-test="org-about-social-link"]';
  private orgAboutSocialLinkEmpty =
    'input[data-test="org-about-social-link"][placeholder="Social link"]';
  private orgTest = 'input[data-test="test"]';
  private addDomainInput = 'input[data-test="add-domain-input"]';
  private addDomain = 'button[data-test="add-domain"]';
  private mergeOrganizations = 'button[data-test="merge-organizations"]';

  async updateDomainToOrg(domains: string) {
    await this.addDomainToOrg(domains, 'update');
  }

  async addDomainToOrg(domain: string, domainToOrgOperation: string = 'add') {
    const domainToOrgOperationLocator =
      domainToOrgOperation === 'add'
        ? this.orgAboutDomainEmpty
        : this.orgAboutDomainFilled;

    await clickLocatorsThatAreVisible(this.page, domainToOrgOperationLocator);

    const addDomainRequestPromise = createRequestPromise(
      this.page,
      'domain',
      domain,
    );

    const addDomainResponsePromise = createResponsePromise(
      this.page,
      'checkDomain?.domainOrganizationId',
      undefined,
    );

    const input = this.page.locator(this.addDomainInput);

    await this.page.waitForTimeout(500);
    await input.pressSequentially(domain, { delay: 200 });
    await clickLocatorThatIsVisible(this.page, this.addDomain);

    await Promise.all([addDomainRequestPromise, addDomainResponsePromise]);

    const mergeOrganizationsResponsePromise = createResponsePromise(
      this.page,
      'organization_Merge',
      undefined,
    );

    await clickLocatorThatIsVisible(this.page, this.mergeOrganizations);

    await Promise.all([mergeOrganizationsResponsePromise]);
    await this.page.waitForTimeout(5000);
    await this.page.reload(); // needed when the test runs locally
  }

  async addRelationshipToOrg(relationship: string) {
    await clickLocatorsThatAreVisible(this.page, this.orgAboutRelationship);

    await this.page.waitForSelector('[role="menuitem"]', { state: 'visible' });

    const requestPromise = createRequestPromise(
      this.page,
      'input?.relationship',
      relationship,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Save?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(this.page, this.relationshipNotAFit);
    await Promise.all([requestPromise, responsePromise]);
  }

  async addIndustryToOrg(industry: string) {
    const orgAboutIndustryInput = this.page.locator(this.orgAboutIndustry);

    await orgAboutIndustryInput.click();

    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const requestPromise = createRequestPromise(
      this.page,
      'input?.industry',
      industry,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Save?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(
      this.page,
      this.industryIndependentPowerAndRenewableElectricityProducers,
    );
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  async addBusinessTypeToOrg(orgAboutBusinessType: string) {
    const orgAboutBusinessTypeInput = this.page.locator(
      this.orgAboutBusinessType,
    );

    await orgAboutBusinessTypeInput.click();

    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const requestPromise = createRequestPromise(
      this.page,
      'input?.market',
      orgAboutBusinessType,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Save?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(this.page, this.businessTypeB2c);
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  async addLastFundingRoundToOrg(orgAboutLastFundingRound: string) {
    const orgAboutLastFundingRoundInput = this.page.locator(
      this.orgAboutLastFundingRound,
    );

    await orgAboutLastFundingRoundInput.click();
    await this.page.keyboard.press('f');
    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const requestPromise = createRequestPromise(
      this.page,
      'input?.lastFundingRound',
      orgAboutLastFundingRound,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Save?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(
      this.page,
      this.lastFundingRoundFriendsAndFamily,
    );
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  async addNumberOfEmployees(orgAboutNumberOfEmployees: number) {
    const orgAboutNumberOfEmployeesInput = this.page.locator(
      this.orgAboutNumberOfEmployees,
    );

    await orgAboutNumberOfEmployeesInput.click();
    await orgAboutNumberOfEmployeesInput.pressSequentially('10');

    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const requestPromise = createRequestPromise(
      this.page,
      'input?.employees',
      orgAboutNumberOfEmployees,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Save?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(this.page, this.numberOfEmployees10KPlus);
    await this.page.waitForTimeout(500);
    await Promise.all([requestPromise, responsePromise]);
  }

  private async addOwner() {
    const orgAboutOwnerInput = this.page.locator(this.orgAboutOwner);

    await orgAboutOwnerInput.click();
    await this.page.keyboard.press('c');
    await this.page.waitForSelector('[role="option"]', { state: 'visible' });

    const responsePromise = createResponsePromise(
      this.page,
      'organization_Save?.metadata?.id',
      undefined,
    );

    await clickLocatorsThatAreVisible(
      this.page,
      this.orgAboutOwnerCustomerosFeTesting,
    );
    await this.page.waitForTimeout(500);
    await Promise.all([responsePromise]);
  }

  async populateAboutFields(update: {
    name: string;
    domain: string;
    orgAboutTags: string;
    orgAboutOwner: string;
    orgAboutIndustry: string;
    orgAboutDescription: string;
    orgAboutBusinessType: string;
    orgAboutSocialLinkEmpty: string;
    orgAboutRelationshipRequest: string;
    orgAboutLastFundingRoundRequest: string;
    orgAboutNumberOfEmployeesRequest: number;
  }) {
    await this.page.fill(this.orgAboutName, update.name);

    await this.updateDomainToOrg(update.domain);

    await this.page.fill(this.orgAboutDescription, update.orgAboutDescription);

    this.page = await writeTextInLocator(
      this.page,
      this.orgAboutTags,
      update.orgAboutTags,
    );
    await this.page.keyboard.press('Enter');
    await clickLocatorsThatAreVisible(this.page, this.orgAboutTags);

    await this.addRelationshipToOrg(update.orgAboutRelationshipRequest);
    await this.addIndustryToOrg(update.orgAboutIndustry);
    await this.addBusinessTypeToOrg(update.orgAboutBusinessType);
    await this.addLastFundingRoundToOrg(update.orgAboutLastFundingRoundRequest);
    await this.addNumberOfEmployees(update.orgAboutNumberOfEmployeesRequest);
    await this.addOwner();
    await this.page.fill(
      this.orgAboutSocialLinkEmpty,
      update.orgAboutSocialLinkEmpty,
    );
    await this.page.keyboard.press('Enter');
  }

  async checkPopulatedAboutFields(update: {
    name: string;
    domain: string;
    orgAboutTags: string;
    orgAboutOwner: string;
    orgAboutIndustry: string;
    orgAboutDescription: string;
    orgAboutRelationship: string;
    orgAboutBusinessType: string;
    orgAboutSocialLinkEmpty: string;
    orgAboutLastFundingRound: string;
    orgAboutNumberOfEmployees: string;
  }) {
    await this.page.reload();
    await this.page.waitForSelector(this.orgAboutName, { state: 'visible' });

    await Promise.all([
      expect
        .soft(this.page.locator(this.orgAboutName))
        .toHaveValue(update.name),
      //TODO: waiting for the fix of the issue [COS-5192: Website save fails to get saved](https://linear.app/customer-os/issue/COS-5192/website-save-fails-to-get-saved)
      // expect
      //   .soft(this.page.locator(this.orgAboutDomainEmpty))
      //   .toHaveValue(update.website),
      expect
        .soft(this.page.locator(this.orgAboutDescription))
        .toHaveValue(update.orgAboutDescription),
      expect
        .soft(this.page.locator(this.orgAboutTags))
        .toContainText(update.orgAboutTags),
      expect
        .soft(this.page.locator(this.orgAboutRelationship))
        .toContainText(update.orgAboutRelationship),
      expect(this.page.locator(this.orgAboutStage)).toHaveCount(0),
      expect
        .soft(this.page.locator(this.orgAboutIndustry))
        .toContainText(update.orgAboutIndustry),
      expect
        .soft(this.page.locator(this.orgAboutBusinessType))
        .toContainText(update.orgAboutBusinessType),
      expect
        .soft(this.page.locator(this.orgAboutLastFundingRound))
        .toContainText(update.orgAboutLastFundingRound),
      expect
        .soft(this.page.locator(this.orgAboutNumberOfEmployees))
        .toContainText(update.orgAboutNumberOfEmployees),
      expect.soft(this.page.locator(this.orgAboutSocialLink).count()).toBe(5),
      // expect
      //   .soft(this.getSocialLinkLocator('facebook.com/cognyte').count())
      //   .toBeGreaterThan(0),
      // expect
      //   .soft(this.getSocialLinkLocator('/silkwebro').count())
      //   .toBeGreaterThan(0),
      // expect
      //   .soft(this.getSocialLinkLocator('youtube.com/c/Gomagro').count())
      //   .toBeGreaterThan(0),
      // expect
      //   .soft(this.getSocialLinkLocator('/992874').count())
      //   .toBeGreaterThan(0),
      // expect
      //   .soft(this.getSocialLinkLocator(update.orgAboutSocialLinkEmpty))
      //   .toHaveCount(1),
      expect
        .soft(this.page.locator(this.orgAboutOwner))
        .toContainText(update.orgAboutOwner),
    ]);
  }

  async enrichCompany(domain: string) {
    await this.addDomainToOrg(domain);
    await expect(this.page.locator(this.orgAboutDomainFilled)).toHaveText(
      domain,
    );
    //   await this.page.waitForLoadState('networkidle');
    //   await this.page.locator(this.orgAboutName).waitFor({ state: 'visible' });
  }

  // private getSocialLinkLocator(exactText: string) {
  //   return this.page.locator(
  //     `p[data-test="${this.socialLinkDataTest}"]:has-text("${exactText}")`,
  //   );
  // }

  async checkEnrichedAboutFields(create: { name: string; domain: string }) {
    // const socialLinks = ['/silkwebro', 'youtube.com/c/Gomagro', '/992874'];
    //
    // // First scroll each social link into view
    // for (const link of socialLinks) {
    //   const locator = this.getSocialLinkLocator(link);
    //
    //   await locator.scrollIntoViewIfNeeded();
    // }

    // Then check all conditions
    await Promise.all([
      // expect
      //   .soft(this.page.locator(this.orgAboutName))
      //   .toHaveValue(create.name),
      expect
        .soft(this.page.locator(this.orgAboutDomainFilled))
        .toHaveText(create.domain),
      expect.soft(this.page.locator(this.orgAboutDescription)).not.toBeEmpty(),
      expect
        .soft(this.page.locator(this.orgAboutRelationship))
        .toContainText('Prospect'),
      expect.soft(this.page.locator(this.orgAboutStage)).toContainText('Lead'),
      expect
        .soft(this.page.locator(this.orgAboutIndustry))
        .toHaveText('Industry not found yet'),
      expect
        .soft(await this.page.locator(this.orgAboutSocialLink).count())
        .toBe(5),
      // expect
      //   .soft(await this.getSocialLinkLocator('/silkwebro').count())
      //   .toBeGreaterThan(0),
      // expect
      //   .soft(await this.getSocialLinkLocator('youtube.com/c/Gomagro').count())
      //   .toBeGreaterThan(0),
      // expect
      //   .soft(await this.getSocialLinkLocator('/992874').count())
      //   .toBeGreaterThan(0),
    ]);
  }
}
