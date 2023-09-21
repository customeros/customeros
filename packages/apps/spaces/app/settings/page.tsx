'use client';
import { useState, useEffect } from 'react';
import { useSession } from 'next-auth/react';
import { IntegrationAppProvider } from '@integration-app/react';

import { TabsContainer } from './Tabs/TabsContainer';
import { Panels } from './Tabs/Panels';
import { SettingsMainSection } from './SettingsMainSection';

interface SettingsPageProps {
  searchParams: { tab?: string };
}

export default function SettingsPage({ searchParams }: SettingsPageProps) {
  const [integrationToken, setIntegrationToken] = useState<
    string | undefined
  >();
  const { data: session } = useSession();

  useEffect(() => {
    if (session?.user) {
      (async () => {
        try {
          const response = await fetch('/api/integration-token');
          const data = await response?.json();
          setIntegrationToken(data);
        } catch (e) {
          // handle error
        }
      })();
    }
  }, [session]);

  return (
    <IntegrationAppProvider token={integrationToken}>
      <SettingsMainSection>
        <TabsContainer>
          <Panels tab={searchParams.tab ?? 'oauth'} />
        </TabsContainer>
      </SettingsMainSection>
    </IntegrationAppProvider>
  );
}
