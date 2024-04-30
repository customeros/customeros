import { Outlet, useLocation } from 'react-router-dom';

import { SettingsSidenav } from '@settings/components/SettingsSidenav';

import { PageLayout } from '@shared/components/PageLayout';
import { RootSidenav } from '@shared/components/RootSidenav/RootSidenav';
import { OrganizationSidenav } from '@organization/components/OrganizationSidenav';

export const Layout = () => {
  const location = useLocation();

  location.pathname;

  const menuRootSelect =
    location.pathname.startsWith('/organizations') ||
    location.pathname.startsWith('/renewals') ? (
      <RootSidenav />
    ) : location.pathname.startsWith('/organization/') ? (
      <OrganizationSidenav />
    ) : (
      <SettingsSidenav />
    );

  return (
    <PageLayout>
      {menuRootSelect}
      <div className='h-full overflow-hidden flex'>
        <Outlet />
      </div>
    </PageLayout>
  );
};
