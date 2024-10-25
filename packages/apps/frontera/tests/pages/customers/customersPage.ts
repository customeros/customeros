import { Page, expect } from '@playwright/test';

export class CustomersPage {
  private organizationNameInAllOrgsTable =
    'p[data-test="organization-name-in-all-orgs-table"]';

  constructor(private page: Page) {}

  async ensureNumberOfCustomersExist(numberOfCustomers: number) {
    const elements = await this.page.$$(
      '.flex.flex-1.relative.w-full > .top-0.left-0.inline-flex.items-center.flex-1.w-full.text-sm.absolute.border-b.bg-white.border-gray-100.transition-all.animate-fadeIn.group[data-index="0"][data-selected="false"][data-focused="false"]',
    );

    expect(
      elements.length,
      `Expected to have ${numberOfCustomers} customer(s) and found ${elements.length}`,
    ).toBe(numberOfCustomers);
  }

  async ensureCustomerExists(organizationName: string, exists: boolean) {
    const orgNames = this.page.locator(this.organizationNameInAllOrgsTable);

    if (exists) {
      await expect(orgNames).toContainText(organizationName);
    } else {
      const count = await orgNames
        .filter({ hasText: organizationName })
        .count();

      expect(count).toBe(0);
    }
  }
}
