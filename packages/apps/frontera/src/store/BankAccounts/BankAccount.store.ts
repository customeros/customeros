import type { RootStore } from '@store/root';

import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { makePayload } from '@store/util.ts';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import {
  DataSource,
  BankAccount,
  BankAccountUpdateInput,
} from '@graphql/types';

import { BankAccountService } from './BankAccount.service';

export class BankAccountStore implements Store<BankAccount> {
  value: BankAccount = getDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<BankAccount>();
  update = makeAutoSyncable.update<BankAccount>();
  private service: BankAccountService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'BankAccount',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
    this.service = BankAccountService.getInstance(transport);
  }

  get id() {
    return this.value.metadata?.id;
  }

  setId(id: string) {
    this.value.metadata.id = id;
  }

  private async save(operation: Operation) {
    const payload = makePayload<BankAccountUpdateInput>(operation);

    await this.service.updateBankAccount({
      input: { ...payload, id: this.value.metadata.id },
    });
  }

  invalidate() {
    // todo
    return Promise.resolve();
  }

  init(data: BankAccount) {
    return merge(this.value, data);
  }
}

const getDefaultValue = (): BankAccount => ({
  accountNumber: null,
  allowInternational: false,
  bankName: '',
  bankTransferEnabled: false,
  bic: null,
  currency: null,
  iban: null,
  metadata: {
    source: DataSource.Openline,
    appSource: DataSource.Openline,
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    sourceOfTruth: DataSource.Openline,
  },
  otherDetails: null,
  routingNumber: null,
  sortCode: null,
});
