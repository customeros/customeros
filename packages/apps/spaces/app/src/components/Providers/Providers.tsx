'use client';

import { useState } from 'react';

import { RecoilRoot } from 'recoil';
import { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { PersistQueryClientProvider } from '@tanstack/react-query-persist-client';

import { createIDBPersister } from '@shared/util/indexedDBPersister';
import { AnalyticsProvider } from '@shared/components/Providers/AnalyticsProvider';

import { NextAuthProvider } from './SessionProvider';
interface ProvidersProps {
  isProduction?: boolean;
  children: React.ReactNode;
  sessionEmail?: string | null;
}

const hostname =
  typeof window !== 'undefined' ? window?.location?.hostname : 'platform';

export const Providers = ({
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
            cacheTime: 1000 * 60 * 60 * 24, // 24 hours
          },
        },
      }),
  );

  return (
    <PersistQueryClientProvider
      client={queryClient}
      persistOptions={{ persister }}
    >
      <ReactQueryDevtools initialIsOpen={false} position='bottom-right' />
      <RecoilRoot>
        <NextAuthProvider>
          <AnalyticsProvider isProduction={isProduction}>
            {children}
          </AnalyticsProvider>
        </NextAuthProvider>
      </RecoilRoot>
    </PersistQueryClientProvider>
  );
};
