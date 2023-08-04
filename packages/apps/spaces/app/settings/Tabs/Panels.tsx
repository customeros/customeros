
import {BillingPanel} from "./panels/BillingPanel";
import {AuthPanel} from "./panels/AuthPanel/AuthPanel";


interface PanelsProps {
  tab: string;
}

export const Panels = ({ tab }: PanelsProps) => {
  switch (tab) {
    case 'auth':
      return <AuthPanel />;
    case 'billing':
      return <BillingPanel />;
    default:
      return <AuthPanel />;
  }
};
