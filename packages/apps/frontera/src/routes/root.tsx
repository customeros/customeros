import { RouteObject } from 'react-router-dom';

import { Layout } from '@shared/components/Layout/Layout';

import { AuthRoute } from './auth/route';
import { FinderRoute } from './finder/route';
import { Error } from './src/components/Error';
import { SettingsRoute } from './settings/route';
import { ProspectsRoute } from './prospects/route';
import { NotFound } from './src/components/NotFound';
import { CustomerMapRoute } from './customer-map/route';
import { OrganizationRoute } from './organization/route';
import { FlowEditorRoute } from './flow-editor/route.tsx';

const NotFoundRoute: RouteObject = {
  path: '*',
  element: <NotFound />,
};

export const RootRoute: RouteObject = {
  path: '/',
  element: <Layout />,
  children: [
    AuthRoute,
    SettingsRoute,
    OrganizationRoute,
    FinderRoute,
    CustomerMapRoute,
    ProspectsRoute,
    FlowEditorRoute,
    NotFoundRoute,
  ],
  errorElement: <Error />,
};
