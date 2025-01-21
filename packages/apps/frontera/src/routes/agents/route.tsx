import { RouteObject } from 'react-router-dom';

import { AgentsPage } from './page';

export const AgentsRoute: RouteObject = {
  path: '/agents',
  element: <AgentsPage />,
};
