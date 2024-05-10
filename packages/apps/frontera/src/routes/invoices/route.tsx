import { RouteObject } from 'react-router-dom';

import { InvoicesPage } from './page';

export const InvoicesRoute: RouteObject = {
  path: '/invoices',
  element: <InvoicesPage />,
};
