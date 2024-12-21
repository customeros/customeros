import { useMemo, createContext } from 'react';

import { RootStore } from '@store/root';

export const StoreContext = createContext<RootStore>({} as RootStore);

export const StoreProvider = ({ children }: { children: React.ReactNode }) => {
  const rootStore = useMemo(() => {
    return RootStore.getInstance();
  }, []);

  return (
    <StoreContext.Provider value={rootStore}>{children}</StoreContext.Provider>
  );
};
