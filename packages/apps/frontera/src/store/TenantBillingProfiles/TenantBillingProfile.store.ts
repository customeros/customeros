import type { RootStore } from '@store/root';

import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import { DataSource, TenantBillingProfile } from '@graphql/types';

export class TenantBillingProfileStore implements Store<TenantBillingProfile> {
  value: TenantBillingProfile = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<TenantBillingProfile>();
  update = makeAutoSyncable.update<TenantBillingProfile>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'TenantBillingProfile',
      mutator: this.save,
      getId: (d) => d?.id,
    });
  }

  get id() {
    return this.value?.id;
  }

  set id(id: string) {
    this.value.id = id;
  }

  private async save() {
    // TODO
  }

  invalidate() {
    // todo
    return Promise.resolve();
  }

  init(data: TenantBillingProfile) {
    return merge(this.value, data);
  }
}

const defaultValue: TenantBillingProfile = {
  id: crypto.randomUUID(),
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
  appSource: DataSource.Openline,
  source: DataSource.Openline,
  sourceOfTruth: DataSource.Openline,

  phone: '',
  addressLine1: '',
  addressLine2: '',
  addressLine3: '',
  locality: '',
  country: '',
  region: '',
  zip: '',
  legalName: '',
  vatNumber: '',
  sendInvoicesFrom: '',
  sendInvoicesBcc: '',
  canPayWithBankTransfer: false,
  canPayWithPigeon: false,
  check: false,

  //deprecated
  email: '',
  domesticPaymentsBankInfo: '',
  internationalPaymentsBankInfo: '',
  canPayWithCard: false,
  canPayWithDirectDebitSEPA: false,
  canPayWithDirectDebitACH: false,
  canPayWithDirectDebitBacs: false,
};
