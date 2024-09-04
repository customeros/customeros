import { chromium } from "playwright";
import type {
  Page,
  BrowserContextOptions,
  Browser as BrowserType,
} from "playwright";

const bcatUrl = "wss://api.browsercat.com/connect";

export class Browser {
  private static instance: Browser;
  public browser: BrowserType | null = null;

  private constructor(private debug?: boolean) {}

  public static async getInstance(option?: {
    debug?: boolean;
  }): Promise<Browser> {
    if (!Browser.instance) {
      Browser.instance = new Browser(option?.debug);
      await Browser.instance.init();
    }
    return Browser.instance;
  }

  private async init() {
    if (!this.browser) {
      if (this.debug) {
        this.browser = await chromium.launch({
          headless: false,
        });

        return;
      }

      this.browser = await chromium.connect(bcatUrl, {
        headers: {
          "api-key":
            "mECU022AvJTbliSlmmx2QeOWFuYnLDpFQx3V1bAkxU2MqIbcvko26C5bo1yH3iNM",
          "browsercat-opts": JSON.stringify({
            proxy: {
              server: "168.158.96.227:10323",
              username: "14aa16eea9d0e",
              password: "edf2135afe",
            },
          }),
        },
      });
    }
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

  public newContext(options?: BrowserContextOptions) {
    if (!this.browser) {
      throw new Error("Browser is not initialized");
    }

    return this.browser.newContext(options);
  }
}
