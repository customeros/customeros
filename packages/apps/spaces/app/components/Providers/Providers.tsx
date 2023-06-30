'use client';

import { useState } from 'react';
import { RecoilRoot } from 'recoil';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';

import { SessionProvider } from './SessionProvider';

export const Providers = ({ children }: { children: React.ReactNode }) => {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <RecoilRoot>
        <SessionProvider>{children}</SessionProvider>
      </RecoilRoot>
    </QueryClientProvider>
  );
};
