import { RouteObject } from 'react-router-dom';

import { Layout } from '@shared/components/Layout/Layout';

import { AuthRoute } from './auth/route';
import { AgentRoute } from './agent/route';
import { AgentsRoute } from './agents/route';
import { FinderRoute } from './finder/route';
import { Error } from './src/components/Error';
import { SettingsRoute } from './settings/route';
import { ProspectsRoute } from './prospects/route';
import { NotFound } from './src/components/NotFound';
import { OnboardingRoute } from './onboarding/route';
import { FlowEditorRoute } from './flow-editor/route';
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
    AgentRoute,
    AgentsRoute,
    SettingsRoute,
    OrganizationRoute,
    FinderRoute,
    CustomerMapRoute,
    ProspectsRoute,
    FlowEditorRoute,
    OnboardingRoute,
    NotFoundRoute,
  ],
  errorElement: <Error />,
};
