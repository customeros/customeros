import { makeAutoObservable } from 'mobx';

import { InvoiceSimulate, InvoiceLineSimulate } from '@graphql/types';
import ServiceLineItemStore from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/Service.store';

import InvoicePreviewStore from './InvoicePreview.store';

export interface ISimulatedInvoiceLineItems extends InvoiceLineSimulate {
  serviceLineItemStore: ServiceLineItemStore | null;
}
export interface InvoiceSimulateServiceLineInput extends InvoiceSimulate {
  invoiceLineItems: ISimulatedInvoiceLineItems[];
}
class InvoiceListStore {
  simulatedInvoices: Array<InvoicePreviewStore> = [];
  isPending = false;
  previewedInvoiceIndex: number = 0;

  constructor() {
    makeAutoObservable(this);
  }

  setPreviewedInvoice(invoiceIndex: number) {
    this.previewedInvoiceIndex = invoiceIndex;
  }
  initializeSimulatedInvoices(
    invoices: Array<InvoiceSimulateServiceLineInput>,
  ) {
    this.simulatedInvoices = invoices.map(
      (invoice) => new InvoicePreviewStore(invoice),
    );

    this.setPreviewedInvoice(0);
  }
  resetSimulatedInvoices() {
    this.simulatedInvoices = [];
  }
}

export default InvoiceListStore;
