import { MouseEventHandler } from 'react';

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

import { cn } from '@ui/utils/cn.ts';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

import { StepViewportPortal } from './StepViewportPortal';

export const BasicEdge: React.FC<
  EdgeProps & { data: Record<string, boolean | string> }
> = observer(({ id, data, ...props }) => {
  const { setEdges } = useReactFlow();
  const [edgePath, labelX, labelY] = getSmoothStepPath({
    ...props,
  });
  const { ui } = useStore();

  const selected = props.selected;

  const toggleOpen: MouseEventHandler<HTMLButtonElement> = (e) => {
    e.stopPropagation();

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
        interactionWidth={80}
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
            icon={
              <Plus
                className='text-inherit transition-transform duration-100'
                style={{
                  transform:
                    ui.flowCommandMenu.isOpen &&
                    id === ui.flowCommandMenu.context.id
                      ? 'rotate(45deg)'
                      : 'initial',
                }}
              />
            }
            className={cn(
              'text-white bg-gray-300 hover:bg-gray-700 hover:text-white focus:bg-gray-700 focus:text-white rounded-full scale-50 hover:scale-100 transition-all ease-in-out ',
              {
                'text-transparent':
                  !data?.isHovered &&
                  (!ui.flowCommandMenu.isOpen ||
                    id !== ui.flowCommandMenu.context.id),

                'scale-100 bg-gray-700':
                  data?.isHovered ||
                  (ui.flowCommandMenu.isOpen &&
                    id === ui.flowCommandMenu.context.id),
              },
            )}
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
});
