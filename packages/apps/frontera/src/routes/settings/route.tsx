import { RouteObject } from 'react-router-dom';

import { SettingsPage } from './page';

export const SettingsRoute: RouteObject = {
  path: '/settings',
  element: <SettingsPage />,
};
