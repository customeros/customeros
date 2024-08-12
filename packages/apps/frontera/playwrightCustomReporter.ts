import * as fs from 'fs';
import * as path from 'path';
import {
  Reporter,
  TestCase,
  TestResult,
  FullResult,
} from '@playwright/test/reporter';

class PlaywrightCustomReporter implements Reporter {
  private results: { name: string; status: string }[] = [];
  private failedTests: string[] = [];

  onTestEnd(test: TestCase, result: TestResult) {
    const emoji = result.status === 'passed' ? '✅' : '❌';
    const statusLine = `${test.title}: ${emoji} ${result.status}`;

    this.results.push({ name: test.title, status: statusLine });

    if (result.status === 'failed' || result.status === 'timedOut') {
      this.failedTests.push(`❌ failed: ${test.title}`);
    }
  }

  onEnd(result: FullResult) {
    const outputDir = process.env.PLAYWRIGHT_OUTPUT_DIR || '.';

    const output = this.results.map((r) => r.status).join('\n');

    fs.writeFileSync(path.join(outputDir, 'test-results.txt'), output);

    const failedOutput = this.failedTests.join('\n');

    fs.writeFileSync(path.join(outputDir, 'failed-tests.txt'), failedOutput);
  }
}

export default PlaywrightCustomReporter;
