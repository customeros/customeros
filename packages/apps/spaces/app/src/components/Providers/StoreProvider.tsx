'use client';

import { useEffect, useContext, createContext } from 'react';

import { RootStore } from '@store/root';
import { TransportLayer } from '@store/transport';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { EnvContext } from './EnvProvider';

export const StoreContext = createContext<RootStore>({} as RootStore);

export const StoreProvider = ({ children }: { children: React.ReactNode }) => {
  const client = getGraphQLClient();
  const env = useContext(EnvContext);

  const { data } = useGlobalCacheQuery(client);

  const user = data?.global_Cache?.user;
  const user_id = user?.id ?? '';
  const username = (() => {
    if (!user) return '';
    const fullName = [user?.firstName, user?.lastName].join(' ').trim();
    const email = user?.emails?.[0]?.email;

    return fullName || email || '';
  })();

  const transportLayer = new TransportLayer({
    token: env?.REALTIME_WS_API_KEY,
    socketPath: `${env?.REALTIME_WS_PATH}/socket`,
  });

  useEffect(() => {
    user && transportLayer.setMetadata({ user_id, username });
  }, [user]);

  return (
    <StoreContext.Provider value={new RootStore(transportLayer)}>
      {children}
    </StoreContext.Provider>
  );
};
