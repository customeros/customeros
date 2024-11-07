import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { SyncableGroup } from '@store/syncable-group.ts';
import { when, override, observable, runInAction, makeObservable } from 'mobx';

import { Invoice } from '@graphql/types';

import mock from './mock.json';
import { InvoiceStore } from './Invoice.store.ts';
import { InvoicesService } from './__service__/Invoices.service.ts';

export class InvoicesStore extends SyncableGroup<Invoice, InvoiceStore> {
  private service: InvoicesService;
  totalElements = 0;

  get channelName() {
    return 'Invoices';
  }

  get persisterKey() {
    return 'Invoices';
  }

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, InvoiceStore);
    this.service = InvoicesService.getInstance(transport);

    makeObservable<InvoicesStore>(this, {
      totalElements: observable,
      channelName: override,
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

      const { invoices } = await this.service.getInvoices({
        pagination: { limit: 1000, page: 1 },
      });

      this.load(invoices.content as Invoice[], {
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
        const { invoices } = await this.service.getInvoices({
          pagination: { limit: 1000, page },
        });

        runInAction(() => {
          page++;
          this.load(invoices.content as Invoice[], {
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
