import type { RootStore } from '@store/root';

import set from 'lodash/set';
import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makePayload } from '@store/util.ts';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { runInAction, makeAutoObservable } from 'mobx';

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
    makeAutoSyncable(this, {
      channelName: 'Contract',
      mutator: this.save,
      getId: (d: Contract) => d?.metadata?.id,
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

  private transformHistoryToContractUpdateInput = (
    history: Operation[],
  ): ContractUpdateInput => {
    const contractUpdate: Partial<ContractUpdateInput> = {
      contractId: this.id,
    };

    history.forEach((change) => {
      change.diff.forEach((diffItem) => {
        const { path, val } = diffItem;
        const [fieldName, subField] = path;
        if (fieldName === 'contractLineItems') {
          return;
        }

        if (subField) {
          if (!contractUpdate[fieldName]) {
            contractUpdate[fieldName] = {};
          }
          (contractUpdate[fieldName] as Record<string, unknown>)[subField] =
            val;
        } else {
          (contractUpdate as Record<string, unknown>)[fieldName] = val;
        }
      });
    });

    return contractUpdate as ContractUpdateInput;
  };

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const path = diff?.path;

    if (!path) {
      const payload = this.transformHistoryToContractUpdateInput(this.history);
      await this.updateContract({
        patch: true,
        ...payload,
      });
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

  init(data: Contract) {
    if (data.contractName === 'Schimb numele acu') {
      console.log('init', data);
    }
    const output = merge(this.value, data);
    const contracts = data.contractLineItems?.map((item) => {
      this.root.contractLineItems.load([item]);

      return this.root.contractLineItems.value.get(item.metadata.id)?.value;
    });
    const opportunities = data.opportunities?.map((item) => {
      this.root.opportunities.load([item]);

      return this.root.opportunities.value.get(item.id)?.value;
    });
    const upcomingInvoices = data.upcomingInvoices?.map((item) => {
      const upcomingInvoice = this.root.invoices.value.get(
        item.metadata.id,
      )?.value;

      if (!upcomingInvoice) {
        this.root.invoices.load([item]);
      }

      return this.root.invoices.value.get(item.metadata.id)?.value;
    });

    set(output, 'contractLineItems', contracts);
    set(output, 'opportunities', opportunities);
    set(output, 'upcomingInvoices', upcomingInvoices);

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
