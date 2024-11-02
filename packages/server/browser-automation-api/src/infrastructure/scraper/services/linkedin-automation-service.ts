import AdmZip from "adm-zip";
import { Writable } from "stream";
import { FrameLocator } from "playwright";
import { setTimeout } from "timers/promises";
import { setTimeout as setTimeoutSync } from "timers";
import { TimeUtils } from '@/util/utilities';
import { Browser } from "../browser";
import { logger } from "@/infrastructure";
import { ErrorParser, StandardError } from "@/util/error";

const Selectors = {
  profileNameHeading: "h1.inline.t-24.v-align-middle.break-words",
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
    private proxyConfig: string,
  ) {}

  async sendConenctionInvite(
    profileUrl: string,
    message?: string,
    options?: { dryRun?: boolean },
  ) {
    const browser = await Browser.getFreshInstance(this.proxyConfig);
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

      await setTimeout(10 * 1000);

      const connectButtons = await page.$$(
        Selectors.connectButton(profileNameText ?? ""),
      );
      const moreActionsButtons = await page.$$(Selectors.moreActionsButton);

      if (connectButtons.length > 0) {
        await connectButtons?.[1].click();
      } else if (moreActionsButtons.length > 0) {
        await moreActionsButtons?.[1].click();

        const dropdown = page.locator(Selectors.moreActionsDropdown)?.last();
        await dropdown.waitFor({ timeout: 10000 });

        const connectButtons = await page.$$(
          Selectors.connectDiv(profileNameText ?? ""),
        );

        if (connectButtons.length === 0) {
          logger.warn(
            "Connect button not found. Profile might be already a connection.",
            {
              source: "LinkedinAutomationService",
            },
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
          Selectors.sendInviteButton,
        );

        await setTimeout(1000);
        if (!options?.dryRun) {
          await sendInviteButton.click();
        }
      } else {
        const sendWithoutNoteButton = sendInviteModal.locator(
          Selectors.sendWithoutNoteButton,
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
    startPage?: number,
  ): Promise<
    [
      result: string[],
      error: StandardError | undefined,
      lastPageVisited?: number,
    ]
  > {
    const browser = await Browser.getFreshInstance(this.proxyConfig);
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
      initialPage?: number,
    ): Promise<
      [
        result: string[],
        error: StandardError | undefined,
        lastPageVisited?: number,
      ]
    > => {
      let accumulator: string[] = [];
      let error: StandardError | undefined;

      // Initial page load
      let currentPage = initialPage ?? 1;
      await page.goto(
        `https://www.linkedin.com/search/results/people/?network=%5B%22F%22%5D&origin=FACETED_SEARCH&page=${currentPage}`,
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
    options?: { dryRun?: boolean },
  ) {
    const browser = await Browser.getFreshInstance(this.proxyConfig);
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

  async retrieveMessages(profileUrl: string) {
    const browser = await Browser.getFreshInstance(this.proxyConfig);
    const context = await browser.newContext({
      userAgent: this.userAgent,
    });

    let lastKnownTime: string | null = null;
    await context.addCookies(this.cookies);
    const page = await context.newPage();

    try {
      await page.goto(profileUrl, { timeout: 60 * 1000 });

      const btn = page.locator('button.pvs-profile-actions__action', { hasText: 'Message' });
      await btn.waitFor({ timeout: 10000 });
      await btn.click();

      await page.waitForSelector('.msg-s-message-list', { timeout: 10000 });

      await page.evaluate(() => {
        const messageList = document.querySelector('.msg-s-message-list');
        if (messageList) {
          messageList.scrollTop = 0;
        }
      });
      await page.waitForTimeout(1000);

      await page.evaluate(() => {
        const messageList = document.querySelector('.msg-s-message-list');
        if (messageList) {
          messageList.scrollTop = messageList.scrollHeight;
        }
      });
      await page.waitForTimeout(2000);

      const messages = [];
      let lastValidName = '';
      let lastValidTime = '';
      let currentTime = '';
      let timeIndex = 1;
      const thisYear = new Date().getFullYear();  // Current year
      let currentYear = thisYear;  // Default to current year
      let currentDate: Date | null = null;

      const messageElements = await page.locator('li.msg-s-message-list__event').all();
      logger.info(`Found ${messageElements.length} message elements`, { source: "LinkedinService" });

      for (const element of messageElements) {
        try {
          const elementInfo = await element.evaluate((el) => {
            const nameEl = el.querySelector('.msg-s-message-group__name');
            const timeEl = el.querySelector('time.msg-s-message-group__timestamp');
            const msgEl = el.querySelector('p.msg-s-event-listitem__body');
            const linkPreviewEl = el.querySelector('.msg-s-event-listitem__content-preview-container');
            const dateHeadingEl = el.querySelector('.msg-s-message-list__time-heading');

            return {
              name: nameEl?.textContent?.trim() || '',
              time: timeEl?.textContent?.trim() || '',
              message: (() => {
                let sanitizedMessage = msgEl?.textContent?.trim() || '';
                let previous;
                do {
                  previous = sanitizedMessage;
                  sanitizedMessage = sanitizedMessage.replace(/<!---->|<.*?>/g, '');
                } while (sanitizedMessage !== previous);
                return sanitizedMessage;
              })() || '',
              altName: el.querySelector('.msg-s-message-group__profile-link')?.textContent?.trim() || '',
              altMessage: el.querySelector('.msg-s-event-listitem__content-preview-container')?.textContent?.trim() || '',
              hasLinkPreview: !!linkPreviewEl,
              dateHeading: dateHeadingEl?.textContent?.trim() || ''
            };
          });

          // Parse the date heading and update current date context
          if (elementInfo.dateHeading) {
            // Check if this heading contains a full date with year
            const hasExplicitYear = elementInfo.dateHeading.match(/\d{4}/);

            const parsedDate = TimeUtils.parseDateHeading(elementInfo.dateHeading, currentYear);
            if (parsedDate) {
              if (hasExplicitYear) {
                // If the date has an explicit year, use it
                currentYear = parsedDate.year;
              } else {
                // If no explicit year, use current year
                parsedDate.year = thisYear;
              }
              currentDate = TimeUtils.createDate(parsedDate);
            }
          }

          const finalName = elementInfo.name || elementInfo.altName || lastValidName;

          const timeToConvert = elementInfo.time || lastKnownTime || '';
          const finalTime = TimeUtils.convertToZuluTime(timeToConvert, currentDate);

// Update last known time only if we have a non-empty time
          if (elementInfo.time) {
            lastKnownTime = elementInfo.time;
          }

          let finalMessage;
          if (elementInfo.message) {
            finalMessage = elementInfo.message;
          } else if (elementInfo.hasLinkPreview) {
            finalMessage = "External link";
          } else {
            finalMessage = elementInfo.altMessage || "Unable to parse message";
          }

          if (!finalMessage) {
            continue;
          }

          // Reset timeIndex if time changes
          if (finalTime !== currentTime) {
            currentTime = finalTime;
            timeIndex = 1;
          }

          messages.push({
            name: finalName,
            time: finalTime,
            timeIndex,
            message: finalMessage
          });

          lastValidName = finalName;
          lastValidTime = finalTime;
          timeIndex++;

        } catch (err) {
          logger.error(`Error processing message element: ${err}`, { source: "LinkedinService" });
        }
      }

      return messages;
    } catch (err) {
      throw LinkedinAutomationService.handleError(err);
    } finally {
      await browser.close();
    }
  }

  async getConnectionsNew(): Promise<
    [results: string[], error: StandardError | null]
  > {
    const browser = await Browser.getFreshInstance(this.proxyConfig, {
      debugBrowserCat: true,
    });

    const context = await browser.newContext({
      userAgent: this.userAgent,
    });

    await context.addCookies(this.cookies);
    const page = await context.newPage();

    let results: string[] = [];
    let error: StandardError | null = null;

    try {
      await page.goto(
        `https://www.linkedin.com/mynetwork/invite-connect/connections/`,
        { timeout: 60 * 1000 },
      );

      const totalConnectionsText = await page
        .locator("header.mn-connections__header")
        .innerText();

      const totalConnections = parseInt(
        totalConnectionsText.replace(/\D/g, ""),
        10,
      );

      let hasMoreResults = true;
      // let lastScrollHeight = 0;

      const getRandomDelay = (min: number, max: number) =>
        Math.floor(Math.random() * (max - min + 1)) + min;

      const getRandomScrollHeight = async () => {
        return await page.evaluate(() => {
          return (
            Math.floor(Math.random() * window.innerHeight * 0.5) +
            window.innerHeight * 0.5
          );
        });
      };

      const smoothScroll = async (
        distance: number,
        direction: "up" | "down",
      ) => {
        let scrolled = 0;
        const step = distance / 30; // Divide the distance into smaller steps

        while (scrolled < distance) {
          await page.evaluate(
            ({ scrollStep, direction }) => {
              window.scrollBy(0, direction === "up" ? -scrollStep : scrollStep);
            },
            { scrollStep: step, direction },
          );

          scrolled += step;
          await page.waitForTimeout(getRandomDelay(0, 200)); // Wait a bit between steps
        }
      };

      const closeChatBubbles = async () => {
        const activeChatBubbles = page.locator(
          "div.msg-overlay-conversation-bubble--is-active",
        );
        const bubbleCount = await activeChatBubbles.count();

        if (bubbleCount === 0) {
          return;
        }

        const bubbles = await activeChatBubbles.all();

        for (const bubble of bubbles) {
          const chatOverlayHeader = bubble.locator(
            "div.msg-overlay-bubble-header__badge-container",
          );
          await chatOverlayHeader.waitFor();
          await chatOverlayHeader.click();
          await page.waitForTimeout(getRandomDelay(1000, 3000));
        }
      };

      const humanizedMouseMove = async (
        startX: number,
        startY: number,
        endX: number,
        endY: number,
        steps = 30,
      ) => {
        let curX = startX;
        let curY = startY;

        for (let i = 0; i < steps; i++) {
          const progress = i / steps;
          const deltaX =
            (endX - startX) * (progress + Math.sin(progress * Math.PI) * 0.2); // Curve
          const deltaY =
            (endY - startY) * (progress + Math.cos(progress * Math.PI) * 0.2);

          const jitterX = Math.random() * 2 - 1;
          const jitterY = Math.random() * 2 - 1;

          curX = startX + deltaX + jitterX;
          curY = startY + deltaY + jitterY;

          await page.mouse.move(curX, curY);
          await page.waitForTimeout(getRandomDelay(10, 50)); // Vary speed
        }

        // Final move to the exact endpoint
        await page.mouse.move(endX, endY);
      };

      const moveMouseNaturally = async () => {
        const startX = Math.floor(Math.random() * 100) + 50;
        const startY = Math.floor(Math.random() * 100) + 50;
        const endX = Math.floor(Math.random() * 100) + 50;
        const endY = Math.floor(Math.random() * 100) + 50;

        // Randomly move the mouse across the page
        await humanizedMouseMove(startX, startY, endX, endY, 30);

        // Wait a random amount of time
        await page.waitForTimeout(getRandomDelay(1000, 3000));

        // Define the return position (close to top-left but avoiding the navbar)
        const returnX = Math.floor(Math.random() * 50) + 10; // Small x value close to left
        const returnY = Math.floor(Math.random() * 30) + 60; // Just below navbar (52px height)

        // Move back to the top-left, avoiding the navbar
        await humanizedMouseMove(endX, endY, returnX, returnY, 30);

        // Wait after the final move
        await page.waitForTimeout(getRandomDelay(1000, 3000));
      };

      // Loop until we have collected all connections or no more results can be loaded
      while (hasMoreResults && results.length < totalConnections) {
        await closeChatBubbles();

        const upDistance = await getRandomScrollHeight();
        const downDistance = await getRandomScrollHeight();

        await smoothScroll(downDistance, "down");
        await page.waitForTimeout(getRandomDelay(1506, 5210));
        await smoothScroll(upDistance, "up");
        if (Math.random() > 0.8) {
          await smoothScroll(upDistance, "up");
          await page.waitForTimeout(getRandomDelay(1006, 2010));
        }
        if (Math.random() > 0.5) {
          await smoothScroll(downDistance, "up");
          if (Math.random() > 0.8) {
            await moveMouseNaturally();
          } else {
            await page.waitForTimeout(getRandomDelay(1506, 5210));
          }
        }
        await page.waitForTimeout(getRandomDelay(2506, 6210));
        logger.info("Evaluate total scroll height", {
          source: "LinkedinAutomationService",
        });
        const totalScrollHeight = await page.evaluate(
          () => document.body.scrollHeight,
        );
        logger.info("Scroll to bottom", {
          source: "LinkedinAutomationService",
        });
        await smoothScroll(totalScrollHeight, "down");
        await page.waitForTimeout(getRandomDelay(1506, 5431));

        // Randomly simulate user interaction, such as a click or mouse movement (not too frequently)
        if (Math.random() > 0.8) {
          await moveMouseNaturally();
        }

        // Check if the "Show more results" button is visible, click if present
        const showMoreResultsButton = page.locator(
          "button.scaffold-finite-scroll__load-button",
        );

        await showMoreResultsButton.waitFor();
        await showMoreResultsButton.scrollIntoViewIfNeeded();

        if (await showMoreResultsButton.isVisible()) {
          logger.info("Click on 'Show more results' button", {
            source: "LinkedinAutomationService",
          });
          await showMoreResultsButton.click();
          await page.waitForTimeout(getRandomDelay(2120, 5300)); // Give time for new results to load
        }

        // Get newly loaded profile links
        logger.info("Evaluate all connection links", {
          source: "LinkedinAutomationService",
        });

        const newProfileUrls = await page
          .locator("a.mn-connection-card__link")
          .evaluateAll((links) => {
            return links.map((link, idx, arr) => {
              const url = link.getAttribute("href")?.split("?")?.[0] ?? "";

              if (idx !== 0 || idx !== arr.length - 1) {
                link.parentElement?.parentElement?.parentElement?.parentElement?.remove();
              }

              return `https://www.linkedin.com${url}`;
            });
          });

        // Add new profile URLs to the list, avoiding duplicates
        results.push(...newProfileUrls.filter((url) => !results.includes(url)));

        logger.info(
          `Collected ${results.length} connections out of ${totalConnections}`,
          {
            source: "LinkedinAutomationService",
          },
        );

        // If we've collected all profiles, stop scrolling
        if (results.length >= 100) {
          logger.info("No more results", {
            source: "LinkedinAutomationService",
          });
          hasMoreResults = false;
          break;
        }

        // Alternatively, check if we can scroll more, if not, stop.
        // logger.info("Evaluate scroll height", {
        //   source: "LinkedinAutomationService",
        // });
        // const currentScrollHeight = await page.evaluate(
        //   () => document.body.scrollHeight
        // );
        // if (currentScrollHeight === lastScrollHeight) {
        //   logger.info("No more results", {
        //     source: "LinkedinAutomationService",
        //   });
        //   hasMoreResults = false;
        // } else {
        //   lastScrollHeight = currentScrollHeight;
        // }
      }
    } catch (err) {
      error = LinkedinAutomationService.handleError(err);
    } finally {
      await page.close();
      return [results, error];
    }
  }

  async downloadAllConnections() {
    const browser = await Browser.getFreshInstance(this.proxyConfig, {
      debug: true,
    });
    const context = await browser.newContext({
      userAgent: this.userAgent,
    });

    await context.addCookies(this.cookies);
    const page = await context.newPage();

    const checkDownloadButton = async (locator: FrameLocator) => {
      try {
        const downloadButton = locator.locator("button.download-btn");
        await downloadButton.waitFor();

        const textContent = await downloadButton.textContent();

        if (textContent?.includes("Download")) {
          return downloadButton;
        }

        if (await downloadButton.isDisabled()) {
          logger.info("Download button was found but it's disabled.", {
            source: "LinkedinAutomationService",
          });

          return "disabled";
        }

        logger.info("Download button was found.", {
          source: "LinkedinAutomationService",
        });

        return downloadButton;
      } catch (error) {
        logger.info("Download button was not found.", {
          source: "LinkedinAutomationService",
        });
        return null;
      }
    };

    try {
      await page.goto(
        "https://www.linkedin.com/mypreferences/d/download-my-data",
        { timeout: 60 * 1000 },
      );

      const iframe = page.frameLocator(".settings-iframe--frame");

      let downloadButton = await checkDownloadButton(iframe);

      if (!downloadButton) {
        const fastFileLabel = iframe.locator('label[for="fast-file-only"]', {
          hasText:
            "Want something in particular? Select the data files you're most interested in.",
        });
        await fastFileLabel.waitFor();
        await fastFileLabel.click();

        page.waitForTimeout(1000);

        const connectionsLabel = iframe.locator(
          'label[for="file_group_CONNECTIONS"]',
        );
        await connectionsLabel.click();

        page.waitForTimeout(1000);

        const requestArchiveButton = iframe.locator("button#download-button");
        await requestArchiveButton.waitFor();
        await requestArchiveButton.click();
      }

      while (!downloadButton || downloadButton === "disabled") {
        await page.reload();
        logger.info("Reloading the page to check if download can start...", {
          source: "LinkedinAutomationService",
        });
        await page.waitForTimeout(180 * 1000);
        const _iframe = page.frameLocator(".settings-iframe--frame");

        downloadButton = await checkDownloadButton(_iframe);
      }

      if (downloadButton && typeof downloadButton !== "string") {
        downloadButton?.click();
      }
      const download = await page.waitForEvent("download");
      const downloadStream = await download.createReadStream();

      let zipBuffer = Buffer.alloc(0);

      await new Promise((resolve, reject) => {
        const writable = new Writable({
          write(chunk, _encoding, callback) {
            zipBuffer = Buffer.concat([zipBuffer, chunk]);
            callback();
          },
        });

        downloadStream.pipe(writable);
        downloadStream.on("end", resolve);
        downloadStream.on("error", reject);
      });

      const zip = new AdmZip(zipBuffer);
      const zipEntries = zip.getEntries();

      let csvFileContent: string[] = [];

      for (const entry of zipEntries) {
        if (entry.entryName.endsWith(".csv")) {
          csvFileContent = entry
            .getData()
            .toString("utf8")
            .split("\n")
            .slice(4)
            .map((row) => {
              const cells = row.split(",");
              const urlIndex = cells.findIndex((v) => v.startsWith("http"));

              return cells?.[urlIndex];
            })
            .filter(Boolean);
          break;
        }
      }

      logger.info("CSV file successfully processed", {
        source: "LinkedinAutomationService",
      });

      return csvFileContent;
    } catch (error) {
      LinkedinAutomationService.handleError(error);
    } finally {
      page.close();
    }
  }

  async getCompanyPeople(companyName: string) {
    const browser = await Browser.getFreshInstance(this.proxyConfig);
    const context = await browser.newContext({
      userAgent: this.userAgent,
    });

    await context.addCookies(this.cookies);
    const page = await context.newPage();

    try {
      await page.goto(
        `https://www.linkedin.com/company/${companyName}/people/`,
        { timeout: 60 * 1000 },
      );

      const totalContactsText = await page
        .locator("h2.text-heading-xlarge")
        .first()
        .innerText();

      const totalContacts = parseInt(totalContactsText.replace(/\D/g, ""), 10);

      const profileUrls: string[] = [];
      let hasMoreResults = true;
      let lastScrollHeight = 0;

      // Loop until we have collected all contacts or no more results can be loaded
      while (hasMoreResults && profileUrls.length < totalContacts) {
        // Scroll to the bottom of the page to load more results
        await page.evaluate(() => {
          window.scrollBy(0, window.innerHeight);
        });

        // Wait for a short period to allow more profiles to load
        await page.waitForTimeout(5000);

        // Check if the "Show more results" button is visible, click if present
        const showMoreResultsButton = page.locator(
          'button:has-text("Show more results")',
        );

        if (await showMoreResultsButton.isVisible()) {
          console.log("am gasit aici butonul.");

          await showMoreResultsButton.click();
          await page.waitForTimeout(5000); // Give time for new results to load
        }

        // Get newly loaded profile links
        const newProfileUrls = await page
          .locator('a.app-aware-link[aria-label*="View"][href*="/in/"]')
          .evaluateAll((links) => {
            return links.map(
              (link) => link.getAttribute("href")?.split("?")?.[0] ?? "",
            );
          });

        // Add new profile URLs to the list, avoiding duplicates
        profileUrls.push(
          ...newProfileUrls.filter((url) => !profileUrls.includes(url)),
        );

        // If we've collected all profiles, stop scrolling
        if (profileUrls.length >= totalContacts) {
          hasMoreResults = false;
          break;
        }

        // Alternatively, check if we can scroll more, if not, stop.
        const currentScrollHeight = await page.evaluate(
          () => document.body.scrollHeight,
        );
        if (currentScrollHeight === lastScrollHeight) {
          hasMoreResults = false;
        } else {
          lastScrollHeight = currentScrollHeight;
        }
      }

      return profileUrls;
    } catch (err) {
      LinkedinAutomationService.handleError(err);
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
  delay: number = 3000,
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
