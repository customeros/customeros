import { setTimeout } from "timers/promises";
import { setTimeout as setTimeoutSync } from "timers";

import { Browser } from "../browser";
import { logger } from "@/infrastructure";
import { ErrorParser, StandardError } from "@/util/error";

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
    private userAgent: string,
    private proxyConfig: string,
  ) {}

  async sendConenctionInvite(
    profileUrl: string,
    message: string,
    options?: { dryRun?: boolean },
  ) {
    const browser = await Browser.getInstance(this.proxyConfig);
    const context = await browser.newContext({
      userAgent: this.userAgent,
    });
    await context.addCookies(this.cookies);

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
      throw LinkedinAutomationService.handleError(err);
    } finally {
      await page.close();
    }
  }

  async getConnections(): Promise<[string[], StandardError | undefined]> {
    const browser = await Browser.getInstance(this.proxyConfig);
    const context = await browser.newContext({
      userAgent: this.userAgent,
    });
    context.addCookies(this.cookies);

    const page = await context.newPage();

    const goToPage = async (currentPage: number) => {
      return await retry(async () => {
        const url = `https://www.linkedin.com/search/results/people/?network=%5B%22F%22%5D&origin=FACETED_SEARCH&page=${currentPage}`;
        return await page.goto(url);
      });
    };

    const scrapeConnections = async (): Promise<
      [string[], StandardError | undefined]
    > => {
      let accumulator: string[] = [];
      let error: StandardError | undefined;

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
        const scrapeCurrentPage = async () => {
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
                .filter((href) => href.includes("/in/"))
                .map((raw) => raw.split("?")[0] + "/");

              return Array.from(new Set(hrefs)); // Remove duplicates
            });

          accumulator = [...accumulator, ...current];
        };

        try {
          await retry(scrapeCurrentPage);

          const delayTime = Math.floor(Math.random() * 3000) + 2000;
          await setTimeout(delayTime);

          currentPage++;
          if (currentPage <= lastPage) {
            await goToPage(currentPage);
          }
        } catch (err) {
          error = LinkedinAutomationService.handleError(err);
          logger.error(`Error scraping page ${currentPage}`, {
            source: "LinkedinAutomationService",
          });

          break;
        }
      }

      return [accumulator, error];
    };

    try {
      return await scrapeConnections();
    } catch (err) {
      throw LinkedinAutomationService.handleError(err);
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
      userAgent: this.userAgent,
    });

    await context.addCookies(this.cookies);
    const page = await context.newPage();

    try {
      await page.goto(profileUrl, { timeout: 60 * 1000 });
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
      throw LinkedinAutomationService.handleError(err);
    } finally {
      await page.close();
    }
  }

  private static handleError(err: unknown): StandardError {
    const error = ErrorParser.parse(err);

    const isTooManyRedirectsErr = error.details?.includes(
      "ERR_TOO_MANY_REDIRECTS",
    );

    if (isTooManyRedirectsErr) {
      const tooManyRedirects = new StandardError({
        code: "EXTERNAL_ERROR",
        reference: "S001",
        details: error.details,
        message: "Too many redirects: session token might be invalid.",
        severity: "critical",
      });

      logger.error("Too many redirects: session token might be invalid.", {
        error: error.message,
        details: error.reference,
        source: "LinkedinAutomationService",
      });

      return tooManyRedirects;
    }

    logger.error("Error in LinkedinAutomationService", {
      error: error.message,
      details: error.details,
      source: "LinkedinAutomationService",
    });

    return error;
  }
}

const retry = async (
  fn: () => Promise<any>,
  retries: number = 4,
  delay: number = 1000,
) => {
  let attempt = 0;
  while (attempt < retries) {
    try {
      return await fn();
    } catch (err) {
      attempt++;
      if (attempt >= retries) {
        throw err; // If all retries fail, throw the error
      }

      const exponentialBackoff = delay * Math.pow(2, attempt);
      logger.info(
        `Retrying after ${exponentialBackoff}ms... (${attempt}/${retries})`,
        {
          source: "LinkedinAutomationService",
        },
      );
      await new Promise((resolve) =>
        setTimeoutSync(resolve, exponentialBackoff),
      );
    }
  }
};
