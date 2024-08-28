import { Page } from '@playwright/test';

export class OrganizationAboutPage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgAboutName = 'input[data-test="org-about-name"]';
  private orgAboutWww = 'input[data-test="org-about-www"]';
  private orgAaboutDescription = 'textarea[data-test="org-about-description"]';
  private orgAboutTags = 'div[data-test="org-about-tags"]';
  private orgAboutRrelationship = 'div[data-test="org-about-relationship"]';
  private orgAboutStage = 'div[data-test="org-about-stage"]';
  private orgAboutIndustry = 'div[data-test="org-about-industry"]';
  private orgAboutBusinessType = 'div[data-test="org-about-business-type"]';
  private orgAboutLastFundingRound =
    'div[data-test="org-about-last-funding-round"]';
  private orgAboutNumberOfEmployees =
    'div[data-test="org-about-number-of-employees"]';
  private orgAboutOrgOwner = 'div[data-test="org-about-org-owner"]';
  private orgAboutSocialLink = 'input[data-test="org-about-social-link"]';

  async populateAboutFields() {
    await this.page.fill(this.orgAboutWww, 'www.qweasdzxc123ads.com');
    await this.page.fill(
      this.orgAaboutDescription,
      'This org is simply the best, better than all the rest',
    );
    await this.page.fill(
      this.orgAboutSocialLink,
      'www.linkedin.com/in/qweasdzxc123ads',
    );
  }
}
