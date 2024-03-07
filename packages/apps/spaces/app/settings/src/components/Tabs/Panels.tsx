import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { AuthPanel } from './panels/AuthPanel';
import { BillingPanel } from './panels/BillingPanel';
import { MasterPlansPanel } from './panels/MasterPlansPanel';
import { IntegrationsPanel } from './panels/IntegrationsPanel';

interface PanelsProps {
  tab: string;
}

export const Panels = ({ tab }: PanelsProps) => {
  const isMasterPlansEnabled = useFeatureIsOn('settings-master-plans-view');

  switch (tab) {
    case 'auth':
      return <AuthPanel />;
    case 'billing':
      return <BillingPanel />;
    case 'integrations':
      return <IntegrationsPanel />;
    case 'master-plans':
      return isMasterPlansEnabled ? <MasterPlansPanel /> : <AuthPanel />;
    default:
      return <AuthPanel />;
  }
};
