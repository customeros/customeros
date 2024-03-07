'use client';

import React from 'react';

import { NovuProvider } from '@novu/notification-center';

import { useEnv } from '@shared/hooks/useEnv';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

interface ProvidersProps {
  isProduction?: boolean;
  children: React.ReactNode;
}

export const NotificationsProvider = ({
  children,
  isProduction,
}: ProvidersProps) => {
  const env = useEnv();
  const client = getGraphQLClient();

  const { data: globalCacheQuery } = useGlobalCacheQuery(client);

  const id = globalCacheQuery?.global_Cache?.user.id ?? 'temp-id';
  const applicationIdentifier = isProduction
    ? env.NOTIFICATION_PROD_APP_IDENTIFIER
    : env.NOTIFICATION_TEST_APP_IDENTIFIER;

  return (
    <NovuProvider
      subscriberId={id}
      applicationIdentifier={applicationIdentifier}
    >
      {children}
    </NovuProvider>
  );
};
