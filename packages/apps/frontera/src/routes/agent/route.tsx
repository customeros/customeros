import { RouteObject } from 'react-router-dom';

import { AgentPage } from './page';

export const AgentRoute: RouteObject = {
  path: '/agents/:id',
  element: <AgentPage />,
};
