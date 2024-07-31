import { Page } from '@playwright/test';

export class OrganizationSideNavPage {
  constructor(page: Page) {
    this.page = page;
  }

  private page: Page;

  private aboutInTheOrganizationPage =
    'button[data-test="org-side-nav-item-about"]';
  private peopleInTheOrganizationPage =
    'button[data-test="org-side-nav-item-people"]';
  private accountInTheOrganizationPage =
    'button[data-test="org-side-nav-item-account"]';
  private successInTheOrganizationPage =
    'button[data-test="org-side-nav-item-success"]';
  private issuesInTheOrganizationPage =
    'button[data-test="org-side-nav-item-issues"]';

  async goToAbout() {
    await this.page.click(this.aboutInTheOrganizationPage);
  }

  async goToPeople() {
    await this.page.click(this.peopleInTheOrganizationPage);
  }

  async goToAccount() {
    await this.page.click(this.accountInTheOrganizationPage);
  }

  async goToSuccess() {
    await this.page.click(this.successInTheOrganizationPage);
  }

  async goToIssues() {
    await this.page.click(this.issuesInTheOrganizationPage);
  }
}
