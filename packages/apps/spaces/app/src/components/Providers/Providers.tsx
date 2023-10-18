'use client';

import { useState } from 'react';
import { RecoilRoot } from 'recoil';
import { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { PersistQueryClientProvider } from '@tanstack/react-query-persist-client';
import { createSyncStoragePersister } from '@tanstack/query-sync-storage-persister';

import { NextAuthProvider } from './SessionProvider';
import { AnalyticsProvider } from '@shared/components/Providers/AnalyticsProvider';

export const Providers = ({ children }: { children: React.ReactNode }) => {
  const [persister] = useState(() =>
    createSyncStoragePersister({
      storage: typeof window !== 'undefined' ? window?.localStorage : null,
    }),
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
          <AnalyticsProvider>{children}</AnalyticsProvider>
        </NextAuthProvider>
      </RecoilRoot>
    </PersistQueryClientProvider>
  );
};
