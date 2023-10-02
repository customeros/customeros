import React from 'react';

import { TabsContainer } from './Tabs/TabsContainer';
import { Panels } from './Tabs/Panels';

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
