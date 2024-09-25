import { setTimeout } from "timers/promises";
import { setTimeout as setTimeoutSync } from "timers";

import { Browser } from "../browser";
import { logger } from "@/infrastructure";
import { ErrorParser, StandardError } from "@/util/error";

const Selectors = {
  profileNameHeading: "h1.text-heading-xlarge",
  connectButton: (profileNameText: string) =>
    `button[aria-label="Invite ${profileNameText} to connect"]`,
  connectDiv: (profileNameText: string) =>
    `div[aria-label="Invite ${profileNameText} to connect"]`,
  moreActionsButton: 'button[aria-label="More actions"]',
  sendInviteModal: "div.send-invite",
  noteInput: "textarea#custom-message",
  addNoteButton: "button.artdeco-button--secondary",
  sendInviteButton:
    'button.artdeco-button--primary[aria-label="Send invitation"]',
  sendWithoutNoteButton:
    'button.artdeco-button--primary[aria-label="Send without a note"]',
  moreActionsDropdown: "div.artdeco-dropdown__content-inner",
};

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
    private proxyConfig: string
  ) {}

  async sendConenctionInvite(
    profileUrl: string,
    message?: string,
    options?: { dryRun?: boolean }
  ) {
    const browser = await Browser.getInstance(this.proxyConfig, {
      debug: true,
    });
    const context = await browser.newContext({
      userAgent: this.userAgent,
    });
    await context.addCookies(this.cookies);

    const page = await context.newPage();
    await page.goto(profileUrl);

    try {
      const profileName = page.locator(Selectors.profileNameHeading);
      await profileName.waitFor({ timeout: 10000 });
      const profileNameText = await profileName.textContent();

      const connectButtons = await page.$$(
        Selectors.connectButton(profileNameText ?? "")
      );
      const moreActionsButtons = await page.$$(Selectors.moreActionsButton);

      if (connectButtons.length > 0) {
        await connectButtons?.[1].click();
      } else if (moreActionsButtons.length > 0) {
        await moreActionsButtons?.[1].click();

        const dropdown = page.locator(Selectors.moreActionsDropdown)?.last();
        await dropdown.waitFor({ timeout: 10000 });

        const connectButtons = await page.$$(
          Selectors.connectDiv(profileNameText ?? "")
        );

        if (connectButtons.length === 0) {
          logger.warn(
            "Connect button not found. Profile might be already a connection.",
            {
              source: "LinkedinAutomationService",
            }
          );
          return;
        }

        await connectButtons?.[1]?.scrollIntoViewIfNeeded();
        await connectButtons?.[1].click();
      } else {
        throw new StandardError({
          code: "INTERNAL_ERROR",
          message: "Connect button and More button missing.",
          severity: "high",
        });
      }

      const sendInviteModal = page.locator(Selectors.sendInviteModal);
      await sendInviteModal.waitFor({ timeout: 10000 });

      if (message) {
        const addNoteButton = sendInviteModal.locator(Selectors.addNoteButton);
        await addNoteButton.click();
        const noteInput = sendInviteModal.locator(Selectors.noteInput);
        await noteInput.fill(message);

        const sendInviteButton = sendInviteModal.locator(
          Selectors.sendInviteButton
        );

        await setTimeout(1000);
        if (!options?.dryRun) {
          await sendInviteButton.click();
        }
      } else {
        const sendWithoutNoteButton = sendInviteModal.locator(
          Selectors.sendWithoutNoteButton
        );

        await setTimeout(1000);
        if (!options?.dryRun) {
          await sendWithoutNoteButton.click();
        }
      }
    } catch (err) {
      throw LinkedinAutomationService.handleError(err);
    } finally {
      await page.close();
    }
  }

  async getConnections(
    startPage?: number
  ): Promise<
    [
      result: string[],
      error: StandardError | undefined,
      lastPageVisited?: number
    ]
  > {
    const browser = await Browser.getInstance(this.proxyConfig);
    const context = await browser.newContext({
      userAgent: this.userAgent,
    });
    context.addCookies(this.cookies);

    const page = await context.newPage();

    const scrollToFooter = async () => {
      const footer = page.locator("footer.global-footer");
      await footer.scrollIntoViewIfNeeded();
    };

    const clickNextButton = async () => {
      return await retry(async () => {
        const nextButton = page.locator('button[aria-label="Next"]');
        if (await nextButton.isEnabled()) {
          await nextButton.click();
          await page
            .locator("ul.reusable-search__entity-result-list")
            .first()
            .waitFor({ timeout: 60 * 1000 });
        } else {
          // return;
        }
      });
    };

    const scrapeConnections = async (
      initialPage?: number
    ): Promise<
      [
        result: string[],
        error: StandardError | undefined,
        lastPageVisited?: number
      ]
    > => {
      let accumulator: string[] = [];
      let error: StandardError | undefined;

      // Initial page load
      let currentPage = initialPage ?? 1;
      await page.goto(
        `https://www.linkedin.com/search/results/people/?network=%5B%22F%22%5D&origin=FACETED_SEARCH&page=${currentPage}`
      );

      // Scroll to bottom to load pagination
      await scrollToFooter();

      // Find out the last page number
      const pagination = page.locator("li.artdeco-pagination__indicator");
      const lastPageBtn = await pagination.last().textContent();
      const lastPage = parseInt(lastPageBtn?.trim() ?? "1");

      while (currentPage <= lastPage) {
        const scrapeCurrentPage = async () => {
          // Wait for results to load on the current page
          const results = page.locator(
            "ul.reusable-search__entity-result-list"
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
                      "reusable-search-simple-insight__wrapping-link"
                    ) &&
                    !link.parentElement?.classList.contains(
                      "reusable-search-simple-insight__text"
                    )
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
            await scrollToFooter();
            await clickNextButton();
          }
        } catch (err) {
          error = LinkedinAutomationService.handleError(err);
          logger.error(`Error scraping page ${currentPage}`, {
            source: "LinkedinAutomationService",
          });

          break;
        }
      }

      return [accumulator, error, currentPage];
    };

    try {
      return await scrapeConnections(startPage);
    } catch (err) {
      throw LinkedinAutomationService.handleError(err);
    } finally {
      await page.close();
    }
  }

  async sendMessageToConnection(
    profileUrl: string,
    message: string,
    options?: { dryRun?: boolean }
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
        'button.pvs-profile-actions__action[aria-label*="Message"]'
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
      "ERR_TOO_MANY_REDIRECTS"
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
  delay: number = 3000
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
        }
      );
      await new Promise((resolve) =>
        setTimeoutSync(resolve, exponentialBackoff)
      );
    }
  }
};
