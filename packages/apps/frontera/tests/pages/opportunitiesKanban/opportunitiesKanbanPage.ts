import { Page, expect, ElementHandle } from '@playwright/test';

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
  // private kanbanColumnIdentified = 'div[data-test="kanban-column-Identified"]';
  // private kanbanColumnQualified = 'div[data-test="kanban-column-Qualified"]';
  private kanbanColumnQualifiedCards =
    'div[data-test="kanban-column-Qualified-cards"]';
  // private kanbanColumnCommittedCards =
  //   'div[data-test="kanban-column-Committed"]';
  private kanbanCards = 'div[data-test="opp-kanban-card"]';
  private oppKanbanCardDots = 'button[data-test="opp-kanban-card-dots"]';
  private addOppPlusIdentified = 'button[data-test="add-opp-plus-Identified"]';
  private addOppPlusQualified = 'button[data-test="add-opp-plus-Qualified"]';
  private oppKanbanChooseOrganization =
    'input[data-test="opp-kanban-choose-organization"]';
  private kanbanClock = 'path[data-test="kanban-clock"]';
  private oppKanbanIcon = 'div[data-test="opp-kanban-icon"]';

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
  }

  async moveOpportunityCard(
    organizationName: string,
    destinationColumn?: string,
  ) {
    const card = this.page.locator(
      `${this.kanbanCards}:has(input[value*="${organizationName}"])`,
    );

    await card.waitFor({ state: 'attached' });

    expect(
      await card.isVisible(),
      `Card '${organizationName}' should be visible`,
    ).toBe(true);

    const { cardsColumnXCenter, cardsColumnYCenter } =
      await this.getCardsColumnCoordinates(destinationColumn);

    await this.dragCardToColumn(
      organizationName,
      cardsColumnXCenter,
      cardsColumnYCenter,
    );
  }

  private async dragCardToColumn(
    organizationName: string,
    cardsColumnXCenter: number,
    cardsColumnYCenter: number,
  ) {
    const cardElement: ElementHandle = await this.page.$(
      `${this.kanbanCards}:has(input[value*="${organizationName}"])`,
    );

    if (cardElement) {
      const cardBoundingBox = await cardElement.boundingBox();

      if (cardBoundingBox) {
        const { x, y, width, height } = cardBoundingBox;

        const clickX = x + width / 2;
        const clickY = y + height - 40;

        await this.page.mouse.move(clickX, clickY);
        await this.page.mouse.down();
        await this.page.mouse.move(cardsColumnXCenter, cardsColumnYCenter, {
          steps: 20,
        });
        await this.page.mouse.up();
      } else {
        process.stdout.write('Element is not visible or has no dimensions');
      }
    } else {
      process.stdout.write('Element not found');
    }
    await this.page.waitForTimeout(1000);
  }

  private async getCardsColumnCoordinates(destinationColumn: string) {
    const cardsColumnElement: ElementHandle = await this.page.$(
      destinationColumn,
    );
    let cardsColumnXCenter: number, cardsColumnYCenter: number;

    if (cardsColumnElement) {
      const cardsColumnBoundingBox = await cardsColumnElement.boundingBox();

      if (cardsColumnBoundingBox) {
        const { x, y, width, height } = cardsColumnBoundingBox;

        cardsColumnXCenter = x + width / 2;
        cardsColumnYCenter = y + height / 2;
      } else {
        process.stdout.write('Element is not visible or has no dimensions');
      }
    } else {
      process.stdout.write('Element not found');
    }

    return { cardsColumnXCenter, cardsColumnYCenter };
  }
}
