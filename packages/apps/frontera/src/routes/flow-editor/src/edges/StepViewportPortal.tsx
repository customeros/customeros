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
              style={{
                transform: `translate(${positionAbsoluteX / 2}px,${
                  positionAbsoluteY + 15
                }px)`,
                position: 'absolute',
                pointerEvents: 'all',
                zIndex: 50000,
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
