import type { RootStore } from '@store/root.ts';

import set from 'lodash/set';
import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { Operation } from '@store/types.ts';
import { makePayload } from '@store/util.ts';
import { Transport } from '@store/transport.ts';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store.ts';

import {
  Contract,
  Currency,
  DataSource,
  ContractStatus,
  ContractUpdateInput,
  ContractRenewalCycle,
  ContractRenewalInput,
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
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'Contract',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
  }

  get id() {
    return this.value.metadata.id;
  }
  set id(id: string) {
    this.value.metadata.id = id;
  }
  get invoices() {
    return this.root.invoices
      .toArray()
      .filter((invoice) => invoice.value.contract.metadata.id === this.id);
  }

  get upcomingInvoices() {
    return this.root.invoices
      .toArray()
      .filter((invoice) => invoice.value.contract.metadata.id === this.id);
  }

  async invalidate() {
    try {
      this.isLoading = true;
      const { contract } = await this.transport.graphql.request<
        CONTRACT_QUERY_RESULT,
        { id: string }
      >(CONTRACT_QUERY, { id: this.id });
      this.load(contract);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async updateContract(payload: ContractUpdateInput) {
    try {
      this.isLoading = true;

      await this.transport.graphql.request<unknown, CONTRACT_UPDATE_PAYLOAD>(
        UPDATE_CONTRACT_DEF,
        {
          input: {
            ...payload,
            contractId: this.id,
            patch: true,
          },
        },
      );
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
        setTimeout(() => {
          this.invalidate();
        }, 600);
      });
    }
  }
  private async updateContractRenewalDate(payload: ContractRenewalInput) {
    try {
      this.isLoading = true;

      await this.transport.graphql.request<unknown, CONTRACT_RENEW_PAYLOAD>(
        RENEW_CONTRACT,
        {
          input: {
            ...payload,
            contractId: this.id,
          },
        },
      );
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const path = diff?.path;
    if (this.history.every((e) => !e.diff?.length)) {
      return;
    }
    if (!path) {
      return;
    }

    match(path)
      .with(['renewalDate', ...P.array()], () => {
        const payload = makePayload<ContractRenewalInput>(operation);
        this.updateContractRenewalDate(payload);
      })
      .with(['contractStatus', ...P.array()], () => {
        const { contractStatus, ...payload } = makePayload<
          ContractUpdateInput & { contractStatus: ContractStatus }
        >(operation);
        this.updateContract(payload);
      })

      .otherwise(() => {
        const payload = makePayload<ContractUpdateInput>(operation);
        this.updateContract(payload);
      });
  }

  async updateContractValues() {
    try {
      this.isLoading = true;

      await this.transport.graphql.request<unknown, CONTRACT_UPDATE_PAYLOAD>(
        UPDATE_CONTRACT_DEF,
        {
          input: {
            committedPeriodInMonths: this.value?.committedPeriodInMonths,
            serviceStarted: this.value?.serviceStarted,
            autoRenew: this.value.autoRenew,
            currency: this.value.currency,
            billingDetails: {
              invoicingStarted: this.value?.billingDetails?.invoicingStarted,
              billingCycleInMonths:
                this.value?.billingDetails?.billingCycleInMonths,
              dueDays: this.value?.billingDetails?.dueDays,
              payAutomatically: this.value?.billingDetails?.payAutomatically,
              canPayWithCard: this.value?.billingDetails?.canPayWithCard,
              canPayWithDirectDebit:
                this.value?.billingDetails?.canPayWithDirectDebit,
              payOnline: this.value?.billingDetails?.payOnline,
              canPayWithBankTransfer:
                this.value?.billingDetails?.canPayWithBankTransfer,
              check: this.value?.billingDetails?.check,
            },
            billingEnabled: this.value.billingEnabled,

            contractId: this.id,
            patch: true,
          },
        },
      );
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async updateBillingAddress() {
    try {
      this.isLoading = true;

      await this.transport.graphql.request<unknown, CONTRACT_UPDATE_PAYLOAD>(
        UPDATE_CONTRACT_DEF,
        {
          input: {
            billingDetails: {
              organizationLegalName:
                this.value?.billingDetails?.organizationLegalName,
              country: this.value?.billingDetails?.country,
              addressLine1: this.value?.billingDetails?.addressLine1,
              addressLine2: this.value?.billingDetails?.addressLine2,
              locality: this.value?.billingDetails?.locality,
              postalCode: this.value?.billingDetails?.postalCode,
              region: this.value?.billingDetails?.region,
              canPayWithBankTransfer:
                this.value?.billingDetails?.canPayWithBankTransfer,
              billingEmail: this.value?.billingDetails?.billingEmail,
              billingEmailCC: this.value?.billingDetails?.billingEmailCC,
              billingEmailBCC: this.value?.billingDetails?.billingEmailBCC,
            },
            contractId: this.id,
            patch: true,
          },
        },
      );
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  init(data: Contract) {
    const output = merge(this.value, data);
    const contractLineItems = data.contractLineItems?.map((item) => {
      this.root.contractLineItems.load([item]);

      return this.root.contractLineItems.value.get(item?.metadata?.id)?.value;
    });
    const opportunities = data.opportunities?.map((item) => {
      this.root.opportunities.load([item]);

      return this.root.opportunities.value.get(item?.metadata?.id)?.value;
    });

    contractLineItems && set(output, 'contractLineItems', contractLineItems);
    opportunities && set(output, 'opportunities', opportunities);

    return output;
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
    id: crypto.randomUUID(),
    appSource: DataSource.Openline,
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
  },
  upcomingInvoices: [],
  attachments: [],
  billingDetails: {
    billingCycleInMonths: 1,
    invoicingStarted: new Date().toISOString(),
    nextInvoicing: new Date().toISOString(),
    addressLine1: '',
    addressLine2: '',
    locality: '',
    region: '',
    country: '',
    postalCode: '',
    organizationLegalName: '',
    billingEmail: '',
    billingEmailCC: [],
    billingEmailBCC: [],
    invoiceNote: '',
    canPayWithCard: false,
    canPayWithDirectDebit: false,
    canPayWithBankTransfer: false,
    payOnline: false,
    payAutomatically: false,
    check: false,
    dueDays: 30,
  },
  committedPeriodInMonths: 1,
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
  id: crypto.randomUUID(),
  name: '',
  renewalCycle: ContractRenewalCycle.None,
  source: DataSource.Openline,
  sourceOfTruth: DataSource.Openline,
  status: ContractStatus.Undefined,
  updatedAt: '',
};

type CONTRACT_UPDATE_PAYLOAD = { input: ContractUpdateInput };
const UPDATE_CONTRACT_DEF = gql`
  mutation updateContract($input: ContractUpdateInput!) {
    contract_Update(input: $input) {
      id
    }
  }
`;

type CONTRACT_RENEW_PAYLOAD = { input: ContractUpdateInput };
const RENEW_CONTRACT = gql`
  mutation renewContract($input: ContractRenewalInput!) {
    contract_Renew(input: $input) {
      id
    }
  }
`;

type CONTRACT_QUERY_RESULT = {
  contract: Contract;
};
const CONTRACT_QUERY = gql`
  query Contract($id: ID!) {
    contract(id: $id) {
      metadata {
        id
        created
        source
        lastUpdated
      }

      contractName
      serviceStarted
      contractSigned
      contractEnded
      contractStatus
      committedPeriodInMonths
      approved

      contractUrl
      billingCycle
      billingEnabled
      currency
      invoiceEmail
      autoRenew

      billingDetails {
        nextInvoicing
        billingCycleInMonths
        postalCode
        country
        locality
        addressLine1
        addressLine2
        invoiceNote
        organizationLegalName
        canPayWithBankTransfer
        billingCycle
        invoicingStarted
        payAutomatically
        region
        dueDays
        canPayWithCard
        billingEmail
        billingEmailCC
        billingEmailBCC
      }
      upcomingInvoices {
        metadata {
          id
        }
      }
      opportunities {
        metadata {
          id
          created
          lastUpdated
          source
          sourceOfTruth
          appSource
        }
        name
        amount
        maxAmount
        internalType
        externalType
        internalStage
        externalStage
        estimatedClosedAt
        generalNotes
        nextSteps
        renewedAt
        renewalApproved
        renewalLikelihood
        renewalUpdatedByUserId
        renewalUpdatedByUserAt
        renewalAdjustedRate
        comments
        organization {
          metadata {
            id
            created
            lastUpdated
            sourceOfTruth
          }
        }
        createdBy {
          id
          firstName
          lastName
          name
        }
        owner {
          id
          firstName
          lastName
          name
        }
        externalLinks {
          externalUrl
          externalId
        }
        id
        createdAt
        updatedAt
        source
        appSource
      }
      contractLineItems {
        metadata {
          id
          created
          lastUpdated
          source
          appSource
          sourceOfTruth
        }
        description
        billingCycle
        price
        quantity
        comments
        serviceEnded
        parentId
        serviceStarted
        tax {
          salesTax
          vat
          taxRate
        }
      }
    }
  }
`;
