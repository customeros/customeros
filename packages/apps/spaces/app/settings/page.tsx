import React from 'react';

import { TabsContainer } from './src/components/Tabs/TabsContainer';
import { Panels } from './src/components/Tabs/Panels';

interface SettingsPageProps {
  searchParams: { tab?: string };
}

export default async function SettingsPage({
  searchParams,
}: SettingsPageProps) {
  return (
    <TabsContainer>
      <Panels tab={searchParams.tab ?? 'oauth'} />
    </TabsContainer>
  );
}
