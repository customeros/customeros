import { Page, expect } from '@playwright/test';

import {
  writeTextInLocator,
  ensureLocatorIsVisible,
  clickLocatorThatIsVisible,
  clickLocatorsThatAreVisible,
  clickFirstLocatorThatIsVisible,
} from '../../helper';

export class MailboxesPage {
  constructor(private page: Page) {
    this.page = page;
  }

  sideNavSettingsMailboxes = 'button[data-test="sideNav-settings-mailboxes"]';
  setUpMailboxes = 'button[data-test="set-up-mailboxes"]';
  settingsMailboxesAddNew = 'span[data-test="settings-mailboxes-add-new"]';
  settingsMailboxesDomainInput = 'input[data-test="mailboxes-domain-input"]';
  settingsMailboxesDomainsSuggestNew =
    'button[data-test="mailboxes-domains-suggest-new"]';
  settingsMailboxesDomainsList = 'div[data-test="mailboxes-domain-list"]';
  settingsmMailboxesDomainAddToCart =
    'button[data-test="mailboxes-domains-add-to-cart"]';
  settingsMailboxesRedirectUrl =
    'input[data-test="settings-mailboxes-redirect-url"]';
  settingsMailboxesRedirectUrlId = 'settings-mailboxes-redirect-url';
  settingsMailboxesFirstUsername =
    'input[data-test="settings-mailboxes-first-username"]';
  settingsMailboxesFirstUsernameId = 'settings-mailboxes-first-username';
  settingsMailboxesSecondUsername =
    'input[data-test="settings-mailboxes-second-username"]';
  settingsMailboxesSecondUsernameId = 'settings-mailboxes-second-username';
  settingsMailboxesCheckout = 'button[data-test="settings-mailboxes-checkout"]';
  stripeFieldNumberInputId = '[id="Field-numberInput"]';
  stripeFieldExpiryInputId = '[id="Field-expiryInput"]';
  stripeFieldFieldCvcInputId = '[id="Field-cvcInput"]';
  stripeFieldFieldCountryInputId = '[id="Field-countryInput"]';
  stripeFieldFieldPostalCodeInputId = '[id="Field-postalCodeInput"]';
  submitButton = 'button[typeof="submit"]';
  newMailboxes = 'button[data-test="new-mailboxes"]';
  newMailboxesHeader = 'h2[data-test="mailboxes-header"]';

  async goToMailboxes() {
    await clickLocatorsThatAreVisible(this.page, this.sideNavSettingsMailboxes);
    await ensureLocatorIsVisible(this.page, this.newMailboxesHeader);
  }

  async setupMailboxes() {
    await clickLocatorThatIsVisible(this.page, this.newMailboxes);
    await ensureLocatorIsVisible(this.page, this.settingsMailboxesAddNew);
  }

  async searchForDomains() {
    await writeTextInLocator(
      this.page,
      this.settingsMailboxesDomainInput,
      'oreoxyz',
    );
    await this.page.keyboard.press('Enter');
    await ensureLocatorIsVisible(this.page, this.settingsMailboxesDomainsList);

    const initialDomains = await this.getMailboxDomains();

    await clickLocatorThatIsVisible(
      this.page,
      this.settingsMailboxesDomainsSuggestNew,
    );

    const updatedDomains = await this.getMailboxDomains();

    const hasOverlap = updatedDomains.some((domain) =>
      initialDomains.includes(domain),
    );

    expect(
      hasOverlap,
      'Updated domains should be completely different from initial domains',
    ).toBe(false);

    expect(
      updatedDomains.length,
      'Should maintain the same number of domain suggestions',
    ).toBe(initialDomains.length);
  }

  async getMailboxDomains(): Promise<string[]> {
    await this.page.waitForSelector(
      '[data-test="mailboxes-domain-list"] .group\\/item span.text-sm',
    );

    return await this.page
      .locator('[data-test="mailboxes-domain-list"] .group\\/item span.text-sm')
      .evaluateAll((elements) => {
        return elements.map(
          (el) => el.textContent?.replace(/\s+/g, '').toLowerCase() || '',
        );
      });
  }

  async addDomainsToCart(numberOfDomains: number = 5) {
    await ensureLocatorIsVisible(this.page, this.settingsMailboxesDomainsList);

    for (let i = 0; i < numberOfDomains; i++) {
      const groupItem = this.page.locator('.group\\/item').first();

      await groupItem.hover();
      await clickFirstLocatorThatIsVisible(
        this.page,
        this.settingsmMailboxesDomainAddToCart,
      );

      await this.page.waitForTimeout(500);
    }
  }

  async setRedirectUrl() {
    await writeTextInLocator(
      this.page,
      this.settingsMailboxesRedirectUrl,
      'https://redirurl.com',
    );
    await this.page.keyboard.press('Tab');
    await expect(
      this.page.getByTestId(this.settingsMailboxesRedirectUrlId),
    ).toHaveValue('https://redirurl.com');
  }

  async setUsernames() {
    await writeTextInLocator(
      this.page,
      this.settingsMailboxesFirstUsername,
      'FirstUsername',
    );
    await this.page.keyboard.press('Tab');
    await writeTextInLocator(
      this.page,
      this.settingsMailboxesSecondUsername,
      'SecondUsername',
    );
    await this.page.keyboard.press('Tab');
    await expect(
      this.page.getByTestId(this.settingsMailboxesFirstUsernameId),
    ).toHaveValue('firstusername');
    await expect(
      this.page.getByTestId(this.settingsMailboxesSecondUsernameId),
    ).toHaveValue('secondusername');
  }

  async checkout() {
    await clickLocatorThatIsVisible(this.page, this.settingsMailboxesCheckout);
  }

  async fillInPaymentForm() {
    const stripeIframe = this.page.frameLocator(
      'iframe[title="Secure payment input frame"]',
    );

    await stripeIframe
      .locator(this.stripeFieldNumberInputId)
      .fill('4242424242424242');
    await stripeIframe.locator(this.stripeFieldExpiryInputId).fill('1234');
    await stripeIframe.locator(this.stripeFieldFieldCvcInputId).fill('123');
    await stripeIframe.locator(this.stripeFieldFieldCountryInputId).click();
    await stripeIframe.getByLabel('Country').selectOption('US');
    await stripeIframe
      .locator(this.stripeFieldFieldPostalCodeInputId)
      .fill('99999');
    await clickLocatorThatIsVisible(this.page, this.submitButton);

    const response = await this.page.waitForResponse(
      (response) =>
        response.url().includes('/customer-os-api') &&
        response.request().method() === 'POST',
      { timeout: 30000 },
    );

    const responseData = await response.json();

    expect(
      responseData?.data?.mailstack_GetPaymentIntent?.clientSecret,
    ).toBeTruthy();
  }
}
