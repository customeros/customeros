import type { RootStore } from '@store/root';

import omit from 'lodash/omit';
import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import {
  Contract,
  Currency,
  DataSource,
  ContractStatus,
  ContractUpdateInput,
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
      storeMapper: {
        contractLineItems: {
          storeName: 'contractLineItems',
          getItemId: (data) => data?.metadata?.id as string,
        },
      },
    });
    makeAutoObservable(this);
  }

  async invalidate() {}

  set id(id: string) {
    this.value.metadata.id = id;
  }

  private async save() {
    const payload: PAYLOAD = {
      input: {
        ...omit(this.value, 'metadata', 'owner'),
        contractId: this.value.metadata.id,
      },
    };
    try {
      this.isLoading = true;
      await this.transport.graphql.request(UPDATE_CONTRACT_DEF, payload);
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  public async updateStatus(status: ContractStatus) {
    this.value.contractStatus = status;
    await this.save();
  }
  public async updateName(name: string) {
    this.value.contractName = name;
    await this.save();
  }
}

const defaultValue: Contract = {
  approved: false,
  autoRenew: false,
  billingEnabled: false,
  contractName: 'Unnamed Contract',
  contractStatus: ContractStatus.Draft,
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
  contractRenewalCycle: ContractRenewalCycle.MonthlyRenewal,
  createdAt: '',
  id: '',
  name: '',
  renewalCycle: ContractRenewalCycle.None,
  source: DataSource.Openline,
  sourceOfTruth: DataSource.Openline,
  status: ContractStatus.Undefined,
  updatedAt: '',
};

type PAYLOAD = { input: ContractUpdateInput };
const UPDATE_CONTRACT_DEF = gql`
  mutation updateContract($input: ContractUpdateInput!) {
    contract_Update(input: $input) {
      id
    }
  }
`;
