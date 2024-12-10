import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { NodeProps, ViewportPortal } from '@xyflow/react';

import { FlowStatus } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { UserPlus01 } from '@ui/media/icons/UserPlus01';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { Lightning01 } from '@ui/media/icons/Lightning01';

import { Handle } from '../components';
import { DropdownCommandMenu } from '../commands/Commands';

export const TriggerNode = (
  props: NodeProps & { data: Record<string, string> },
) => {
  const { ui, flows } = useStore();
  const flowId = useParams()?.id as string;
  const flow = flows.value.get(flowId)?.value;
  const flowWasStarted = flow?.status === FlowStatus.On || flow?.firstStartedAt;

  return (
    <>
      <div
        className={`h-[83px] w-[300px] bg-white border border-grayModern-300 rounded-lg group relative cursor-pointer flex flex-col items-center`}
      >
        <div
          data-test={'flow-trigger-block'}
          className='px-4 bg-gray-25 text-xs h-full flex items-center w-full rounded-t-lg justify-center border-b border-dashed border-gray-300 text-gray-500'
        >
          Flow triggers when
        </div>

        <div className='flex items-center justify-between w-full p-4  h-[56px]'>
          <div className='truncate text-sm flex items-center'>
            <div className='size-6 mr-2 bg-gray-50 border border-gray-100 rounded flex items-center justify-center'>
              {props.data.entity && props.data?.triggerType ? (
                <UserPlus01 className='text-gray-500 ' />
              ) : (
                <Lightning01 className='text-gray-500' />
              )}
            </div>

            {props.data.entity && props.data.triggerType ? (
              <span className='font-medium '>
                <span className='capitalize mr-1'>
                  {props.data.entity?.toLowerCase() ?? 'Record'}
                </span>
                is added to this flow
              </span>
            ) : (
              <span className='text-gray-400'>Choose a trigger…</span>
            )}
          </div>

          {!flowWasStarted && !ui.flowActionSidePanel.isOpen && (
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='Edit'
              icon={<ChevronDown />}
              dataTest={'flow-trigger-block-options'}
              className='ml-2 opacity-0 group-hover:opacity-100 pointer-events-all'
            />
          )}
        </div>
        <Handle type='target' />
        <Handle type='source' />
      </div>
      <TriggerViewportPortal
        id={props.id}
        positionAbsoluteX={props.positionAbsoluteX}
        positionAbsoluteY={props.positionAbsoluteY}
      />
    </>
  );
};

export const TriggerViewportPortal = observer(
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

    const showTriggerDropdown =
      ui.flowCommandMenu?.isOpen &&
      id === ui.flowCommandMenu.context.id &&
      ui.flowCommandMenu.context.entity === 'Trigger';

    return (
      <>
        {showTriggerDropdown && (
          <ViewportPortal>
            <div
              className='border border-gray-200 rounded-lg shadow-lg'
              style={{
                transform: `translate(calc(${positionAbsoluteX}px + 150px - 180px), ${
                  positionAbsoluteY + 83 + 4 // 83 is height of the node, 4 is desired spacing
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
