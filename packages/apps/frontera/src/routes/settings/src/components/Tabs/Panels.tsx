import { AuthPanel } from './panels/AuthPanel';
import { Workspace } from './panels/Workspace';
import { BillingPanel } from './panels/BillingPanel';
import { IntegrationsPanel } from './panels/IntegrationsPanel';

interface PanelsProps {
  tab: string;
}

export const Panels = ({ tab }: PanelsProps) => {
  switch (tab) {
    case 'auth':
      return <AuthPanel />;
    case 'billing':
      return <BillingPanel />;
    case 'integrations':
      return <IntegrationsPanel />;
    case 'workspace':
      return <Workspace />;

    default:
      return <AuthPanel />;
  }
};
