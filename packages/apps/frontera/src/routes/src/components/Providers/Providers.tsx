import React, { useState } from 'react';
import { ToastContainer } from 'react-toastify';

import { RecoilRoot } from 'recoil';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { SplashScreen } from '@shared/components/SplashScreen/SplashScreen';

import { StoreProvider } from './StoreProvider';
import { Env, EnvProvider } from './EnvProvider';
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
          <SplashScreen>
            <ReactQueryDevtools initialIsOpen={false} position='bottom' />
            <PhoenixSocketProvider>
              <RecoilRoot>
                <IntegrationsProvider>
                  <GrowthbookProvider>
                    <NotificationsProvider isProduction={isProduction}>
                      <AnalyticsProvider isProduction={isProduction}>
                        {children}
                        <ToastContainer
                          position='bottom-right'
                          autoClose={8000}
                          limit={3}
                          closeOnClick={true}
                          hideProgressBar={true}
                          theme='colored'
                        />
                      </AnalyticsProvider>
                    </NotificationsProvider>
                  </GrowthbookProvider>
                </IntegrationsProvider>
              </RecoilRoot>
            </PhoenixSocketProvider>
          </SplashScreen>
        </StoreProvider>
      </QueryClientProvider>
    </EnvProvider>
  );
};
