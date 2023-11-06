'use client';
import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';

import { useSession } from 'next-auth/react';
import { IntegrationAppProvider } from '@integration-app/react';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';

import { Panels } from './src/components/Tabs/Panels';
import { TabsContainer } from './src/components/Tabs/TabsContainer';

export default function SettingsPage() {
  const client = getGraphQLClient();
  const searchParams = useSearchParams();
  const { data: tenant } = useTenantNameQuery(client);
  const [integrationToken, setIntegrationToken] = useState<
    string | undefined
  >();
  const { data: session } = useSession();
  const tab = searchParams?.get('tab') ?? 'oauth';

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
        <Panels tab={tab} />
      </TabsContainer>
    </IntegrationAppProvider>
  );
}
