import { useMemo, createContext } from 'react';

import { RootStore } from '@store/root';
import { Transport } from '@store/transport';

export const StoreContext = createContext<RootStore>({} as RootStore);

export const StoreProvider = ({ children }: { children: React.ReactNode }) => {
  const rootStore = useMemo(() => {
    return new RootStore(new Transport());
  }, []);

  return (
    <StoreContext.Provider value={rootStore}>{children}</StoreContext.Provider>
  );
};
