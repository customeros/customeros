import { chromium, test as base, BrowserContext } from '@playwright/test';

type VideoFixture = {
  context: BrowserContext;
};

interface TestFixtures {}

export const test = base.extend<TestFixtures & VideoFixture>({
  context: async (
    { context: _ },
    use: (context: BrowserContext) => Promise<void>,
  ) => {
    const browser = await chromium.launch();
    const context = await browser.newContext({
      recordVideo: {
        dir: 'videos/',
        size: { width: 1280, height: 720 },
      },
    });

    await use(context);

    await context.close();
    await browser.close();
  },
  page: async ({ context }, use) => {
    const page = await context.newPage();

    await use(page);
  },
});

export { expect } from '@playwright/test';
