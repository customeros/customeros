import { RouteObject } from 'react-router-dom';

import { FlowEditor } from './page';

export const FlowEditorRoute: RouteObject = {
  path: '/flow-editor',
  element: <FlowEditor />,
};
