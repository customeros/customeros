import type { RootStore } from '@store/root';

import merge from 'lodash/merge';
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
} from '@graphql/types';

import { InvoicesService } from './__service__/Invoices.service';

export class InvoiceStore extends Syncable<Invoice> {
  private service: InvoicesService;

  constructor(
    public root: RootStore,
    public transport: Transport,
    data: Invoice,
  ) {
    super(root, transport, data ?? getDefaultValue());
    this.service = InvoicesService.getInstance(transport);
    makeObservable<InvoiceStore>(this, {
      id: override,
      save: override,
      getId: override,
      setId: override,
      invalidate: action,
      contract: computed,
      getChannelName: override,
      provider: computed,
      bankAccounts: computed,
      number: computed,
      updateInvoiceStatus: action,
      getInvoiceStatus: action,
    });
  }

  get id() {
    return this.value.invoiceNumber;
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

      const { invoice } = await this.service.getInvoice(this.number);

      await this.load(invoice as Invoice);
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

  getInvoiceStatus() {
    return this.value.status;
  }

  init(data: Invoice) {
    return merge(this.value, data);
  }

  updateInvoiceStatus(status: InvoiceStatus) {
    this.value.status = status;
  }

  static getDefaultValue(): Invoice {
    return {
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
    };
  }
}

export const getDefaultValue = (): Invoice => ({
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
