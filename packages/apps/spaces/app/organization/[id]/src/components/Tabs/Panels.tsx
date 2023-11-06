import { AboutPanel } from './panels/AboutPanel';
import { UpNextPanel } from './panels/UpNextPanel';
import { PeoplePanel } from './panels/PeoplePanel';
import { IssuesPanel } from './panels/IssuesPanel';
import { AccountPanel } from './panels/AccountPanel';
import { SuccessPanel } from './panels/SuccessPanel';

interface PanelsProps {
  tab: string;
}

export const Panels = ({ tab }: PanelsProps) => {
  switch (tab) {
    case 'up-next':
      return <UpNextPanel />;
    case 'account':
      return <AccountPanel />;
    case 'success':
      return <SuccessPanel />;
    case 'people':
      return <PeoplePanel />;
    case 'issues':
      return <IssuesPanel />;
    default:
      return <AboutPanel />;
  }
};
