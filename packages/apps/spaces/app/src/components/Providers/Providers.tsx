'use client';

import { useState } from 'react';
import { RecoilRoot } from 'recoil';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { AnalyticsProvider } from '@shared/components/Providers/AnalyticsProvider';
import { NextAuthProvider } from './SessionProvider';

interface ProvidersProps {
  sessionEmail?: string | null;
  children: React.ReactNode;
}

export const Providers = ({ children, sessionEmail }: ProvidersProps) => {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <ReactQueryDevtools initialIsOpen={false} position='bottom-right' />
      <RecoilRoot>
        <NextAuthProvider>
          <AnalyticsProvider>{children}</AnalyticsProvider>
        </NextAuthProvider>
      </RecoilRoot>
    </QueryClientProvider>
  );
};
