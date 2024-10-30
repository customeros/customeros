import { randomUUID } from 'crypto';
import { Page, expect, ElementHandle } from '@playwright/test';

import {
  createRequestPromise,
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
  sideNavItemOpportunitiesSelected =
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
  private oppThreeDotsMenuWinProb =
    'div[role="menuitem"]:has-text("Set win probability")';
  private winRateSlider = 'span[role="slider"]';
  private winRateSliderBar = 'span[data-test="slider-bar"]';
  private winRateConfirm = 'button[data-test="win-rate-confirm"]';
  private winRateModal = 'div[data-test="win-rate-modal"]';
  private oppKanbanCardOpportunityName =
    'input[data-test="opp-kanban-card-opportunity-name"]';
  private cardArrEstimate = 'input[data-test="card-arr-estimate"]';

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
          parseFloat(trimmedWeightedArrEstimate).toFixed(2),
          `Expected to have ${expectedWeightedArrEstimate} Weighted Arr Estimate and found ${trimmedWeightedArrEstimate}.`,
        )
        .toBe(expectedWeightedArrEstimate.toFixed(2)),
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
    await this.page
      .locator(this.oppKanbanChooseOrganization)
      .pressSequentially(organizationName);
    await this.page.locator(this.oppKanbanChooseOrganization).press('Enter');

    const responseMutationPromise = createResponsePromise(
      this.page,
      'opportunity_Create?.metadata?.id',
      undefined,
    );
    const queryResponsePromise1 = createResponsePromise(
      this.page,
      'opportunity?.metadata?.id',
      undefined,
    );
    const queryResponsePromise2 = createResponsePromise(
      this.page,
      'opportunity?.metadata?.id',
      undefined,
    );

    await responseMutationPromise;
    await queryResponsePromise1;
    await queryResponsePromise2;
    await this.page.waitForTimeout(1500);
  }

  async moveOpportunityCard(
    opportunityName: string,
    destinationColumn?: string,
  ) {
    const card = this.page.locator(
      `${this.kanbanCards}:has(input[value*="${opportunityName}"])`,
    );

    await card.waitFor({ state: 'attached' });
    await card.scrollIntoViewIfNeeded();

    expect(
      await card.isVisible(),
      `Card '${opportunityName}' should be visible`,
    ).toBe(true);

    const { cardsColumnXCenter, cardsColumnYCenter } =
      await this.getCardsColumnCoordinates(destinationColumn);

    await this.dragCardToColumn(
      opportunityName,
      cardsColumnXCenter,
      cardsColumnYCenter,
    );
  }

  private async dragCardToColumn(
    opportunityName: string,
    cardsColumnXCenter: number,
    cardsColumnYCenter: number,
  ) {
    const cardElement: ElementHandle = await this.page.$(
      `${this.kanbanCards}:has(input[value*="${opportunityName}"])`,
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

  async setWinRates(winRateColumn: string, number: number) {
    await clickLocatorsThatAreVisible(
      this.page,
      winRateColumn,
      this.oppThreeDotsMenuWinProb,
    );
    await this.page.waitForSelector(`${this.winRateModal}[data-state="open"]`, {
      state: 'attached',
    });

    let widthPercentage: number | null = null;

    const winRateSliderBarElement: ElementHandle = await this.page.$(
      this.winRateSliderBar,
    );

    if (winRateSliderBarElement) {
      const winRateSliderBarBoundingBox =
        await winRateSliderBarElement.boundingBox();

      if (winRateSliderBarBoundingBox) {
        const { width } = winRateSliderBarBoundingBox;

        widthPercentage = (width * number) / 100;
      } else {
        process.stdout.write('Element is not visible or has no dimensions');
      }
    } else {
      process.stdout.write('Element not found');
    }

    const winRateSliderElement: ElementHandle = await this.page.$(
      this.winRateSlider,
    );

    if (winRateSliderElement) {
      const winRateSliderBoundingBox = await winRateSliderElement.boundingBox();

      if (winRateSliderBoundingBox) {
        const { x, y, width, height } = winRateSliderBoundingBox;

        const clickXSource = x + width / 2;
        const clickXDestination = x + widthPercentage;
        const clickY = y + height / 2;

        await this.page.mouse.move(clickXSource, clickY);
        await this.page.mouse.down();
        await this.page.mouse.move(clickXDestination, clickY, {
          steps: 400,
        });
        await this.page.mouse.up();
      } else {
        process.stdout.write('Element is not visible or has no dimensions');
      }
    } else {
      process.stdout.write('Element not found');
    }
    await this.page.waitForTimeout(1000);
    await this.page.waitForSelector(this.winRateConfirm, {
      state: 'attached',
    });
    expect(await this.page.isEnabled(this.winRateConfirm)).toBe(true);
    await this.page.locator(this.winRateConfirm).click();
    await this.page.waitForTimeout(3000);
  }

  private async clickConfirm() {
    const button: ElementHandle = await this.page.$(this.winRateConfirm);

    if (button) {
      const box = await button.boundingBox();

      if (box) {
        const { x, y, width, height } = box;

        const clickX = x + width / 2;
        const clickY = y + height / 2;

        await this.page.mouse.move(clickX, clickY);
        await this.page.mouse.down();
        await this.page.mouse.up();
      }
    }
  }

  async updateOpportunityName(organizationName: string) {
    const card = this.page.locator(
      `${this.kanbanCards}:has(p:text("${organizationName}"))`,
    );

    await card.waitFor({ state: 'attached' });

    expect(
      await card.isVisible(),
      `Card '${organizationName}' should be visible`,
    ).toBe(true);

    const opportunityName = randomUUID();

    expect(1, 'opportunityName: ' + opportunityName).toBe(1);

    await this.page.waitForTimeout(3000);

    const requestPromise = createRequestPromise(
      this.page,
      'name',
      opportunityName,
    );

    const responsePromise = createResponsePromise(
      this.page,
      'opportunity_Save?.metadata?.id',
      undefined,
    );

    const nameInput = card.locator(this.oppKanbanCardOpportunityName);

    await nameInput.dblclick();
    await nameInput.pressSequentially(opportunityName, { delay: 200 });

    expect(1, 'ONE IS ONE').toBe(1);
    await Promise.all([requestPromise, responsePromise]);

    return opportunityName;
  }

  async setOpportunityArrEstimate(opportunityName: string) {
    const card = this.page.locator(
      `[data-test="opp-kanban-card"]:has(input[data-test="opp-kanban-card-opportunity-name"][value^="${opportunityName}"])`,
    );

    await card.waitFor({ state: 'attached' });

    expect(
      await card.isVisible(),
      `Card '${opportunityName}' should be visible`,
    ).toBe(true);

    const arrEstimateInput = card.locator('[data-test="card-arr-estimate"]');

    await arrEstimateInput.dblclick();
    await arrEstimateInput.pressSequentially('5');
    await arrEstimateInput.press('Tab');
  }
}
