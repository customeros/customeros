import { RouteObject } from 'react-router-dom';

import { OrganizationsPage } from './page';

export const OrganizationsRoute: RouteObject = {
  path: '/organizations',
  element: <OrganizationsPage />,
};
