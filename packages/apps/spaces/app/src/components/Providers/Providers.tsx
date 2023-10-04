'use client';

import { useState } from 'react';
import { RecoilRoot } from 'recoil';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { NextAuthProvider } from './SessionProvider';
import { JuneProvider } from '@shared/components/Providers/JuneProvider';

export const Providers = ({ children }: { children: React.ReactNode }) => {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <ReactQueryDevtools initialIsOpen={false} position='bottom-right' />
      <RecoilRoot>
        <NextAuthProvider>
          <JuneProvider>{children}</JuneProvider>
        </NextAuthProvider>
      </RecoilRoot>
    </QueryClientProvider>
  );
};
