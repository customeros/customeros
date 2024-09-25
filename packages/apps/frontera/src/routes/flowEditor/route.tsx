import { RouteObject } from 'react-router-dom';

import { NotFound } from '@shared/components/NotFound';

import { FlowEditor } from './page';

export const FlowEditorRoute: RouteObject = {
  path: '/flow-editor/:id',
  element: <FlowEditor />,
  errorElement: <NotFound />,
};
