'use client';

import React, { useState, useEffect } from 'react';

import { RecoilRoot } from 'recoil';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
// import { compress, decompress } from 'lz-string';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
// import { PersistQueryClientProvider } from '@tanstack/react-query-persist-client';
// import { createSyncStoragePersister } from '@tanstack/query-sync-storage-persister';

import { StoreProvider } from './StoreProvider';
import { Env, EnvProvider } from './EnvProvider';
import { NextAuthProvider } from './SessionProvider';
import { AnalyticsProvider } from './AnalyticsProvider';
import { PhoenixSocketProvider } from './SocketProvider';
import { GrowthbookProvider } from './GrowthbookProvider';
import { IntegrationsProvider } from './IntegrationsProvider';
import { NotificationsProvider } from './NotificationsProvider';
interface ProvidersProps {
  env: Env;
  isProduction?: boolean;
  children: React.ReactNode;
}

export const Providers = ({ env, children, isProduction }: ProvidersProps) => {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            gcTime: 1000 * 60 * 60 * 24, // 24 hours
          },
        },
      }),
  );

  useEffect(() => {
    // Temporary: should be removed after a few weeks
    if (typeof window === 'undefined') return;
    localStorage.removeItem('REACT_QUERY_OFFLINE_CACHE');
    indexedDB.deleteDatabase('keyval-store');
  }, []);

  return (
    <EnvProvider env={env}>
      <QueryClientProvider client={queryClient}>
        <StoreProvider>
          <ReactQueryDevtools initialIsOpen={false} position='bottom' />
          <PhoenixSocketProvider>
            <RecoilRoot>
              <NextAuthProvider>
                <IntegrationsProvider>
                  <GrowthbookProvider>
                    <NotificationsProvider isProduction={isProduction}>
                      <AnalyticsProvider isProduction={isProduction}>
                        {children}
                      </AnalyticsProvider>
                    </NotificationsProvider>
                  </GrowthbookProvider>
                </IntegrationsProvider>
              </NextAuthProvider>
            </RecoilRoot>
          </PhoenixSocketProvider>
        </StoreProvider>
      </QueryClientProvider>
    </EnvProvider>
  );
};
