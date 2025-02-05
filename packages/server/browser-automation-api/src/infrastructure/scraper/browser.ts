import { chromium } from "playwright";
import type {
  Page,
  BrowserContextOptions,
  Browser as BrowserType,
} from "playwright";
import { Hyperbrowser } from "@hyperbrowser/sdk";
import { ErrorParser, StandardError } from "@/util/error";

import { logger } from "../logger";

export type ProxyConfig = {
  proxyServer: string;
  proxyServerUsername: string;
  proxyServerPassword: string;
};

const bcatUrl = process.env.HYPERBROWSER_API_URL;
const apiKey = process.env.HYPERBROWSER_API_KEY ?? "";

const hyperbrowser = new Hyperbrowser({
  apiKey,
});

export class Browser {
  private static instances: Map<string, Browser>;
  public browser: BrowserType | null = null;

  constructor(private debug?: boolean, private debugRemote?: boolean) {}

  public static async getFreshInstance(
    proxyConfig: ProxyConfig,
    options?: {
      debug?: boolean;
      debugRemote?: boolean;
    }
  ): Promise<Browser> {
    logger.info("Creating fresh browser instance.", {
      source: "Browser",
    });
    const instance = new Browser(options?.debug, options?.debugRemote);
    await instance.init(proxyConfig);
    logger.info("Fresh browser instance created ok.", {
      source: "Browser",
    });

    return instance;
  }

  private async init(proxyConfig: ProxyConfig) {
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
                message: "Remote Browser API key or url is not provided",
                severity: "critical",
              });
            }

            const session = await hyperbrowser.sessions.create({
              useStealth: true,
              ...proxyConfig,
            });

            if (!session.wsEndpoint) {
              throw new StandardError({
                code: "EXTERNAL_ERROR",
                message: `Failed to create a remote browser session.`,
                severity: "critical",
              });
            }

            logger.info("Connecting to remote browser", {
              source: "Browser",
            });

            const browser = await chromium.connectOverCDP(session.wsEndpoint, {
              logger: {
                isEnabled: (_name, severity) => !!this.debugRemote,
                log: (name, _severity, message, _args) => {
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
            })
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
