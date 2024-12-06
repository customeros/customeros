import { RouteObject } from 'react-router-dom';

import { NotFound } from '@shared/components/NotFound';

import { OnboardingPage } from './page.tsx';

export const OnboardingRoute: RouteObject = {
  path: '/welcome',
  element: <OnboardingPage />,
  errorElement: <NotFound />,
};
