import { MouseEventHandler } from 'react';

import { observer } from 'mobx-react-lite';
import { NodeProps, ViewportPortal } from '@xyflow/react';

import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { UserPlus01 } from '@ui/media/icons/UserPlus01';
import { Lightning01 } from '@ui/media/icons/Lightning01.tsx';

import { Handle } from '../components';
import { DropdownCommandMenu } from '../commands/Commands.tsx';

// const triggerEventMapper: Record<string, string> = {
//   RecordAddedManually: 'added manually',
// };

export const TriggerNode = (
  props: NodeProps & { data: Record<string, string> },
) => {
  const { ui } = useStore();

  const handleOpen: MouseEventHandler<HTMLButtonElement> = (e) => {
    e.stopPropagation();
    ui.flowCommandMenu.setOpen(true);
    ui.flowCommandMenu.setType('TriggersHub');
    ui.flowCommandMenu.setContext({
      id: props.id,
      entity: 'Trigger',
    });
  };

  return (
    <>
      <div
        className={`aspect-[9/1] h-[56px] w-[300px] bg-white border border-grayModern-300 p-3 rounded-lg group relative cursor-pointer`}
      >
        {/*<div className='flex items-center text-gray-400 uppercase text-xs mb-1 absolute top-[-20px]'>*/}
        {/*  Trigger*/}
        {/*</div>*/}

        <div className='truncate  text-sm flex items-center justify-between '>
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
                added manually
              </span>
            ) : (
              <span role={'button'} onClick={handleOpen}>
                What should trigger this flow?
              </span>
            )}
          </div>

          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Edit'
            icon={<Edit03 />}
            onClick={handleOpen}
            className='ml-2 opacity-0 group-hover:opacity-100 pointer-events-all'
          />
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

    return (
      <>
        {ui.flowCommandMenu?.isOpen && id === ui.flowCommandMenu.context.id && (
          <ViewportPortal>
            <div
              style={{
                transform: `translate(calc(${positionAbsoluteX}px + 150px - 180px), ${
                  positionAbsoluteY + 70
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
