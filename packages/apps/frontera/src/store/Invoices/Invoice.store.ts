import type { RootStore } from '@store/root';

import merge from 'lodash/merge';
import { match } from 'ts-pattern';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';
import { Syncable } from '@store/syncable.ts';
import { action, override, computed, runInAction, makeObservable } from 'mobx';

import {
  Invoice,
  Contract,
  Currency,
  Metadata,
  DataSource,
  Organization,
  InvoiceStatus,
  InvoiceUpdateInput,
} from '@graphql/types';

export class InvoiceStore extends Syncable<Invoice> {
  constructor(
    public root: RootStore,
    public transport: Transport,
    data: Invoice,
  ) {
    super(root, transport, data ?? getDefaultValue());
    makeObservable<InvoiceStore, 'updateInvoiceStatus' | 'contract'>(this, {
      id: override,
      getId: override,
      setId: override,
      invalidate: action,
      updateInvoiceStatus: action,
      contract: computed,
      provider: computed,
      bankAccounts: computed,
      number: computed,
    });
  }

  get id() {
    return this.value.metadata?.id;
  }

  get number() {
    return this.value.invoiceNumber;
  }

  set id(id: string) {
    this.value.metadata.id = id;
  }

  get contract() {
    return this.root.contracts.value.get(this.value.contract.metadata.id)
      ?.value;
  }

  get provider() {
    return this.root.settings.tenantBillingProfiles.toArray()?.[0]?.value;
  }

  get bankAccounts() {
    return this.root.settings.bankAccounts?.toArray();
  }

  getId() {
    return this.value.metadata.id;
  }

  getChannelName() {
    return 'Invoices';
  }

  async invalidate() {
    try {
      this.isLoading = true;

      const { invoice } = await this.transport.graphql.request<
        INVOICE_QUERY_RESULT,
        { number: string }
      >(INVOICES_QUERY, { number: this.number });

      await this.load(invoice);
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

  private async updateInvoiceStatus(payload: InvoiceUpdateInput) {
    try {
      this.isLoading = true;
      await this.transport.graphql.request<UPDATE_INVOICE_STATUS_MUTATION_PAYLOAD>(
        UPDATE_INVOICE_STATUS_MUTATION,
        {
          input: {
            ...payload,
            id: this.id,
            patch: true,
          },
        },
      );

      runInAction(() => {
        this.invalidate();
      });
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

  async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const path = diff?.path;

    match(path)
      .with(['status'], () => {
        const payload = makePayload<InvoiceUpdateInput>(operation);

        this.updateInvoiceStatus(payload);
      })

      .otherwise(() => {});
  }

  init(data: Invoice) {
    return merge(this.value, data);
  }
}

type INVOICE_QUERY_RESULT = {
  invoice: Invoice;
};
const INVOICES_QUERY = gql`
  query Invoice($number: String!) {
    invoice_ByNumber(number: $number) {
      issued
      metadata {
        id
        created
      }
      organization {
        metadata {
          id
        }
      }
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
        logoUrl
        logoRepositoryFileId
        name
        addressLine1
        addressLine2
        addressZip
        addressLocality
        addressCountry
        addressRegion
      }
      contract {
        metadata {
          id
        }
        contractName
        billingDetails {
          canPayWithBankTransfer
          billingCycleInMonths
        }
      }
      invoiceUrl
      invoiceNumber
      invoicePeriodStart
      invoicePeriodEnd
      due
      issued
      amountDue
      currency
      dryRun
      status
      subtotal
      invoiceLineItems {
        metadata {
          id
          created
          lastUpdated
          source
          sourceOfTruth
          appSource
        }
        contractLineItem {
          serviceStarted
          price
          billingCycle
        }
        quantity
        subtotal
        taxDue
        total
        price
        description
      }
    }
  }
`;

type UPDATE_INVOICE_STATUS_MUTATION_PAYLOAD = {
  input: InvoiceUpdateInput;
};

const UPDATE_INVOICE_STATUS_MUTATION = gql`
  mutation UpdateInvoiceStatus($input: InvoiceUpdateInput!) {
    invoice_Update(input: $input) {
      metadata {
        id
      }
    }
  }
`;

const getDefaultValue = (): Invoice => ({
  metadata: {
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    appSource: DataSource.Openline,
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
  },
  organization: {
    metadata: {
      id: crypto.randomUUID(),
    } as Metadata,
  } as Organization,
  contract: {
    metadata: {
      id: crypto.randomUUID(),
    } as Metadata,
  } as Contract,
  issued: new Date().toISOString(),
  invoiceNumber: '',
  billingCycleInMonths: 1,
  invoicePeriodStart: new Date().toISOString(),
  invoicePeriodEnd: new Date().toISOString(),
  due: new Date().toISOString(),
  amountDue: 0,
  currency: Currency.Usd,
  dryRun: false,
  status: InvoiceStatus.Due,
  invoiceLineItems: [],
  paid: false,
  subtotal: 0,
  taxDue: 0,
  paymentLink: '',
  repositoryFileId: '',
  note: '',
  amountPaid: 0,
  amountRemaining: 0,
  preview: false,
  offCycle: false,
  postpaid: false,
  invoiceUrl: '',
  customer: {
    name: '',
    email: '',
    addressLine1: '',
    addressLine2: '',
    addressZip: '',
    addressLocality: '',
    addressCountry: '',
    addressRegion: '',
  },
  provider: {
    logoUrl: '',
    logoRepositoryFileId: '',
    name: '',
    addressLine1: '',
    addressLine2: '',
    addressZip: '',
    addressLocality: '',
    addressCountry: '',
    addressRegion: '',
  },
});
