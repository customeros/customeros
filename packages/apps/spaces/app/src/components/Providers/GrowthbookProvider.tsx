import { useEffect } from 'react';

import { GrowthBook } from '@growthbook/growthbook-react';
import { GrowthBookProvider } from '@growthbook/growthbook-react';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';

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
  const { data } = useTenantNameQuery(client);

  useEffect(() => {
    growthbook.loadFeatures();
  }, []);

  useEffect(() => {
    growthbook.setAttributes({
      tenant: data?.tenant,
    });
  }, [data?.tenant]);

  return (
    <GrowthBookProvider growthbook={growthbook}>{children}</GrowthBookProvider>
  );
};
