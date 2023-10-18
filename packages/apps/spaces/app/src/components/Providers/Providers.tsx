'use client';

import { useState } from 'react';
import { RecoilRoot } from 'recoil';
import { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import {
  PersistQueryClientProvider,
  Persister,
} from '@tanstack/react-query-persist-client';

import { createIDBPersister } from '@shared/util/createIDBPersister';
import { AnalyticsProvider } from '@shared/components/Providers/AnalyticsProvider';

import { NextAuthProvider } from './SessionProvider';

let persister: Persister;
if (typeof window !== 'undefined') {
  persister = createIDBPersister(`cos-${window?.location?.hostname}`);
}

export const Providers = ({ children }: { children: React.ReactNode }) => {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 1000 * 60 * 1, // 1 minutes
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
