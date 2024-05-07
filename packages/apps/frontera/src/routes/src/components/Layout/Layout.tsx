import { useEffect } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

import { P, match } from 'ts-pattern';
import { SettingsSidenav } from '@settings/components/SettingsSidenav';

import { PageLayout } from '@shared/components/PageLayout';
import { RootSidenav } from '@shared/components/RootSidenav/RootSidenav';
import { OrganizationSidenav } from '@organization/components/OrganizationSidenav';

export const Layout = () => {
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (location.pathname === '/') {
      navigate('/organizations');
    }
  }, [location.pathname]);

  if (location.pathname === '/') {
    return null;
  }
  const sidenav = match(location.pathname)
    .with(
      P.string.startsWith('/organizations'),
      P.string.startsWith('/renewals'),
      P.string.startsWith('/invoices'),
      () => <RootSidenav />,
    )
    .with(P.string.startsWith('/organization'), () => <OrganizationSidenav />)
    .with(P.string.startsWith('/settings'), () => <SettingsSidenav />)
    .otherwise(() => null);

  return (
    <PageLayout
      unstyled={location.pathname.startsWith('/auth')}
      className='w-screen h-screen'
    >
      {sidenav}
      <div className='h-full w-full overflow-hidden flex'>
        <Outlet />
      </div>
    </PageLayout>
  );
};
