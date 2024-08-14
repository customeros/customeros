import { useSearchParams } from 'react-router-dom';

import { useStore } from '@shared/hooks/useStore';

import { Panels } from './src/components/Tabs/Panels';
import { TabsContainer } from './src/components/Tabs/TabsContainer';

export const SettingsPage = () => {
  const [searchParams] = useSearchParams();
  const store = useStore();
  const tab = searchParams?.get('tab') ?? 'workspace';

  store.ui.commandMenu.setType('GlobalHub');

  return (
    <TabsContainer>
      <Panels tab={tab} />
    </TabsContainer>
  );
};
