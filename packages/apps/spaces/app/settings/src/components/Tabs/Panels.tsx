import { AuthPanel } from './panels/AuthPanel';
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
    default:
      return <AuthPanel />;
  }
};
