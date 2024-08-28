import { Page } from '@playwright/test';

import { clickLocatorsThatAreVisible } from '../../helper';

export class OrganizationSideNavPage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private orgSideNavItemAbout = 'button[data-test="org-side-nav-item-about"]';
  private orgSideNavItemPeople = 'button[data-test="org-side-nav-item-people"]';
  private orgSideNavItemAccount =
    'button[data-test="org-side-nav-item-account"]';
  private orgSideNavItemSuccess =
    'button[data-test="org-side-nav-item-success"]';
  private orgSideNavItemIssues = 'button[data-test="org-side-nav-item-issues"]';

  async goToAbout() {
    await this.page.click(this.orgSideNavItemAbout);
  }

  async goToPeople() {
    await clickLocatorsThatAreVisible(this.page, this.orgSideNavItemPeople);
  }

  async goToAccount() {
    await clickLocatorsThatAreVisible(this.page, this.orgSideNavItemAccount);
  }

  async goToSuccess() {
    await this.page.click(this.orgSideNavItemSuccess);
  }

  async goToIssues() {
    await this.page.click(this.orgSideNavItemIssues);
  }
}
