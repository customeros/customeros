import { useSearchParams } from 'react-router-dom';

import { Panels } from './src/components/Tabs/Panels';
import { TabsContainer } from './src/components/Tabs/TabsContainer';

export const SettingsPage = () => {
  const [searchParams] = useSearchParams();

  const tab = searchParams?.get('tab') ?? 'oauth';

  return (
    <TabsContainer>
      <Panels tab={tab} />
    </TabsContainer>
  );
};
