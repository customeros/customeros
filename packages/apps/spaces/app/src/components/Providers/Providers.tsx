'use client';

import React, { useState } from 'react';

import { RecoilRoot } from 'recoil';
import { compress, decompress } from 'lz-string';
import { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { PersistQueryClientProvider } from '@tanstack/react-query-persist-client';
import { createSyncStoragePersister } from '@tanstack/query-sync-storage-persister';

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
  const [persister] = useState(() =>
    createSyncStoragePersister({
      storage: typeof window !== 'undefined' ? window?.localStorage : null,
      serialize: (data) => compress(JSON.stringify(data)),
      deserialize: (data) => JSON.parse(decompress(data)),
    }),
  );

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

  return (
    <EnvProvider env={env}>
      <PersistQueryClientProvider
        client={queryClient}
        persistOptions={{ persister, maxAge: 1000 * 60 * 60 * 12 }}
      >
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
      </PersistQueryClientProvider>
    </EnvProvider>
  );
};
