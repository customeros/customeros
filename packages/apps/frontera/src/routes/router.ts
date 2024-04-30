import { createBrowserRouter } from 'react-router-dom';

import { RootRoute } from './root';
import { AuthRoute } from './auth/route';
// import { SettingsRoute } from './settings/route';
import { OrganizationRoute } from './organization/route';
import { OrganizationsRoute } from './organizations/route';

export const router = createBrowserRouter([
  RootRoute,
  AuthRoute,
  // SettingsRoute,
  OrganizationRoute,
  OrganizationsRoute,
]);
