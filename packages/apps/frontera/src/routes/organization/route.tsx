import { RouteObject } from 'react-router-dom';

import NotFound from '@organization/components/NotFound/NotFound.tsx';

import { OrganizationPage } from './page';

export const OrganizationRoute: RouteObject = {
  path: '/organization/:id',
  element: <OrganizationPage />,
  errorElement: <NotFound />,
};
