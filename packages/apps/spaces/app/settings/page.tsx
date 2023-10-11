'use client';
import { useState, useEffect } from 'react';
import { useSession } from 'next-auth/react';
import { IntegrationAppProvider } from '@integration-app/react';

import { toastError } from '@ui/presentation/Toast';
import { Panels } from './src/components/Tabs/Panels';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { TabsContainer } from './src/components/Tabs/TabsContainer';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';

interface SettingsPageProps {
  searchParams: { tab?: string };
}

export default function SettingsPage({ searchParams }: SettingsPageProps) {
  const client = getGraphQLClient();
  const { data: tenant } = useTenantNameQuery(client);
  const [integrationToken, setIntegrationToken] = useState<
    string | undefined
  >();
  const { data: session } = useSession();

  useEffect(() => {
    if (session?.user && tenant?.tenant) {
      (async () => {
        try {
          const response = await fetch(
            `/api/integration/token?tenant=${tenant.tenant}`,
          );
          const data = await response?.json();
          setIntegrationToken(data.token);
        } catch (e) {
          toastError('Failed to fetch integration token', 'integration-token');
        }
      })();
    }
  }, [session, tenant]);

  return (
    <IntegrationAppProvider token={integrationToken}>
      <TabsContainer>
        <Panels tab={searchParams.tab ?? 'oauth'} />
      </TabsContainer>
    </IntegrationAppProvider>
  );
}
