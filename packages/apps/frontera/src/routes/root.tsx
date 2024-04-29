import type { RouteObject } from 'react-router-dom';

import { NotFound } from './not-found';

export const RootRoute: RouteObject = {
  path: '/',
  element: <div>Root route</div>,
  errorElement: <NotFound />,
};
