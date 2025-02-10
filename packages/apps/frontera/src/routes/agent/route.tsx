import { RouteObject } from 'react-router-dom';

import { AgentPage } from './page';
import { AgentSubRoutesWrapper } from './components';

export const AgentRoute: RouteObject = {
  path: '/agents/:id/*',
  element: <AgentPage />,
  children: [
    {
      path: '*',
      element: <AgentSubRoutesWrapper />,
    },
  ],
};
