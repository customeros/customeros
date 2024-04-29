import { RouteObject } from 'react-router-dom';

import { SignIn } from './signin';

export const SignInRoute: RouteObject = {
  path: '/auth/signin',
  element: <SignIn />,
};
