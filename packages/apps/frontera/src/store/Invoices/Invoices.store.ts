import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Store } from '@store/store.ts';
import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { GroupOperation } from '@store/types.ts';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store.ts';

import {
  Filter,
  SortBy,
  Invoice,
  Pagination,
  Organization,
  SortingDirection,
} from '@graphql/types';

import { InvoiceStore } from './Invoice.store.ts';

export class InvoicesStore implements GroupStore<Invoice> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<Invoice>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<Invoice>();
  totalElements = 0;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Invoices',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: InvoiceStore,
    });
  }
  toArray() {
    return Array.from(this.value.values());
  }
  toComputedArray<T extends Store<Invoice>>(
    compute: (arr: Store<Invoice>[]) => T[],
  ) {
    const arr = this.toArray();

    return compute(arr);
  }

  async bootstrap() {
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

      this.load(invoices.content);
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
          this.load(invoices.content);
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
