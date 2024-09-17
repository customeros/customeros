import { RouteObject } from 'react-router-dom';

import { MarketingFlowBuilder } from './page';

export const FlowEditorRoute: RouteObject = {
  path: '/flow-editor',
  element: <MarketingFlowBuilder />,
};
