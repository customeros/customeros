import { chromium } from "playwright";
import type {
  Page,
  BrowserContextOptions,
  Browser as BrowserType,
} from "playwright";
import { ErrorParser, StandardError } from "@/util/error";

import { logger } from "../logger";

const bcatUrl = process.env.BROWSERCAT_API_URL;
const apiKey = process.env.BROWSERCAT_API_KEY;

export class Browser {
  private static instances: Map<string, Browser>;
  public browser: BrowserType | null = null;

  private constructor(private debug?: boolean) {}

  public static async getInstance(
    proxyConfig: string,
    options?: {
      debug?: boolean;
    },
  ): Promise<Browser> {
    if (!Browser.instances) {
      Browser.instances = new Map();
    }

    if (!Browser.instances.has(proxyConfig)) {
      const instance = new Browser(options?.debug);
      Browser.instances.set(proxyConfig, instance);
      await instance.init(proxyConfig);
    }

    const instance = Browser.instances.get(proxyConfig);

    if (!instance) {
      throw new StandardError({
        code: "INTERNAL_ERROR",
        message: "Browser instance not found for the given proxy config",
        severity: "critical",
      });
    }

    logger.debug("returning browser instance", {
      source: "Browser",
    });
    return instance;
  }

  private async init(proxyConfig: string) {
    return new Promise<void>(async (resolve, reject) => {
      if (!this.browser) {
        try {
          if (this.debug) {
            this.browser = await chromium.launch({
              headless: false,
              logger: {
                isEnabled: (_name, severity) => true,
                log: (_name, _severity, message, _args) => {
                  if (message instanceof Error) {
                    return logger.error(message.message, {
                      source: "Playwright",
                    });
                  }

                  return logger.info(message, {
                    source: "Playwright",
                  });
                },
              },
            });
          } else {
            if (!apiKey || !bcatUrl) {
              throw new StandardError({
                code: "INTERNAL_ERROR",
                message: "Browsercat API key or url is not provided",
                severity: "critical",
              });
            }

            logger.info("Connecting to Browsercat", {
              source: "Browser",
            });
            const browser = await chromium.connect(bcatUrl, {
              headers: {
                "api-key": apiKey,
                "browsercat-opts": proxyConfig,
              },
            });
            this.browser = browser;
          }
          logger.info("Browser initialized successfully", {
            source: "Browser",
          });
          resolve();
        } catch (err) {
          const error = ErrorParser.parse(err);
          logger.error("Error in Browser", {
            error: error.message,
            details: error.details,
          });

          reject(
            new StandardError({
              code: "EXTERNAL_ERROR",
              message: `Failed to initialize the browser.`,
              details: error.details,
              severity: "critical",
            }),
          );
        }
      }
    });
  }

  public async getPage(): Promise<Page> {
    if (!this.browser) {
      throw new Error("Browser is not initialized");
    }
    return await this.browser.newPage();
  }

  public async close() {
    if (this.browser) {
      await this.browser.close();
      this.browser = null;
    }
  }

  public async newContext(options?: BrowserContextOptions) {
    if (!this.browser) {
      throw new Error("Browser is not initialized");
    }

    return await this.browser.newContext(options);
  }
}
