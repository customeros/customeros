import type { RootStore } from '@store/root';

import set from 'lodash/set';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makePayload } from '@store/util.ts';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
import { Store, makeAutoSyncable } from '@store/store';

import {
  Contract,
  Currency,
  DataSource,
  Organization,
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
    makeAutoSyncable(this, {
      channelName: 'Contract',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
    makeAutoObservable(this);
  }

  get id() {
    return this.value.metadata.id;
  }
  set id(id: string) {
    this.value.metadata.id = id;
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
    // const type = diff?.op;
    const path = diff?.path;
    // const value = diff?.val;
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
  init(data: Contract) {
    const contracts = data.contractLineItems?.map((item) => {
      this.root.contractLineItems.load([item]);

      return this.root.contractLineItems.value.get(item.metadata.id)?.value;
    });

    set(data, 'contractLineItems', contracts);

    return data;
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
        postalCode
        country
        locality
        addressLine1
        addressLine2
        invoiceNote
        organizationLegalName
        billingCycle
        invoicingStarted
        region
        dueDays
        billingEmail
        billingEmailCC
        billingEmailBCC
      }
      upcomingInvoices {
        metadata {
          id
        }
        invoicePeriodEnd
        invoicePeriodStart
        status
        issued
        amountDue
        due
        currency
        invoiceLineItems {
          metadata {
            id
            created
          }

          quantity
          subtotal
          taxDue
          total
          price
          description
        }
        contract {
          billingDetails {
            canPayWithBankTransfer
          }
        }
        status
        invoiceNumber
        invoicePeriodStart
        invoicePeriodEnd
        invoiceUrl
        due
        issued
        subtotal
        taxDue
        currency
        note
        customer {
          name
          email
          addressLine1
          addressLine2
          addressZip
          addressLocality
          addressCountry
          addressRegion
        }
        provider {
          name
          addressLine1
          addressLine2
          addressZip
          addressLocality
          addressCountry
        }
      }
      opportunities {
        id
        comments
        internalStage
        internalType
        amount
        maxAmount
        name
        renewalLikelihood
        renewalAdjustedRate
        renewalUpdatedByUserId
        renewedAt
        updatedAt

        owner {
          id
          firstName
          lastName
          name
        }
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
