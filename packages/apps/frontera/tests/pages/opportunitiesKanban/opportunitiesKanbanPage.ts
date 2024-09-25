import { Page, expect } from '@playwright/test';

import {
  createResponsePromise,
  clickLocatorsThatAreVisible,
} from '../../helper';

export class OpportunitiesKanbanPage {
  private page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  private sideNavItemOpportunitiesKanban =
    'button[data-test="side-nav-item-Opportunities"]';
  private sideNavItemOpportunitiesSelected =
    'button[data-test="side-nav-item-Opportunities"] div[aria-selected="true"]';
  private oppsKanbanHeaderOppsCount =
    'span[data-test="opps-kanban-header-opps-count"]';
  private oppsKanbanHeaderOppsEstimate =
    'span[data-test="opps-kanban-header-arr-estimate"]';
  private oppsKanbanHeaderWeightedArrEstimate =
    'span[data-test="opps-kanban-header-weighted-arr-estimate"]';
  private oppsFinderCount = 'span[data-test="opps-finder-count"]';
  private kanbanColumnIdentified = 'div[data-test="kanban-column-Identified"]';
  private kanbanColumnQualified = 'div[data-test="kanban-column-Qualified"]';
  private kanbanColumnQualifiedCards =
    'div[data-test="kanban-column-Qualified-cards"]';
  private kanbanCards = 'div[data-test="opp-kanban-card"]';
  private oppKanbanCardDots = 'button[data-test="opp-kanban-card-dots"]';
  private addOppPlusIdentified = 'button[data-test="add-opp-plus-Identified"]';
  private oppKanbanChooseOrganization =
    'input[data-test="opp-kanban-choose-organization"]';
  private kanbanClock = 'path[data-test="kanban-clock"]';

  async goToOpportunitiesKanban() {
    await clickLocatorsThatAreVisible(
      this.page,
      this.sideNavItemOpportunitiesKanban,
    );
  }

  async checkOpportunitiesKanbanHeaderValues(
    expectedActualOppsFinderCount: number,
    expectedActualOppsHeaderCount: number,
    expectedArrEstimate: number,
    expectedWeightedArrEstimate: number,
  ) {
    const actualOppsFinderCount = await this.page
      .locator(this.oppsFinderCount)
      .innerText();
    const trimmedActualOppsFinderCount =
      actualOppsFinderCount.match(/^\d+/)?.[0] || actualOppsFinderCount;
    const actualOppsHeaderCount = await this.page
      .locator(this.oppsKanbanHeaderOppsCount)
      .innerText();
    const actualArrEstimate = await this.page
      .locator(this.oppsKanbanHeaderOppsEstimate)
      .innerText();
    const trimmedArrEstimate =
      actualArrEstimate.match(/\$?([\d,]+\.?\d*)/)?.[1]?.replace(/,/g, '') ||
      actualArrEstimate;
    const actualWeightedArrEstimate = await this.page
      .locator(this.oppsKanbanHeaderWeightedArrEstimate)
      .innerText();
    const trimmedWeightedArrEstimate =
      actualWeightedArrEstimate
        .match(/\$?([\d,]+\.?\d*)/)?.[1]
        ?.replace(/,/g, '') || actualWeightedArrEstimate;

    await Promise.all([
      expect
        .soft(
          trimmedActualOppsFinderCount,
          `Expected to have ${expectedActualOppsFinderCount} opportunities in finder and ${trimmedActualOppsFinderCount} were found.`,
        )
        .toBe(expectedActualOppsFinderCount.toString()),
      expect
        .soft(
          actualOppsHeaderCount,
          `Expected to have ${expectedActualOppsHeaderCount} opportunities in header and ${actualOppsHeaderCount} were found.`,
        )
        .toBe(expectedActualOppsHeaderCount.toString()),
      expect
        .soft(
          trimmedArrEstimate,
          `Expected to have ${expectedArrEstimate} ARR Estimate and found ${trimmedArrEstimate}.`,
        )
        .toBe(expectedArrEstimate.toString()),
      expect
        .soft(
          trimmedWeightedArrEstimate,
          `Expected to have ${expectedWeightedArrEstimate} Weighted Arr Estimate and found ${trimmedWeightedArrEstimate}.`,
        )
        .toBe(expectedWeightedArrEstimate.toString()),
    ]);
  }

  // async removeOpportunity(organizationName: string) {
  //   await this.page.waitForSelector(this.kanbanColumnIdentified);
  //
  //   // Get all cards within the Identified column
  //   const cards = await this.page.$$(
  //     `${this.kanbanColumnIdentified} ${this.kanbanCards}`,
  //   );
  //
  //   console.log(`Number of cards detected: ${cards.length}`);
  //
  //   // Double-check by getting the count using evaluate
  //   const cardCount = await this.page.evaluate((selector) => {
  //     return document.querySelectorAll(selector).length;
  //   }, `${this.kanbanColumnIdentified} ${this.kanbanCards}`);
  //
  //   console.log(`Number of cards detected using evaluate: ${cardCount}`);
  //
  //   await this.page.screenshot({
  //     path: 'debug-screenshot.png',
  //     fullPage: true,
  //   });
  // }

  async addOpportunity(organizationName: string) {
    await this.page.locator(this.addOppPlusIdentified).click();

    const responsePromise = createResponsePromise(
      this.page,
      'opportunity?.metadata?.id',
      undefined,
    );

    await this.page
      .locator(this.oppKanbanChooseOrganization)
      .pressSequentially(organizationName);
    await this.page.locator(this.oppKanbanChooseOrganization).press('Enter');
    await Promise.all([responsePromise]);

    const card = this.page.locator(
      `${this.kanbanCards}:has(input[value*="${organizationName}"])`,
    );
    const kanbanColumnQualifiedCards = this.page.locator(
      this.kanbanColumnQualifiedCards,
    );

    await card.waitFor({ state: 'attached' });

    //STARD DEBUGGING
    // console.log('card visible: ', await card.isVisible());
    // console.log(
    //   'kanbanColumnQualifiedCards visible: ',
    //   await kanbanColumnQualifiedCards.isVisible(),
    // );

    // await clickLocatorsThatAreVisible(
    //   this.page,
    //   `${this.kanbanCards}:has(input[value*="${organizationName}"])`,
    //   this.kanbanColumnQualifiedCards,
    // );

    //END DEBUGGING
    // await this.page.screenshot({
    //   path: 'before.png',
    //   fullPage: true,
    // });

    await this.page.waitForTimeout(5000);
    // await card.dragTo(kanbanColumnQualifiedCards);

    // await this.page.screenshot({ path: 'screenshot.png' });

    // console.log(
    //   'clock styling before: ',
    //   await this.page.locator(this.kanbanClock).getAttribute('style'),
    // );
    await this.page.locator(this.kanbanClock).hover();
    // console.log('Hovered on the clock.');
    await this.page.screenshot({ path: 'clock-hovered.png' });
    await this.page.waitForTimeout(3000);
    await this.page.mouse.down();
    await this.page.waitForTimeout(3000);
    await kanbanColumnQualifiedCards.hover();
    // console.log(
    //   'clock styling after: ',
    //   await this.page.locator(this.kanbanClock).getAttribute('style'),
    // );
    await this.page.waitForTimeout(3000);
    await this.page.mouse.up();
    await this.page.waitForTimeout(3000);

    // await this.page.screenshot({
    //   path: 'after.png',
    //   fullPage: true,
    // });
  }
}
