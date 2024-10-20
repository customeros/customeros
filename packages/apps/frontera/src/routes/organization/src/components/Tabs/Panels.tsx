import { AboutPanel } from './panels/AboutPanel';
import { PeoplePanel } from './panels/PeoplePanel';
import { IssuesPanel } from './panels/IssuesPanel';
import { AccountPanel } from './panels/AccountPanel';
import { SuccessPanel } from './panels/SuccessPanel';
import { InvoicesPanel } from './panels/InvoicesPanel';

interface PanelsProps {
  tab: string;
}

export const Panels = ({ tab }: PanelsProps) => {
  switch (tab) {
    case 'account':
      return <AccountPanel />;
    case 'success':
      return <SuccessPanel />;
    case 'people':
      return <PeoplePanel />;
    case 'issues':
      return <IssuesPanel />;
    case 'invoices':
      return <InvoicesPanel />;
    default:
      return <AboutPanel />;
  }
};
