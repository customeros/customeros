'use client';

import { SessionProvider } from 'next-auth/react';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

type Props = {
  children?: React.ReactNode;
};

export const NextAuthProvider = ({ children }: Props) => {
  const client = getGraphQLClient();
  useGlobalCacheQuery(client);

  return <SessionProvider>{children}</SessionProvider>;
};
