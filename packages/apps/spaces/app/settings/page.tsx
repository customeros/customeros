'use client';
import { useSearchParams } from 'next/navigation';

import { Panels } from './src/components/Tabs/Panels';
import { TabsContainer } from './src/components/Tabs/TabsContainer';

export default function SettingsPage() {
  const searchParams = useSearchParams();

  const tab = searchParams?.get('tab') ?? 'oauth';

  return (
    <TabsContainer>
      <Panels tab={tab} />
    </TabsContainer>
  );
}
