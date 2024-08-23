import { AuthPanel } from './panels/AuthPanel';
import { General } from './panels/Workspace/General';
import { BillingPanel } from './panels/BillingPanel';
import { TagsManager } from './panels/Workspace/TagsManager';
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
    case 'general':
      return <General />;
    case 'tags':
      return <TagsManager />;

    default:
      return <AuthPanel />;
  }
};
