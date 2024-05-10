import { RouteObject } from 'react-router-dom';

import { ProspectsBoardPage } from './page';

export const ProspectsRoute: RouteObject = {
  path: '/prospects',
  element: <ProspectsBoardPage />,
};
