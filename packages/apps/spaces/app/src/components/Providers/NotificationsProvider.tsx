'use client';

import React from 'react';

import { NovuProvider } from '@novu/notification-center';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

interface ProvidersProps {
  isProduction?: boolean;
  children: React.ReactNode;
}

export const NotificationsProvider = ({ children }: ProvidersProps) => {
  const client = getGraphQLClient();

  const { data: globalCacheQuery } = useGlobalCacheQuery(client);

  const id = globalCacheQuery?.global_Cache?.user.id;

  return (
    <NovuProvider
      subscriberId={id}
      applicationIdentifier={
        process.env.NEXT_PUBLIC_NOTIFICATION_TEST_APP_IDENTIFIER as string
      }
    >
      {children}
    </NovuProvider>
  );
};
