import { setTimeout } from "timers/promises";

import { Browser } from "../browser";
import { logger } from "@/infrastructure";
import { ErrorParser } from "@/util/error";

export type Cookies = ReadonlyArray<{
  name: string;
  value: string;
  url?: string;
  domain?: string;
  path?: string;
  expires?: number;
  httpOnly?: boolean;
  secure?: boolean;
  sameSite?: "Strict" | "Lax" | "None";
}>;

export class LinkedinAutomationService {
  constructor(
    private cookies: Cookies,
    private useAgent: string,
    private proxyConfig: string,
  ) {}

  async sendConenctionInvite(
    profileUrl: string,
    message: string,
    options?: { dryRun?: boolean },
  ) {
    const browser = await Browser.getInstance(this.proxyConfig);
    const context = await browser.newContext({
      userAgent: this.useAgent,
    });
    context.setExtraHTTPHeaders({
      "browsercat-opts": this.proxyConfig,
    });
    context.addCookies(this.cookies);

    const page = await context.newPage();
    await page.goto(profileUrl);

    try {
      const profileName = page.locator("h1.text-heading-xlarge");
      await profileName.waitFor({ timeout: 10000 });
      const profileNameText = await profileName.textContent();

      const buttons = page.locator(`button[aria-label*="${profileNameText}"]`);
      await buttons.last().waitFor({ timeout: 10000 });
      await buttons.last().click();

      const sendInviteModal = page.locator("div.send-invite");
      await sendInviteModal.waitFor({ timeout: 10000 });
      const addNoteButton = sendInviteModal.locator(
        "button.artdeco-button--secondary",
      );

      await addNoteButton.click();
      const noteInput = sendInviteModal.locator("textarea#custom-message");
      await noteInput.fill(message);

      const sendInviteButton = sendInviteModal.locator(
        'button.artdeco-button--primary[aria-label="Send invitation"]',
      );

      await setTimeout(1000);
      if (!options?.dryRun) {
        await sendInviteButton.click();
      }
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in LinkedinAutomationService", {
        error: error.message,
        details: error.details,
      });
      throw error;
    } finally {
      await page.close();
    }
  }

  async getConnections() {
    const browser = await Browser.getInstance(this.proxyConfig);
    const context = await browser.newContext({
      userAgent: this.useAgent,
    });
    context.addCookies(this.cookies);

    const page = await context.newPage();

    const goToPage = async (currentPage: number) => {
      const url = `https://www.linkedin.com/search/results/people/?network=%5B%22F%22%5D&origin=FACETED_SEARCH&page=${currentPage}`;
      await page.goto(url);
    };

    const scrapeConnections = async () => {
      let accumulator: string[] = [];

      // Initial page load
      let currentPage = 1;
      await goToPage(currentPage);

      // Scroll to bottom to load pagination
      const footer = page.locator("footer.global-footer");
      await footer.scrollIntoViewIfNeeded();

      // Find out the last page number
      const pagination = page.locator("li.artdeco-pagination__indicator");
      const lastPageBtn = await pagination.last().textContent();
      const lastPage = parseInt(lastPageBtn?.trim() ?? "1");

      while (currentPage <= lastPage) {
        try {
          // Wait for results to load on the current page
          const results = page.locator(
            "ul.reusable-search__entity-result-list",
          );
          await results.first().waitFor({ timeout: 10000 });

          const current = await results
            .first()
            .locator("a.app-aware-link")
            .evaluateAll((links) => {
              const hrefs = links
                .filter(
                  (link) =>
                    !link.classList.contains(
                      "reusable-search-simple-insight__wrapping-link",
                    ) &&
                    !link.parentElement?.classList.contains(
                      "reusable-search-simple-insight__text",
                    ),
                )
                .map((link) => link.getAttribute("href") ?? "")
                .filter((href) => href.includes("/in/"));

              return Array.from(new Set(hrefs)); // Remove duplicates
            });

          accumulator = [...accumulator, ...current];

          const delayTime = Math.floor(Math.random() * 3000) + 2000;
          await setTimeout(delayTime);

          currentPage++;
          if (currentPage <= lastPage) {
            await goToPage(currentPage);
          }
        } catch (error) {
          console.error(`Error scraping page ${currentPage}:`, error);
          break;
        }
      }

      return accumulator;
    };

    try {
      const connectionUrls = await scrapeConnections();
      return connectionUrls;
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in LinkedinAutomationService", {
        error: error.message,
        details: error.details,
      });
      throw error;
    } finally {
      await page.close();
    }
  }

  async sendMessageToConnection(
    profileUrl: string,
    message: string,
    options?: { dryRun?: boolean },
  ) {
    const browser = await Browser.getInstance(this.proxyConfig);
    const context = await browser.newContext({
      userAgent: this.useAgent,
    });
    context.addCookies(this.cookies);
    const page = await context.newPage();
    await page.goto(profileUrl);

    try {
      const messageButtons = page.locator(
        'button.pvs-profile-actions__action[aria-label*="Message"]',
      );
      await messageButtons.waitFor({ timeout: 10000 });
      await messageButtons.click();

      const messageInput = page.locator("div.msg-form__contenteditable");
      await messageInput.waitFor({ timeout: 10000 });
      await messageInput.fill(message);

      const sendButton = page.locator("button.msg-form__send-button");
      await setTimeout(1000);

      if (!options?.dryRun) {
        await sendButton.click();
      }
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error in LinkedinAutomationService", {
        error: error.message,
        details: error.details,
      });
      throw error;
    } finally {
      return await page.close();
    }
  }
}
