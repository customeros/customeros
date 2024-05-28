import React, {
  useState,
  useEffect,
  useContext,
  PropsWithChildren,
} from 'react';

import { reaction } from 'mobx';
import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

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

  runSimulationReaction() {
    return reaction(
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
    const isSimulationEnabled = useFeatureIsOn('invoice-simulation');
    const [store] = useState(() => new Store());

    useEffect(() => {
      if (isSimulationEnabled) {
        const disposer = store.runSimulationReaction();

        return () => {
          disposer();
        };
      }
    }, [store, isSimulationEnabled]);

    return (
      <EditContractModalStoreContext.Provider value={store}>
        {children}
      </EditContractModalStoreContext.Provider>
    );
  },
);
