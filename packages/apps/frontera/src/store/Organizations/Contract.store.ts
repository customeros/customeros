import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import {
  Contract,
  Currency,
  DataSource,
  ContractStatus,
  ContractRenewalCycle,
} from '@graphql/types';

export class ContractStore implements Store<Contract> {
  value: Contract = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<Contract>();
  update = makeAutoSyncable.update<Contract>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'Contract',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
    makeAutoObservable(this);
  }

  async invalidate() {}

  set id(id: string) {
    this.value.metadata.id = id;
  }

  private async save() {}
}

const defaultValue: Contract = {
  approved: false,
  autoRenew: false,
  billingEnabled: false,
  contractName: '',
  contractStatus: ContractStatus.Undefined,
  contractUrl: '',
  externalLinks: [],
  invoices: [],
  metadata: {
    id: '',
    appSource: DataSource.Openline,
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
  },
  upcomingInvoices: [],
  attachments: [],
  billingDetails: {},
  committedPeriodInMonths: 0,
  contractEnded: '',
  contractLineItems: [],
  contractSigned: '',
  ltv: 0,
  serviceStarted: '',
  createdBy: null,
  currency: Currency.Usd,
  opportunities: [],
  owner: null,
  // deprecated fields -> should be removed when schema is updated
  appSource: DataSource.Openline,
  contractRenewalCycle: ContractRenewalCycle.None,
  createdAt: '',
  id: '',
  name: '',
  renewalCycle: ContractRenewalCycle.None,
  source: DataSource.Openline,
  sourceOfTruth: DataSource.Openline,
  status: ContractStatus.Undefined,
  updatedAt: '',
};
