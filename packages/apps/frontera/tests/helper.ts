import { Page, expect, Locator, Response } from '@playwright/test';

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

  await locator.waitFor({ state: 'attached', timeout: 15000 });
  await expect(locator).toBeVisible({ timeout: 15000 });

  await expect(locator).toBeEnabled({ timeout: 5000 });

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
    await retryOperation(
      page,
      async () => {
        const locator = await ensureLocatorIsVisible(page, selector);

        await page.waitForTimeout(300);

        await locator.waitFor({ state: 'attached', timeout: 5000 });
        await locator.scrollIntoViewIfNeeded({ timeout: 30000 });

        await expect(locator).toBeVisible({ timeout: 5000 });

        try {
          await locator.click({ timeout: 5000 });
        } catch (clickError) {
          await page.waitForTimeout(500);
          await locator.click({ force: true, timeout: 5000 });
        }
      },
      3,
      1500,
    );
  }
}

export async function clickLocatorThatIsVisible(page: Page, selector: string) {
  const locator = await ensureLocatorIsVisible(page, selector);

  // Add stability delay after ensuring visibility
  await page.waitForTimeout(300);

  try {
    await locator.click({ timeout: 5000 });
  } catch (error) {
    await page.waitForTimeout(500);
    await locator.click({ force: true, timeout: 5000 });
  }

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
  expectedKey: string,
  expectedValue: string | number,
) {
  return new Promise((resolve, reject) => {
    const timeout = setTimeout(() => {
      reject(
        new Error(
          `The request "${expectedKey}" didn't get set with the "${expectedValue}" value before timeout`,
        ),
      );
    }, 30000);

    page
      .waitForRequest(
        (request) => {
          if (
            request.method() === 'POST' &&
            request.url().includes('customer-os-api')
          ) {
            const postData = request.postData();

            if (postData) {
              const parsedData = JSON.parse(postData);

              return (
                parsedData.variables?.input?.[expectedKey] === expectedValue
              );
            }
          }

          return false;
        },
        { timeout: 30000 },
      )
      .then(() => {
        clearTimeout(timeout);
        resolve(true);
      })
      .catch((error) => {
        clearTimeout(timeout);
        reject(error);
      });
  });
}

export function createResponsePromise(
  page: Page,
  expectedKey: string,
  expectedValue: string | undefined,
) {
  return new Promise<Response>((resolve, reject) => {
    const timeout = setTimeout(() => {
      reject(
        new Error(
          `The response "${expectedKey}" remained set with the "${expectedValue}" value before timeout`,
        ),
      );
    }, 30000);

    page
      .waitForResponse(async (response) => {
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

          const actualValue = getNestedProperty(responseBody.data, expectedKey);

          if (actualValue !== expectedValue) {
            resolve(response);

            return true;
          }

          return false;
        }

        return false;
      })
      .catch((error) => {
        clearTimeout(timeout);
        reject(error);
      });
  });
}

export async function doScreenshot(page: Page, screenshotName: string) {
  await page.screenshot({
    path: screenshotName + '.png',
    fullPage: true,
  });
}
