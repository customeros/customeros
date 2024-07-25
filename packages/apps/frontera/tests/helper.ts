// helper.ts

export async function assertWithRetry(
  assertionFunc: () => Promise<void>,
  maxRetries = 5,
  retryInterval = 3000,
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
      await this.page.waitForTimeout(retryInterval);
      await this.page.reload();
    }
  }
}
