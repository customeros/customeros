import { gql } from 'graphql-request';
import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { SyncableGroup } from '@store/syncable-group.ts';
import { when, observable, runInAction, makeObservable } from 'mobx';

import { Filter, SortBy, Invoice, Pagination } from '@graphql/types';

import mock from './mock.json';
import { InvoiceStore } from './Invoice.store.ts';

export class InvoicesStore extends SyncableGroup<Invoice, InvoiceStore> {
  totalElements = 0;

  get channelName() {
    return 'Invoices';
  }

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, InvoiceStore);

    makeObservable(this, {
      totalElements: observable,
    });

    when(
      () =>
        this.isBootstrapped && this.totalElements > 0 && !this.root.demoMode,
      async () => {
        await this.bootstrapRest();
      },
    );
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray<T extends InvoiceStore>(
    compute: (arr: InvoiceStore[]) => T[],
  ) {
    const arr = this.toArray();

    return compute(arr);
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.load(mock.data.invoices.content as Invoice[], {
        getId: (data) => data.invoiceNumber,
      });
      this.isBootstrapped = true;
      this.totalElements = mock.data.invoices.totalElements;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;

      const { invoices } = await this.transport.graphql.request<
        INVOICES_QUERY_RESPONSE,
        INVOICES_QUERY_PAYLOAD
      >(INVOICES_QUERY, {
        pagination: { limit: 1000, page: 0 },
        sort: [],
      });

      this.load(invoices.content, {
        getId: (data) => data.invoiceNumber, // this change is intentional, preview id changes between updates and number stays stable
      });
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = invoices.totalElements;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async bootstrapRest() {
    let page = 1;

    while (this.totalElements > this.value.size) {
      try {
        const { invoices } = await this.transport.graphql.request<
          INVOICES_QUERY_RESPONSE,
          INVOICES_QUERY_PAYLOAD
        >(INVOICES_QUERY, {
          pagination: { limit: 1000, page },
          sort: [],
        });

        runInAction(() => {
          page++;
          this.load(invoices.content, {
            getId: (data) => data.invoiceNumber,
          });
        });
      } catch (e) {
        runInAction(() => {
          this.error = (e as Error)?.message;
        });
        break;
      }
    }
  }
}

type INVOICES_QUERY_PAYLOAD = {
  where?: Filter;
  sort?: SortBy[];
  pagination: Pagination;
};
type INVOICES_QUERY_RESPONSE = {
  invoices: {
    content: Invoice[];
    totalElements: number;
    totalAvailable: number;
  };
};
const INVOICES_QUERY = gql`
  query getInvoices(
    $pagination: Pagination!
    $where: Filter
    $sort: [SortBy!]
  ) {
    invoices(pagination: $pagination, where: $where, sort: $sort) {
      content {
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
        }
        contract {
          metadata {
            id
          }
          billingDetails {
            billingCycleInMonths
            canPayWithBankTransfer
          }
          contractName
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
        preview
        subtotal
        taxDue
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
      totalElements
      totalAvailable
    }
  }
`;
