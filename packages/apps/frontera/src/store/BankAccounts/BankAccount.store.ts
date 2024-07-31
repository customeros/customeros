import type { RootStore } from '@store/root';

import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import { Currency, DataSource, BankAccount } from '@graphql/types';

export class BankAccountStore implements Store<BankAccount> {
  value: BankAccount = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<BankAccount>();
  update = makeAutoSyncable.update<BankAccount>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'BankAccount',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
  }

  get id() {
    return this.value.metadata?.id;
  }

  set id(id: string) {
    this.value.metadata.id = id;
  }

  private async save() {
    // TODO
  }

  invalidate() {
    // todo
    return Promise.resolve();
  }

  init(data: BankAccount) {
    return merge(this.value, data);
  }
}

const defaultValue: BankAccount = {
  metadata: {
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    appSource: DataSource.Openline,
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
  },
  bankName: '',
  currency: Currency.Usd,
  bankTransferEnabled: false,
  allowInternational: false,
  iban: '',
  bic: '',
  sortCode: '',
  accountNumber: '',
  routingNumber: '',
  otherDetails: '',
};
