import { Outlet, useLocation } from 'react-router-dom';

import { P, match } from 'ts-pattern';
import { SettingsSidenav } from '@settings/components/SettingsSidenav';

import { PageLayout } from '@shared/components/PageLayout';
import { RootSidenav } from '@shared/components/RootSidenav/RootSidenav';
import { OrganizationSidenav } from '@organization/components/OrganizationSidenav';

const allowedPaths = [
  '/auth',
  '/organizations/',
  '/organization',
  '/invoices',
  '/renewals',
  '/customer-map',
  '/settings',
  '/prospects',
];

export const Layout = () => {
  const location = useLocation();

  if (!allowedPaths.some((path) => location.pathname.startsWith(path))) {
    return null;
  }

  const sidenav = match(location.pathname)
    .with(
      P.string.startsWith('/organizations'),
      P.string.startsWith('/renewals'),
      P.string.startsWith('/invoices'),
      P.string.startsWith('/prospects'),
      P.string.startsWith('/customer-map'),
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
      <div className='h-full w-full flex-col overflow-hidden flex'>
        <Outlet />
      </div>
    </PageLayout>
  );
};
