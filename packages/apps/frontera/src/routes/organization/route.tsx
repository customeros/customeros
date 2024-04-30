import { RouteObject } from 'react-router-dom';

import { OrganizationPage } from './page';

export const OrganizationRoute: RouteObject = {
  path: '/organization/:id',
  element: <OrganizationPage />,
};
