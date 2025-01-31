import { useParams } from 'react-router-dom';

import { OrganizationDetails } from '@shared/components/OrganizationDetails';

import { PeoplePanel } from './panels/PeoplePanel';
import { IssuesPanel } from './panels/IssuesPanel';
import { AccountPanel } from './panels/AccountPanel';
import { SuccessPanel } from './panels/SuccessPanel';
import { InvoicesPanel } from './panels/InvoicesPanel';

interface PanelsProps {
  tab: string;
}

export const Panels = ({ tab }: PanelsProps) => {
  const { id } = useParams<{ id: string }>();

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
      return id ? <OrganizationDetails id={id} /> : null;
  }
};
