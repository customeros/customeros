import type { Browser } from "./browser";

export class Scraper {
  constructor(public browser: Browser) {}

  public async scrape(url: string) {
    const page = await this.browser.getPage();

    await page.goto(url);

    const title = await page.title();
    console.log(title);
  }
}
