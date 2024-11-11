import localforage from 'localforage';
import { computed, reaction, observable, makeObservable } from 'mobx';

import type { RootStore } from './root';
import type { Transport } from './transport';
import type { Operation, GroupOperation } from './types';

import { GraphqlService } from './graphql';

type TransactionType = 'SINGLE' | 'GROUP';
type TransactionStatus =
  | 'PENDING'
  | 'PROCESSING'
  | 'COMPLETED'
  | 'RETRYING'
  | 'FAILED';
type OperationType = Operation | GroupOperation;
type TransactionOptions = {
  syncOnly?: boolean;
  persist?: () => void;
  onFailled?: () => void;
  onCompleted?: () => void;
};

class Transaction {
  type: TransactionType;
  status: TransactionStatus = 'PENDING';
  createdAt = new Date();
  completedAt?: Date;
  failedAt?: Date;
  retryDelay = 0;
  lastAttemptAt: Date | null = null;
  operation: OperationType;
  syncOnly: boolean = false;
  retryAttempts = 0;

  constructor(
    type: TransactionType,
    operation: OperationType,
    private options?: TransactionOptions,
  ) {
    this.type = type;
    this.operation = operation;
    this.syncOnly = this?.options?.syncOnly ?? false;
  }

  start() {
    this.status = 'PROCESSING';
  }

  complete() {
    this.status = 'COMPLETED';
    this.completedAt = new Date();
    this?.options?.onCompleted?.();
    this?.options?.persist?.();
  }

  retry() {
    this.status = 'RETRYING';
    this.retryAttempts++;
    this.lastAttemptAt = new Date();
    this.retryDelay = this.calculateBackoff();
  }

  fail() {
    this.status = 'FAILED';
    this.failedAt = new Date();
    this?.options?.onFailled?.();
  }

  private calculateBackoff(): number {
    return 1000 * Math.pow(2, this.retryAttempts);
  }
}

class TransactionQueue {
  private id: string;
  private queue: Transaction[] = [];

  get hasCommits() {
    return this.queue.length > 0;
  }

  get storageKey() {
    return `tx-queue-${this.id}`;
  }

  constructor(id: string) {
    this.id = id;

    makeObservable<TransactionQueue, 'queue'>(this, {
      queue: observable,
      hasCommits: computed,
    });

    window.addEventListener('online', async () => {
      const savedTxs = await this.loadSaved();

      if (!savedTxs) return;

      savedTxs.forEach((tx) => this.push(tx));

      await this.clearSaved();
    });
  }

  push(tx: Transaction) {
    if (navigator?.onLine) {
      this.queue.push(tx);
    } else {
      this.save(tx);
    }
  }

  next() {
    return this.queue.shift();
  }

  private async save(tx: Transaction) {
    try {
      const prev = await localforage.getItem<Transaction[]>(this.storageKey);

      await localforage.setItem(this.storageKey, [...(prev ?? []), tx]);
    } catch (err) {
      console.error('Could not save transaction.', err);
    }
  }

  private async loadSaved() {
    try {
      return await localforage.getItem<Transaction[]>(this.storageKey);
    } catch (err) {
      console.error('Could not load saved transactions.', err);
    }
  }

  private async clearSaved() {
    try {
      await localforage.removeItem(this.storageKey);
    } catch (err) {
      console.error('Could not remove saved transactions', err);
    }
  }
}

class TransactionRunner {
  private isMainRunning = false;
  private isRetryRunning = false;
  private maxRetries = 3;
  private graphqlService: GraphqlService;

  constructor(
    private root: RootStore,
    private transport: Transport,
    private mainQueue: TransactionQueue,
    private retryQueue: TransactionQueue,
  ) {
    this.graphqlService = new GraphqlService(this.root, this.transport);
  }

  // tx should contain a type property to handle group packets too
  async processMainQueue() {
    if (this.isMainRunning) return;
    this.isMainRunning = true;

    while (this.mainQueue.hasCommits) {
      const tx = this.mainQueue.next();

      if (!tx) continue;
      if (!tx.operation.entityId) continue;

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
      if (!tx.operation.entityId) continue;

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
    tx.start();

    if (!tx.syncOnly) {
      await this.processGraphqlMutation(tx);
    }
    await this.processSyncPacket(tx);
    tx.complete();
  }

  private async processSyncPacket(tx: Transaction) {
    return await new Promise<void>((resolve, reject) => {
      const channelBinding =
        tx.type === 'SINGLE' ? 'sync_packet' : 'sync_group_packet';
      const channelKey = tx.operation?.entity;
      const channel = this.transport?.channels?.get(channelKey as string);

      channel
        ?.push(channelBinding, { payload: { operation: tx.operation } })
        ?.receive('ok', () => {
          resolve();
        })
        ?.receive('error', () => {
          reject();
        });
    });
  }

  private async processGraphqlMutation(tx: Transaction) {
    return await this.graphqlService.mutate(tx.operation as Operation);
  }

  private handleRetry(tx: Transaction) {
    if (tx.retryAttempts < this.maxRetries) {
      tx.retry();
      this.retryQueue.push(tx);
    } else {
      tx.fail();
    }
  }

  private hasRetryDelayPassed(tx: Transaction): boolean {
    if (!tx.lastAttemptAt || !tx.retryDelay) return true;
    const now = new Date().getTime();
    const lastAttemptTime = tx.lastAttemptAt.getTime();

    return now - lastAttemptTime >= tx.retryDelay;
  }
}

export class TransactionService {
  private mainQueue: TransactionQueue;
  private retryQueue: TransactionQueue;
  private runner: TransactionRunner;

  constructor(private root: RootStore, private transport: Transport) {
    this.mainQueue = new TransactionQueue('main');
    this.retryQueue = new TransactionQueue('retry');
    this.runner = new TransactionRunner(
      this.root,
      this.transport,
      this.mainQueue,
      this.retryQueue,
    );
  }

  commit(operation: Operation, options?: TransactionOptions) {
    this.mainQueue.push(new Transaction('SINGLE', operation, options));
  }

  groupCommit(operation: GroupOperation, options?: TransactionOptions) {
    this.mainQueue.push(new Transaction('GROUP', operation, options));
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
