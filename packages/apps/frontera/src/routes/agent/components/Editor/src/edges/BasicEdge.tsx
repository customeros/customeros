import { useParams } from 'react-router-dom';
import { FC, useState, useEffect, MouseEventHandler } from 'react';

import { observer } from 'mobx-react-lite';
import {
  BaseEdge,
  EdgeProps,
  useReactFlow,
  EdgeLabelRenderer,
  getSmoothStepPath,
} from '@xyflow/react';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { FlowStatus } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

import { StepViewportPortal } from './StepViewportPortal';

export const BasicEdge: FC<
  EdgeProps & { data: Record<string, boolean | string> }
> = observer(({ id, data, ...props }) => {
  const [edgePath, labelX, labelY] = getSmoothStepPath({
    ...props,
  });

  const [showButton, setShowButton] = useState(false);
  const { getNodes } = useReactFlow();

  useEffect(() => {
    if (getNodes().length === 2) setShowButton(true);
  }, []);

  const { ui, flows } = useStore();
  const flowId = useParams()?.id as string;
  const flow = flows.value.get(flowId)?.value;
  const flowWasStarted = flow?.status === FlowStatus.On || flow?.firstStartedAt;

  const toggleOpen: MouseEventHandler<HTMLButtonElement> = (e) => {
    // do not open when  flow actions panel is open
    if (ui.flowActionSidePanel.isOpen) {
      return;
    }
    e.stopPropagation();

    if (flowWasStarted) {
      ui.commandMenu.setType('ActiveFlowUpdateInfo');
      ui.commandMenu.setContext({
        ...ui.commandMenu.context,
        meta: {
          type: 'steps',
        },
      });
      ui.commandMenu.setOpen(true);

      return;
    }

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

  return (
    <>
      <BaseEdge {...props} path={edgePath} />

      <EdgeLabelRenderer>
        <div
          className={cn('nodrag nopan')}
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
              'bg-gray-300 border-4 border-white text-transparent hover:bg-gray-700 hover:text-white focus:bg-inherit focus:text-inherit rounded-full scale-[0.3635] transition-all ease-in-out ',
              {
                'scale-100 !bg-gray-700 text-white border-2':
                  showButton ||
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
