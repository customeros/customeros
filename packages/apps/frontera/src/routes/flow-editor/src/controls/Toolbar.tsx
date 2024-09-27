import { Controls } from '@xyflow/react';

import { LayoutButton } from './LayoutButton.tsx';

export const Toolbar = () => {
  return (
    <Controls
      position='bottom-left'
      showInteractive={false}
      orientation='horizontal'
      className='bg-white rounded'
    >
      <LayoutButton />
    </Controls>
  );
};
