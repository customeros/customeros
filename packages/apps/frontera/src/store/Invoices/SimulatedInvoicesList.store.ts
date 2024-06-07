// import { Store } from '@store/store.ts';
// import { GroupStore } from '@store/group-store.ts';
// import { when, runInAction, makeAutoObservable } from 'mobx';
//
// import { Invoice, InvoiceSimulate, InvoiceLineSimulate } from '@graphql/types';
// import ServiceLineItemStore from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/Service.store';
//
// import SimulatedInvoiceStore from './SimulatedInvoice.store';
//
// export interface ISimulatedInvoiceLineItems extends InvoiceLineSimulate {
//   serviceLineItemStore: ServiceLineItemStore | null;
// }
// export interface InvoiceSimulateServiceLineInput extends InvoiceSimulate {
//   invoiceLineItems: ISimulatedInvoiceLineItems[];
// }
// class InvoiceListStore {
//   simulatedInvoices: Array<InvoicePreviewStore> = [];
//   isLoading = false;
//   previewedInvoiceIndex: number = 0;
//   value: Map<string, Store<Invoice>> = new Map();
//
//   constructor() {
//     makeAutoObservable(this);
//   }
//
//   load(this: GroupStore<SimulatedInvoiceStore>, data: T[]) {
//     data.forEach((item) => {
//       const id = item.metadata?.id;
//       if (this.value.has(id)) {
//         this.value.get(id)?.load(item);
//
//         return;
//       }
//
//       const itemStore = new SimulatedInvoiceStore(this.root);
//       itemStore.load(item);
//       this.value.set(id, itemStore);
//     });
//
//     this.isBootstrapped = true;
//   }
//
//   async bootstrap() {
//     if (this.isBootstrapped || this.isLoading) return;
//
//     try {
//       this.isLoading = true;
//       const { invoices } = await this.transport.graphql.request<
//         INVOICES_QUERY_RESPONSE,
//         INVOICES_QUERY_PAYLOAD
//       >(INVOICES_QUERY, {
//         pagination: { limit: 1000, page: 0 },
//         sort: [],
//       });
//
//       this.load(invoices.content);
//       runInAction(() => {
//         this.isBootstrapped = true;
//         this.totalElements = invoices.totalElements;
//       });
//     } catch (e) {
//       runInAction(() => {
//         this.error = (e as Error)?.message;
//       });
//     } finally {
//       runInAction(() => {
//         this.isLoading = false;
//       });
//     }
//   }
//   resetSimulatedInvoices() {
//     this.simulatedInvoices = [];
//   }
// }
//
// export default InvoiceListStore;
