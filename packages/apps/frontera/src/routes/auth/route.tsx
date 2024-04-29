import { RouteObject } from 'react-router-dom';

import { Auth } from './src/Auth';
import { SignInRoute } from './signin/route';

export const AuthRoute: RouteObject = {
  path: '/auth',
  element: <Auth />,
  children: [SignInRoute],
};
