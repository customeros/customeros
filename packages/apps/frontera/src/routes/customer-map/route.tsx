import { RouteObject } from 'react-router-dom';

import { DashboardPage } from './page';

export const CustomerMapRoute: RouteObject = {
  path: '/customer-map',
  element: <DashboardPage />,
};
