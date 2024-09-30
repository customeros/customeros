import { observer } from 'mobx-react-lite';
import { ViewportPortal } from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';

import { DropdownCommandMenu } from '../commands/Commands.tsx';

export const StepViewportPortal = observer(
  ({
    id,
    positionAbsoluteX,
    positionAbsoluteY,
  }: {
    id: string;
    positionAbsoluteX: number;
    positionAbsoluteY: number;
  }) => {
    const { ui } = useStore();

    return (
      <>
        {ui.flowCommandMenu?.isOpen && id === ui.flowCommandMenu.context.id && (
          <ViewportPortal>
            <div
              className='border border-gray-200 rounded-lg shadow-lg'
              style={{
                transform: `translate(calc(${positionAbsoluteX}px - 50%), ${
                  positionAbsoluteY + 24 // 24 is desired spacing between dropdown and button
                }px)`,
                position: 'absolute',
                pointerEvents: 'all',
                zIndex: 50000,
                width: '360px',
                left: '0',
                top: '0',
              }}
            >
              <DropdownCommandMenu />
            </div>
          </ViewportPortal>
        )}
      </>
    );
  },
);
