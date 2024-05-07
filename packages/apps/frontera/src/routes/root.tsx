import { RouteObject } from 'react-router-dom';

import { Layout } from '@shared/components/Layout/Layout';

import { NotFound } from './not-found';
import { AuthRoute } from './auth/route';
import { RenewalsRoute } from './renewals/route';
import { SettingsRoute } from './settings/route';
import { OrganizationRoute } from './organization/route';
import { OrganizationsRoute } from './organizations/route';

export const RootRoute: RouteObject = {
  path: '/',
  element: <Layout />,
  children: [
    AuthRoute,
    RenewalsRoute,
    SettingsRoute,
    OrganizationRoute,
    OrganizationsRoute,
  ],
  errorElement: <NotFound />,
};
