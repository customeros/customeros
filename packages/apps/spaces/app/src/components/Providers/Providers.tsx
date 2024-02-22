'use client';

import React, { useState } from 'react';

import { RecoilRoot } from 'recoil';
import { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { PersistQueryClientProvider } from '@tanstack/react-query-persist-client';

import { createIDBPersister } from '@shared/util/indexedDBPersister';
import { AnalyticsProvider } from '@shared/components/Providers/AnalyticsProvider';

import { Env, EnvProvider } from './EnvProvider';
import { NextAuthProvider } from './SessionProvider';
import { PhoenixSocketProvider } from './SocketProvider';
import { GrowthbookProvider } from './GrowthbookProvider';
import { NotificationsProvider } from './NotificationsProvider';
interface ProvidersProps {
  env: Env;
  isProduction?: boolean;
  children: React.ReactNode;
  sessionEmail?: string | null;
}

const hostname =
  typeof window !== 'undefined' ? window?.location?.hostname : 'platform';

export const Providers = ({
  env,
  children,
  sessionEmail,
  isProduction,
}: ProvidersProps) => {
  const [persister] = useState(() =>
    createIDBPersister(`${sessionEmail ?? 'cos'}-${hostname}`),
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
        persistOptions={{ persister }}
      >
        <ReactQueryDevtools initialIsOpen={false} position='bottom' />
        <PhoenixSocketProvider>
          <RecoilRoot>
            <NextAuthProvider>
              <GrowthbookProvider>
                <NotificationsProvider isProduction={isProduction}>
                  <AnalyticsProvider isProduction={isProduction}>
                    {children}
                  </AnalyticsProvider>
                </NotificationsProvider>
              </GrowthbookProvider>
            </NextAuthProvider>
          </RecoilRoot>
        </PhoenixSocketProvider>
      </PersistQueryClientProvider>
    </EnvProvider>
  );
};
