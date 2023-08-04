import {OAuthPanel} from "./panels/OAuthPanel";
import {BillingInfoPanel} from "./panels/BillingInfoPanel";


interface PanelsProps {
  tab: string;
}

export const Panels = ({ tab }: PanelsProps) => {
  switch (tab) {
    case 'oauth':
      return <OAuthPanel />;
    case 'billing':
      return <BillingInfoPanel />;
    default:
      return <OAuthPanel />;
  }
};
