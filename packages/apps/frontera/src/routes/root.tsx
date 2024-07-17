import { RouteObject } from 'react-router-dom';

import { Layout } from '@shared/components/Layout/Layout';

import { AuthRoute } from './auth/route';
import { Error } from './src/components/Error';
import { RenewalsRoute } from './renewals/route';
import { SettingsRoute } from './settings/route';
import { ProspectsRoute } from './prospects/route';
import { FinderRoute } from './organizations/route';
import { NotFound } from './src/components/NotFound';
import { CustomerMapRoute } from './customer-map/route';
import { OrganizationRoute } from './organization/route';

const NotFoundRoute: RouteObject = {
  path: '*',
  element: <NotFound />,
};

export const RootRoute: RouteObject = {
  path: '/',
  element: <Layout />,
  children: [
    AuthRoute,
    RenewalsRoute,
    SettingsRoute,
    OrganizationRoute,
    FinderRoute,
    CustomerMapRoute,
    ProspectsRoute,
    NotFoundRoute,
  ],
  errorElement: <Error />,
};
