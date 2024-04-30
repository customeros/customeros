import { RouteObject } from 'react-router-dom';

import { Layout } from '@shared/components/Layout/Layout';

import { NotFound } from './not-found';
import { RenewalsRoute } from './renewals/route';
import { SettingsRoute } from './settings/route';
import { OrganizationRoute } from './organization/route';
import { OrganizationsRoute } from './organizations/route';

export const RootRoute: RouteObject = {
  path: '/',
  element: <Layout />,
  children: [
    OrganizationsRoute,
    OrganizationRoute,
    RenewalsRoute,
    SettingsRoute,
  ],

  errorElement: <NotFound />,
};
