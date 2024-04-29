import { useContext } from 'react';

import { StoreContext } from '@shared/components/Providers/StoreProvider';

export const useStore = () => {
  return useContext(StoreContext);
};
