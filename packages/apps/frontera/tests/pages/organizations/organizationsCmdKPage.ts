import { Page, expect } from '@playwright/test';

export class OrganizationsCmdKPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  private organizationsHub = 'div[data-test="organization-hub"]';
  private organizationsHubSpan = 'div[data-test="organization-hub"] span';
  private organizationHubInput = 'input[data-test="organization-hub-input"]';
  // private organizationHubNavigate =
  //   'div[data-test="organization-hub-navigate"]';
  private organizationHubAddNewOrgs =
    'div[data-test="organization-hub-add-new-orgs"]';
  private organizationHubGt = 'div[data-test="organization-hub-gt"]';
  private organizationHubGo = 'div[data-test="organization-hub-go"]';
  private organizationHubGc = 'div[data-test="organization-hub-gc"]';
  private organizationHubGz = 'div[data-test="organization-hub-gz"]';
  private organizationHubGn = 'div[data-test="organization-hub-gn"]';
  private organizationHubGi = 'div[data-test="organization-hub-gi"]';
  private organizationHubGr = 'div[data-test="organization-hub-gr"]';
  private organizationHubGs = 'div[data-test="organization-hub-gs"]';
  private organizationHubGd = 'div[data-test="organization-hub-gd"]';

  async accessCmdK() {
    await this.openCmdK();

    const organizationsHub = this.page
      .locator(this.organizationsHubSpan)
      .first();

    const organizationsHubText = await organizationsHub.textContent();

    const organizationHubInput = this.page
      .locator(this.organizationHubInput)
      .first();

    const organizationHubInputText = await organizationHubInput.getAttribute(
      'placeholder',
    );

    const organizationHubAddNewOrgs = this.page
      .locator(this.organizationHubAddNewOrgs)
      .first();

    const organizationHubAddNewOrgsText =
      await organizationHubAddNewOrgs.textContent();

    const navigationGroup = this.page
      .locator(this.organizationsHub)
      .locator('div[data-test="organization-hub-navigate"] div[role="group"]')
      .first();

    const navigationItems = await navigationGroup
      .locator('div[cmdk-item]')
      .all();

    const organizationHubGt = this.page.locator(this.organizationHubGt);
    const organizationHubGtText = await organizationHubGt.textContent();
    const navigationItemTextOne = await navigationItems[0].textContent();

    const organizationHubGo = this.page.locator(this.organizationHubGo);
    const organizationHubGoText = await organizationHubGo.textContent();
    const navigationItemTextTwo = await navigationItems[1].textContent();

    const organizationHubGc = this.page.locator(this.organizationHubGc);
    const organizationHubGcText = await organizationHubGc.textContent();
    const navigationItemTextThree = await navigationItems[2].textContent();

    const organizationHubGz = this.page.locator(this.organizationHubGz);
    const organizationHubGzText = await organizationHubGz.textContent();
    const navigationItemTextFour = await navigationItems[3].textContent();

    const organizationHubGn = this.page.locator(this.organizationHubGn);
    const organizationHubGnText = await organizationHubGn.textContent();
    const navigationItemTextFive = await navigationItems[4].textContent();

    const organizationHubGi = this.page.locator(this.organizationHubGi);
    const organizationHubGiText = await organizationHubGi.textContent();
    const navigationItemTextSix = await navigationItems[5].textContent();

    const organizationHubGr = this.page.locator(this.organizationHubGr);
    const organizationHubGrText = await organizationHubGr.textContent();
    const navigationItemTextSeven = await navigationItems[6].textContent();

    const organizationHubGs = this.page.locator(this.organizationHubGs);
    const organizationHubGsText = await organizationHubGs.textContent();
    const navigationItemTextEight = await navigationItems[7].textContent();

    const organizationHubGd = this.page.locator(this.organizationHubGd);
    const organizationHubGdText = await organizationHubGd.textContent();
    const navigationItemTextNine = await navigationItems[8].textContent();

    await Promise.all([
      expect.soft(organizationsHubText.trim()).toBe('Organizations'),
      expect
        .soft(organizationHubInputText.trim())
        .toBe('Type a command or search'),
      expect
        .soft(organizationHubAddNewOrgsText.trim())
        .toBe('Add new organizations...'),
      expect.soft(navigationItems).toHaveLength(9),
      expect
        .soft(organizationHubGtText.trim())
        .toBe(navigationItemTextOne.trim()),
      expect
        .soft(organizationHubGoText.trim())
        .toBe(navigationItemTextTwo.trim()),
      expect
        .soft(organizationHubGcText.trim())
        .toBe(navigationItemTextThree.trim()),
      expect
        .soft(organizationHubGzText.trim())
        .toBe(navigationItemTextFour.trim()),
      expect
        .soft(organizationHubGnText.trim())
        .toBe(navigationItemTextFive.trim()),
      expect
        .soft(organizationHubGiText.trim())
        .toBe(navigationItemTextSix.trim()),
      expect
        .soft(organizationHubGrText.trim())
        .toBe(navigationItemTextSeven.trim()),
      expect
        .soft(organizationHubGsText.trim())
        .toBe(navigationItemTextEight.trim()),
      expect
        .soft(organizationHubGdText.trim())
        .toBe(navigationItemTextNine.trim()),
    ]);

    await this.page.keyboard.press('Escape');

    const organizationsHubCount = await this.page
      .locator(this.organizationsHubSpan)
      .count();

    expect(organizationsHubCount).toBe(0);
  }

  async verifyOrganizationsHub() {
    await this.page.waitForSelector(this.organizationsHubSpan, {
      state: 'visible',
    });

    const organizationsHub = this.page.locator(this.organizationsHubSpan);

    const organizationsHubText = await organizationsHub.textContent();

    expect(organizationsHubText.trim()).toBe('Organizations');
  }

  async verifyCmdKFinder() {
    await this.openCmdK();

    await this.page
      .locator(this.organizationHubInput)
      .pressSequentially('go to customer');

    const navigationGroup = this.page
      .locator(this.organizationsHub)
      .locator('div[data-test="organization-hub-navigate"] div[role="group"]')
      .first();

    const navigationItems = await navigationGroup
      .locator('div[cmdk-item]')
      .all();

    const organizationHubGc = this.page.locator(this.organizationHubGc);
    const organizationHubGcText = await organizationHubGc.textContent();
    const navigationItemTextOne = await navigationItems[0].textContent();

    const organizationHubGd = this.page.locator(this.organizationHubGd);
    const organizationHubGdText = await organizationHubGd.textContent();
    const navigationItemTextTwo = await navigationItems[1].textContent();

    await Promise.all([
      expect.soft(navigationItems).toHaveLength(2),
      expect
        .soft(organizationHubGcText.trim())
        .toBe(navigationItemTextOne.trim()),
      expect
        .soft(organizationHubGdText.trim())
        .toBe(navigationItemTextTwo.trim()),
    ]);
  }

  async openCmdK() {
    const isMac = process.platform === 'darwin';

    if (isMac) {
      await this.page.keyboard.down('Meta');
      await this.page.keyboard.press('KeyK');
      await this.page.keyboard.up('Meta');
    } else {
      await this.page.keyboard.down('Control');
      await this.page.keyboard.press('KeyK');
      await this.page.keyboard.up('Control');
    }
  }
}
