import { useEffect } from 'react';

import { GrowthBook } from '@growthbook/growthbook-react';
import { GrowthBookProvider } from '@growthbook/growthbook-react';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

export const growthbook = new GrowthBook({
  apiHost: 'https://cdn.growthbook.io',
  clientKey: 'sdk-kU7RLceKTmkcTjxO',
  enableDevMode: true,
  subscribeToChanges: true,
  trackingCallback: (experiment, result) => {
    // TODO: Use real analytics tracking system
  },
});

export const GrowthbookProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const client = getGraphQLClient();
  const { data: tenantQuery } = useTenantNameQuery(client);
  const { data: globalCacheQuery } = useGlobalCacheQuery(client);

  const tenant = tenantQuery?.tenant;
  const id = globalCacheQuery?.global_Cache?.user.id;
  const email = globalCacheQuery?.global_Cache?.user?.emails?.[0]?.email;

  useEffect(() => {
    growthbook.loadFeatures();
  }, []);

  useEffect(() => {
    growthbook.setAttributes({
      id,
      email,
      tenant,
    });
  }, [tenant, id, email]);

  return (
    <GrowthBookProvider growthbook={growthbook}>{children}</GrowthBookProvider>
  );
};
