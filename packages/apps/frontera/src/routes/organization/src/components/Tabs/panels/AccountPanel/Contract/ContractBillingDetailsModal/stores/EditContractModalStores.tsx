import React, { useState, useContext, PropsWithChildren } from 'react';

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
  }
}

const EditContractModalStoreContext =
  React.createContext<EditContractModalStoreContextValue | null>(null);

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
    const [store] = useState(() => new Store());

    return (
      <EditContractModalStoreContext.Provider value={store}>
        {children}
      </EditContractModalStoreContext.Provider>
    );
  },
);
