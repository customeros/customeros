import React, { useContext, PropsWithChildren } from 'react';

import { reaction } from 'mobx';
import { observer } from 'mobx-react-lite';

import { ServiceFormStore } from './Services.store';
import InvoicePreviewListStore from './InvoicePreviewList.store';

interface EditContractModalStoreContextValue {
  serviceFormStore: ServiceFormStore;
  invoicePreviewList: InvoicePreviewListStore;
}

class Store {
  serviceFormStore: ServiceFormStore;
  invoicePreviewList: InvoicePreviewListStore;
  constructor() {
    this.serviceFormStore = new ServiceFormStore();
    this.invoicePreviewList = new InvoicePreviewListStore();
    this.setupReactions();
  }

  private setupReactions() {
    reaction(
      () => this.serviceFormStore.shouldReact(),
      () => {
        this.serviceFormStore.runSimulation(this.invoicePreviewList);
      },
      {
        delay: 500,
      },
    );
  }
}

const store = new Store();

const EditContractModalStoreContext =
  React.createContext<EditContractModalStoreContextValue | null>(store);

export const useEditContractModalStores = () => {
  const context = useContext(EditContractModalStoreContext);
  if (context === null)
    throw new Error(
      'useEditContractModalStores hook must be used within a EditContractModalStoreContextProvider',
    );

  return context;
};

export const EditContractModalStoreContextProvider = observer(
  ({ children }: PropsWithChildren) => {
    return (
      <EditContractModalStoreContext.Provider value={store}>
        {children}
      </EditContractModalStoreContext.Provider>
    );
  },
);
