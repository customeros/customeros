import { useState, useEffect } from 'react';

// import { useSession } from 'next-auth/react';
import { IntegrationAppProvider } from '@integration-app/react';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';

export const IntegrationsProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [integrationToken, setIntegrationToken] = useState<
    string | undefined
  >();

  const client = getGraphQLClient();
  const { data: tenant } = useTenantNameQuery(client);
  // const { data: session } = useSession();

  // useEffect(() => {
  //   if (session?.user && tenant?.tenant) {
  //     (async () => {
  //       try {
  //         const response = await fetch(
  //           `/api/integration/token?tenant=${tenant.tenant}`,
  //         );
  //         const data = await response?.json();
  //         setIntegrationToken(data.token);
  //       } catch (e) {
  //         toastError('Failed to fetch integration token', 'integration-token');
  //       }
  //     })();
  //   }
  // }, [session, tenant]);

  return (
    <IntegrationAppProvider token={integrationToken}>
      {children}
    </IntegrationAppProvider>
  );
};
