'use client';

import { useMemo, useEffect, useContext, createContext } from 'react';

import { RootStore } from '@store/root';
import { TransportLayer } from '@store/transport';
import { enableStaticRendering } from 'mobx-react-lite';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { EnvContext } from './EnvProvider';

enableStaticRendering(typeof window === 'undefined');

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

  const transportLayer = useMemo(
    () =>
      new TransportLayer({
        token: env?.REALTIME_WS_API_KEY,
        socketPath: `${env?.REALTIME_WS_PATH}/socket`,
      }),
    [env?.REALTIME_WS_API_KEY, env?.REALTIME_WS_PATH],
  );

  useEffect(() => {
    if (user_id && username) {
      transportLayer.setMetadata({ user_id, username });
    }

    () => {
      transportLayer.disconnect();
    };
  }, [user_id, username, transportLayer]);

  const rootStore = useMemo(
    () => new RootStore(transportLayer),
    [transportLayer],
  );

  return (
    <StoreContext.Provider value={rootStore}>{children}</StoreContext.Provider>
  );
};
