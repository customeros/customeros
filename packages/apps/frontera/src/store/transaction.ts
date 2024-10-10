import { computed, reaction, observable, makeObservable } from 'mobx';

import type { Operation } from './types';
import type { Transport } from './transport';

type Transaction = {
  status: string;
  createdAt: Date;
  retryDelay: number;
  channelName: string;
  lastAttemptAt: Date;
  operation: Operation;
  retryAttempts: number;
};

class TransactionQueue {
  private queue: Transaction[] = [];

  get hasCommits() {
    return this.queue.length > 0;
  }

  constructor() {
    makeObservable<TransactionQueue, 'queue'>(this, {
      queue: observable,
      hasCommits: computed,
    });
  }

  push(tx: Transaction) {
    this.queue.push(tx);
  }

  next() {
    return this.queue.shift();
  }
}

class TransactionRunner {
  private isMainRunning = false;
  private isRetryRunning = false;
  private maxRetries = 3;
  private baseDelay = 1000;

  constructor(
    private transport: Transport,
    private mainQueue: TransactionQueue,
    private retryQueue: TransactionQueue,
  ) {}

  // tx should contain a type property to handle group packets too
  async processMainQueue() {
    if (this.isMainRunning) return;
    this.isMainRunning = true;

    while (this.mainQueue.hasCommits) {
      const tx = this.mainQueue.next();

      if (!tx) continue;
      if (!tx.channelName) continue;

      try {
        this.processTransaction(tx);
      } catch (err) {
        this.handleRetry(tx);
      }
    }

    this.isMainRunning = false;
  }

  async processRetryQueue() {
    if (this.isRetryRunning) return;
    this.isRetryRunning = true;

    while (this.retryQueue.hasCommits) {
      const tx = this.retryQueue.next();

      if (!tx) continue;
      if (!tx.channelName) continue;

      const delayPassed = this.hasRetryDelayPassed(tx);

      if (!delayPassed) {
        this.retryQueue.push(tx);
        continue;
      }

      try {
        await this.processTransaction(tx);
      } catch (err) {
        this.handleRetry(tx);
      }
    }

    this.isRetryRunning = false;
  }

  private async processTransaction(tx: Transaction) {
    await new Promise((resolve, reject) => {
      this.transport.channels
        .get(tx.channelName)
        ?.push('sync_packet', { payload: { operation: tx.operation } })
        ?.receive('ok', resolve)
        ?.receive('error', reject);
    });
  }

  private handleRetry(tx: Transaction) {
    if (!tx.retryAttempts) {
      tx.retryAttempts = 0;
    }

    if (tx.retryAttempts < this.maxRetries) {
      tx.retryAttempts++;
      tx.lastAttemptAt = new Date();

      const delay = this.calculateBackoff(tx.retryAttempts);

      tx.retryDelay = delay;

      this.retryQueue.push(tx);
    } else {
      // log it or move to a dead letter queue
    }
  }

  private hasRetryDelayPassed(tx: Transaction): boolean {
    if (!tx.lastAttemptAt || !tx.retryDelay) return true;
    const now = new Date().getTime();
    const lastAttemptTime = tx.lastAttemptAt.getTime();

    return now - lastAttemptTime >= tx.retryDelay;
  }

  private calculateBackoff(attempt: number): number {
    return this.baseDelay * Math.pow(2, attempt);
  }
}

export class TransactionService {
  private mainQueue: TransactionQueue;
  private retryQueue: TransactionQueue;
  private runner: TransactionRunner;

  constructor(private transport: Transport) {
    this.mainQueue = new TransactionQueue();
    this.retryQueue = new TransactionQueue();
    this.runner = new TransactionRunner(
      this.transport,
      this.mainQueue,
      this.retryQueue,
    );
  }

  commit(tx: Transaction) {
    this.mainQueue.push(tx);
  }

  startRunners() {
    reaction(
      () => this.mainQueue.hasCommits,
      (hasCommits) => {
        if (!hasCommits) return;
        this.runner.processMainQueue();
      },
    );
    reaction(
      () => this.retryQueue.hasCommits,
      (hasCommits) => {
        if (!hasCommits) return;
        this.runner.processRetryQueue();
      },
    );
  }
}
