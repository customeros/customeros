import { Page, expect, Locator } from '@playwright/test';

export async function assertWithRetry(
  assertionFunc: () => Promise<void>,
  maxRetries = 10,
  retryInterval = 1000,
): Promise<void> {
  let lastError;

  for (let i = 0; i < maxRetries; i++) {
    try {
      await assertionFunc();

      return;
    } catch (error) {
      lastError = error;

      if (i < maxRetries - 1) {
        console.warn(`Assertion failed, retrying in ${retryInterval}ms...`);
        await new Promise((resolve) => setTimeout(resolve, retryInterval));
      }
    }
  }
  throw lastError;
}

export async function retryOperation(
  page: Page,
  operation: () => Promise<void>,
  maxAttempts: number,
  retryInterval: number,
) {
  for (let attempt = 0; attempt < maxAttempts; attempt++) {
    try {
      await operation();
      break; // Success, exit the loop
    } catch (error) {
      if (attempt === maxAttempts - 1) {
        throw error; // If it's the last attempt, throw the error
      }

      console.error(
        `Operation failed. Retrying in ${
          retryInterval / 1000
        } seconds... (Attempt ${attempt + 1}/${maxAttempts})`,
      );
      await page.waitForTimeout(retryInterval);
      await page.reload();
    }
  }
}

export async function ensureLocatorIsVisible(
  page: Page,
  selector: string,
): Promise<Locator> {
  const locator = page.locator(selector);

  await expect(locator).toBeVisible({ timeout: 15000 });

  return locator;
}

export async function ensureLocatorIsVisibleWithIndex(
  page: Page,
  selector: string,
  contractIndex: number,
): Promise<Locator> {
  const locator = page.locator(selector).nth(contractIndex);

  await expect(locator).toBeVisible();

  return locator;
}

export async function ensureLocatorIsVisibleAndHasText(
  page: Page,
  selector: string,
  expectedText: string,
): Promise<Locator> {
  const locator = await ensureLocatorIsVisible(page, selector);

  await expect(locator).toHaveText(expectedText);

  return locator;
}

export async function clickLocatorsThatAreVisible(
  page: Page,
  ...selectors: string[]
) {
  for (const selector of selectors) {
    await page.locator(selector).scrollIntoViewIfNeeded({ timeout: 30000 });
    await ensureLocatorIsVisible(page, selector);

    await page.click(selector);
  }
}

export async function clickLocatorThatIsVisible(page: Page, selector: string) {
  const locator = await ensureLocatorIsVisible(page, selector);

  await page.click(selector);

  return locator;
}

export async function doubleClickLocatorThatIsVisible(
  page: Page,
  selector: string,
) {
  const locator = await ensureLocatorIsVisible(page, selector);

  await page.dblclick(selector);

  return locator;
}

export async function writeTextInLocator(
  page: Page,
  selector: string,
  text: string,
) {
  const locator = await ensureLocatorIsVisible(page, selector);

  await page.click(selector);
  await locator.pressSequentially(text);

  return page;
}

export async function clickLocatorThatIsVisibleWithIndex(
  page: Page,
  selector: string,
  contractIndex: number,
) {
  await ensureLocatorIsVisibleWithIndex(page, selector, contractIndex);

  await page.click(selector);
}

export async function clickLocatorWithSublocatorThatIsVisible(
  page: Page,
  selector: string,
  subSelector: string,
): Promise<Locator> {
  const locator = await ensureLocatorIsVisible(page, selector);
  const subLocator = locator.locator(subSelector);

  await subLocator.click();

  return subLocator;
}

export async function clickLocatorThatIsVisibleAndHasText(
  page: Page,
  selector: string,
  expectedText: string,
) {
  await ensureLocatorIsVisibleAndHasText(page, selector, expectedText);

  await page.click(selector);
}

export function createRequestPromise(
  page: Page,
  requestsKey: string,
  requestValue: string | number,
) {
  return page.waitForRequest(
    (request) => {
      if (
        request.method() === 'POST' &&
        request.url().includes('customer-os-api')
      ) {
        const postData = request.postData();

        if (postData) {
          const parsedData = JSON.parse(postData);

          return parsedData.variables?.input?.[requestsKey] === requestValue;
        }
      }

      return false;
    },
    { timeout: 15000 },
  );
}

export function createResponsePromise(
  page: Page,
  responseKey: string,
  responseValue: string | undefined,
) {
  return page.waitForResponse(async (response) => {
    if (
      response.request().method() === 'POST' &&
      response.url().includes('customer-os-api')
    ) {
      const responseBody = await response.json();

      const getNestedProperty = (obj: unknown, path: string): unknown => {
        return path.split('.').reduce((prev, curr) => {
          if (curr.endsWith('?')) {
            curr = curr.slice(0, -1);
          }

          if (prev && typeof prev === 'object' && curr in prev) {
            return (prev as Record<string, unknown>)[curr];
          } else {
            return undefined;
          }
        }, obj);
      };

      const actualValue = getNestedProperty(responseBody.data, responseKey);

      return actualValue !== responseValue;
    }

    return false;
  });
}

export async function doScreenshot(page: Page, screenshotName: string) {
  await page.screenshot({
    path: screenshotName + '.png',
    fullPage: true,
  });
}
