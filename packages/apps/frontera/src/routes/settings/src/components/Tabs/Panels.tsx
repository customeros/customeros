import { AuthPanel } from './panels/AuthPanel';
import { General } from './panels/Workspace/General';
import { BillingPanel } from './panels/BillingPanel';
import { ApiManager } from './panels/Workspace/ApiManager';
import { TagsManager } from './panels/Workspace/TagsManager';
import { IntegrationsPanel } from './panels/IntegrationsPanel';
import { OrganizationFields } from './panels/Fields/Organizations';

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
    case 'api':
      return <ApiManager />;
    case 'organizations':
      return <OrganizationFields />;
    // case 'contacts':
    //   return <ContactFields />;

    default:
      return <AuthPanel />;
  }
};
