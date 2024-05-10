import { RouteObject } from 'react-router-dom';

import { FailurePage } from './page';

export const FailureRoute: RouteObject = {
  path: '/auth/failure',
  element: <FailurePage />,
};
