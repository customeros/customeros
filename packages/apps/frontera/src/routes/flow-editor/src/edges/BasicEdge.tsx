import { MouseEventHandler } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import {
  BaseEdge,
  EdgeProps,
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

      return;
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
      <BaseEdge {...props} path={edgePath} />

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
            dataTest={'flow-add-step-or-trigger'}
            icon={
              <Plus
                strokeWidth={4}
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
              'bg-gray-300 border-4 border-white text-transparent hover:bg-gray-700 hover:text-white focus:bg-inherit focus:text-inherit  rounded-full scale-[0.3635] transition-all ease-in-out ',
              {
                'scale-100 !bg-gray-700 text-white border-2':
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
