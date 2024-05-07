import { useMemo, useState, useEffect, createContext } from 'react';

import { autorun } from 'mobx';
import { RootStore } from '@store/root';
import { observer } from 'mobx-react-lite';
import { TransportLayer, TransportLayerOptions } from '@store/transport';

export const StoreContext = createContext<RootStore>({} as RootStore);

export const StoreProvider = observer(
  ({ children }: { children: React.ReactNode }) => {
    const [transportOptions, setTransportOptions] =
      useState<TransportLayerOptions>();

    const transportLayer = useMemo(() => {
      return new TransportLayer(transportOptions);
    }, [transportOptions]);

    const rootStore = useMemo(() => {
      return new RootStore(transportLayer);
    }, [transportLayer]);

    useEffect(() => {
      autorun(() => {
        // temporary - will be removed once we drop react-query and getGraphQLClient
        const persisted = window.localStorage.getItem('__COS_SESSION__');
        window.__COS_SESSION__ = JSON.parse(persisted ?? '{}');

        if (rootStore.sessionStore.sessionToken) {
          setTransportOptions({
            email: rootStore.sessionStore.value.profile.email,
            userId: rootStore.sessionStore.value.profile.id,
            sessionToken: rootStore.sessionStore.sessionToken,
          });
        }
      });
    }, []);

    return (
      <StoreContext.Provider value={rootStore}>
        {children}
      </StoreContext.Provider>
    );
  },
);
