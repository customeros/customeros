import * as fs from 'fs';
import * as path from 'path';
import { chromium, test as base, BrowserContext } from '@playwright/test';

type VideoFixture = {
  context: BrowserContext;
};

interface TestFixtures {}

export const test = base.extend<TestFixtures & VideoFixture>({
  context: async ({ context: _ }, use) => {
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
  page: async ({ context }, use, testInfo) => {
    const page = await context.newPage();

    await use(page);

    await page.close();

    const video = page.video();

    if (video) {
      const videoPath = await video.path();
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
      const customVideoName = `${testInfo.title.replace(
        /\s+/g,
        '_',
      )}_${timestamp}.webm`;
      const newVideoPath = path.join('videos', customVideoName);

      testInfo.attachments.push({
        name: 'video',
        path: newVideoPath,
        contentType: 'video/webm',
      });

      await new Promise<void>((resolve) => {
        fs.rename(videoPath, newVideoPath, (err) => {
          if (err) console.error('Error renaming video:', err);
          resolve();
        });
      });
    }
  },
});

export { expect } from '@playwright/test';
