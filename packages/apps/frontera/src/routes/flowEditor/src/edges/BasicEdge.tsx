import { useState } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import {
  BaseEdge,
  EdgeProps,
  MarkerType,
  useReactFlow,
  getSmoothStepPath,
  EdgeLabelRenderer,
} from '@xyflow/react';

import { Plus } from '@ui/media/icons/Plus.tsx';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

import { StepViewportPortal } from './EdgeCommandMenu.tsx';

export const BasicEdge: React.FC<EdgeProps> = observer(
  ({ id, data, ...props }) => {
    const [isOpen] = useState(false);
    const { setEdges } = useReactFlow();
    const [edgePath, labelX, labelY] = getSmoothStepPath({
      ...props,
    });
    const { ui } = useStore();

    const selected = props.selected;

    const toggleOpen = () => {
      if (ui.flowCommandMenu.isOpen) {
        ui.flowCommandMenu.setOpen(false);
      }

      ui.flowCommandMenu.setType('StepsHub');
      ui.flowCommandMenu.setOpen(true);
      ui.flowCommandMenu.setContext({
        entity: 'Step',
        id,
        meta: {
          source: props.source,
          target: props.target,
        },
      });
    };

    const handleDisconnectSteps = () => {
      setEdges((edges) => edges.filter((edge) => edge.id !== id));
    };

    useKey('Backspace', handleDisconnectSteps, {
      when: selected,
    });

    return (
      <>
        <BaseEdge
          {...props}
          path={edgePath}
          markerEnd={MarkerType.ArrowClosed}
        />

        <EdgeLabelRenderer>
          <div
            className='nodrag nopan'
            style={{
              position: 'absolute',
              transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
              fontSize: 12,
              pointerEvents: 'all',
            }}
          >
            <IconButton
              size='xxs'
              onClick={toggleOpen}
              aria-label='Add step or trigger'
              className='text-white bg-gray-700 hover:bg-gray-600 hover:text-white focus:bg-gray-600 focus:text-white rounded-full'
              icon={
                <Plus
                  className='text-inherit transition-transform duration-100'
                  style={{ transform: isOpen ? 'rotate(45deg)' : 'initial' }}
                />
              }
            />
          </div>
        </EdgeLabelRenderer>

        <StepViewportPortal
          id={id}
          positionAbsoluteX={labelX}
          positionAbsoluteY={labelY}
        />
      </>
    );
  },
);
