import { Outlet } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Header } from './components';
import { SubHeader } from './components/SubHeader/SubHeader';

export const AgentPage = observer(() => {
  return (
    <>
      <Header />
      <SubHeader />
      <Outlet />
    </>
  );
});
