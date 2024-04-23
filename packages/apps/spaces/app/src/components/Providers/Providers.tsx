'use client';

import React, { useState } from 'react';

import { RecoilRoot } from 'recoil';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

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
