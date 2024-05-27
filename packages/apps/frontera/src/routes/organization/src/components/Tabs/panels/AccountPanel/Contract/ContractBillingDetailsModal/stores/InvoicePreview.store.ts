import { makeAutoObservable } from 'mobx';

import { InvoiceSimulateServiceLineInput } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/InvoicePreviewList.store';

class InvoiceStore {
  invoice: InvoiceSimulateServiceLineInput | null = null;

  constructor(invoice: InvoiceSimulateServiceLineInput) {
    makeAutoObservable(this);
    this.invoice = invoice;
  }
}

export default InvoiceStore;
