import { Controls } from '@xyflow/react';

import { LayoutButton } from './LayoutButton.tsx';

export const Toolbar = () => {
  return (
    <Controls
      className='px-1 py-1'
      position='bottom-left'
      orientation='horizontal'
    >
      <LayoutButton />
    </Controls>
  );
};
