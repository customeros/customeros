import { RouteObject } from 'react-router-dom';

import { Auth } from './src/Auth';
import { SignInRoute } from './signin/route';
import { SuccessRoute } from './success/route';
import { FailureRoute } from './failure/route';

export const AuthRoute: RouteObject = {
  path: '/auth',
  element: <Auth />,
  children: [SignInRoute, SuccessRoute, FailureRoute],
};
