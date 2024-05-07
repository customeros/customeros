import { RouteObject } from 'react-router-dom';

import { SignIn } from './page';

export const SignInRoute: RouteObject = {
  path: '/auth/signin',
  element: <SignIn />,
};
