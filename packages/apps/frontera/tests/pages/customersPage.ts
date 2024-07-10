import { Page, expect } from '@playwright/test';

export class CustomersPage {
  constructor(private page: Page) {}

  async ensureNoCustomerExists() {
    const noResultsElement = await this.page.waitForSelector(
      'div.flex.flex-1.relative.w-full div.pt-12.mx-auto.text-gray-700.text-center',
      { state: 'visible' },
    );
    const spanText = await this.page
      .locator(
        'div.flex.flex-1.relative.w-full div.pt-12.mx-auto.text-gray-700.text-center span.text-md.font-medium',
      )
      .textContent();
    expect(spanText).toBe('No Resultsville');

    const fullText = await noResultsElement.textContent();
    expect(fullText).toContain('Empty here in No Resultsville');
    expect(fullText).toContain(
      'Try using different keywords, checking for typos, or adjusting your filters.',
    );
  }

  async ensureNumberOfCustomersExist(numberOfCustomers) {
    const elements = await this.page.$$(
      '.flex.flex-1.relative.w-full > .top-0.left-0.inline-flex.items-center.flex-1.w-full.text-sm.absolute.border-b.bg-white.border-gray-100.transition-all.animate-fadeIn.group[data-index="0"][data-selected="false"][data-focused="false"]',
    );
    expect(
      elements.length,
      `Expected to have ${numberOfCustomers} customer(s) and found ${elements.length}`,
    ).toBe(numberOfCustomers);
  }
}
