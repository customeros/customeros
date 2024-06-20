import { useMemo, createContext } from 'react';

import { RootStore } from '@store/root';
import { Transport } from '@store/transport';

export const StoreContext = createContext<RootStore>({} as RootStore);

export const StoreProvider = ({ children }: { children: React.ReactNode }) => {
  const demoMode = window.location.search.includes('demoMode');

  const rootStore = useMemo(() => {
    return new RootStore(new Transport(), demoMode);
  }, []);

  return (
    <StoreContext.Provider value={rootStore}>{children}</StoreContext.Provider>
  );
};
