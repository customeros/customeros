import { ReactElement } from 'react';

import { observer } from 'mobx-react-lite';
import { NodeProps, useReactFlow, ViewportPortal } from '@xyflow/react';

import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { Code01 } from '@ui/media/icons/Code01.tsx';
import { User01 } from '@ui/media/icons/User01.tsx';
import { XSquare } from '@ui/media/icons/XSquare.tsx';
import { UserPlus01 } from '@ui/media/icons/UserPlus01';
import { Building07 } from '@ui/media/icons/Building07.tsx';
import { PlusSquare } from '@ui/media/icons/PlusSquare.tsx';
import { PlusCircle } from '@ui/media/icons/PlusCircle.tsx';
import { Lightning01 } from '@ui/media/icons/Lightning01.tsx';
import { RefreshCw01 } from '@ui/media/icons/RefreshCw01.tsx';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01.tsx';
import { CheckCircleBroken } from '@ui/media/icons/CheckCircleBroken.tsx';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

import { Handle } from '../components';

const triggerEventMapper: Record<string, string> = {
  RecordAddedManually: 'added manually',
};

export const TriggerNode = (
  props: NodeProps & { data: Record<string, string> },
) => {
  const { ui } = useStore();

  if (props.data.triggerType === 'EndFlow') {
    return (
      <div className='max-w-[131px] flex bg-white border-2 border-grayModern-300 p-3 rounded-lg items-center'>
        <div className='size-6 mr-2 bg-gray-100 rounded flex items-center justify-center'>
          <XSquare className='text-gray-500' />
        </div>
        <span className='text-sm'>End Flow</span>
        <Handle
          type='target'
          className={`h-2 w-2 bg-transparent border-transparent`}
        />
      </div>
    );
  }

  const handleOpen = () => {
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
        className={`aspect-[9/1] max-w-[300px] bg-white border-2 border-grayModern-300 p-3 rounded-lg group`}
      >
        <div className='flex items-center text-gray-400 uppercase text-xs mb-1'>
          Trigger
        </div>

        <div className='truncate  text-sm flex items-center '>
          <div className='truncate text-sm flex items-center'>
            <div className='size-6 mr-2 bg-gray-100 rounded flex items-center justify-center'>
              {props.data.triggerEntity && props.data.triggerType ? (
                <UserPlus01 className='text-gray-500 ' />
              ) : (
                <Lightning01 className='text-gray-500' />
              )}
            </div>

            {props.data.triggerEntity && props.data.triggerType ? (
              <span className='font-medium'>
                {props.data.triggerEntity ?? 'Record'}{' '}
                {triggerEventMapper?.[props.data.triggerType]}
              </span>
            ) : (
              <span className='text-gray-400'>Select trigger...</span>
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
                transform: `translate(${positionAbsoluteX + 15}px,${
                  positionAbsoluteY + 90
                }px)`,
                position: 'absolute',
                pointerEvents: 'all',
              }}
            >
              <TriggerMenu />
            </div>
          </ViewportPortal>
        )}
      </>
    );
  },
);

const RecordAddedManually = observer(() => {
  const { ui } = useStore();
  const { setNodes } = useReactFlow();

  const updateSelectedNode = (entity: 'Contact') => {
    setNodes((nodes) =>
      nodes.map((node) => {
        if (node.id === ui.flowCommandMenu.context.id) {
          return {
            ...node,
            data: {
              ...node.data,
              triggerEntity: entity,
            },
          };
        }

        return node;
      }),
    );
  };

  return (
    <Command>
      <CommandInput
        autoFocus
        className='p-1 text-sm'
        placeholder='Search record'
        inputWrapperClassName='min-h-4'
        wrapperClassName='py-2 px-4 mt-2'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />

      <Command.List>
        <CommandItem
          leftAccessory={<User01 />}
          onSelect={() => {
            updateSelectedNode('Contact');
            ui.flowCommandMenu.setOpen(false);
            ui.flowCommandMenu.setType('TriggersHub');
          }}
        >
          Contact
        </CommandItem>
        <CommandItem disabled leftAccessory={<Building07 />}>
          <span className='text-gray-700'>Organization</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>{' '}
        <CommandItem disabled leftAccessory={<CoinsStacked01 />}>
          <span className='text-gray-700'>Opportunity</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>
      </Command.List>
    </Command>
  );
});

const RecordCreated = () => (
  <div>
    <h2>Record Created</h2>
    <p>This feature is coming soon.</p>
  </div>
);

const RecordUpdated = () => (
  <div>
    <h2>Record Updated</h2>
    <p>This feature is coming soon.</p>
  </div>
);

const RecordMatchesCondition = () => (
  <div>
    <h2>Record Matches Condition</h2>
    <p>This feature is coming soon.</p>
  </div>
);

const Webhook = () => (
  <div>
    <h2>Webhook</h2>
    <p>This feature is coming soon.</p>
  </div>
);

export const TriggersHub = observer(() => {
  const { ui } = useStore();
  const { setNodes } = useReactFlow();

  const updateSelectedNode = (triggerType: 'RecordAddedManually') => {
    setNodes((nodes) =>
      nodes.map((node) => {
        if (node.id === ui.flowCommandMenu.context.id) {
          return {
            ...node,
            data: {
              ...node.data,

              triggerType: triggerType,
            },
          };
        }

        return node;
      }),
    );
    ui.flowCommandMenu.setType(triggerType);
  };

  return (
    <Command>
      <CommandInput
        autoFocus
        className='p-1 text-sm'
        placeholder='Search trigger'
        inputWrapperClassName='min-h-4'
        wrapperClassName='py-2 px-4 mt-2'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />

      <Command.List>
        <CommandItem
          leftAccessory={<PlusCircle />}
          onSelect={() => {
            updateSelectedNode('RecordAddedManually');
          }}
        >
          Record added manually...
        </CommandItem>
        <CommandItem disabled leftAccessory={<PlusSquare />}>
          <span className='text-gray-700'>Record created</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>
        <CommandItem disabled leftAccessory={<RefreshCw01 />}>
          <span className='text-gray-700'>Record updated</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>
        <CommandItem disabled leftAccessory={<CheckCircleBroken />}>
          <span className='text-gray-700'>Record matches condition</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>
        <CommandItem disabled leftAccessory={<Code01 />}>
          <span className='text-gray-700'>Webhook</span>{' '}
          <span className='text-gray-500'>(Coming soon)</span>
        </CommandItem>
      </Command.List>
    </Command>
  );
});

const Commands: Record<string, ReactElement> = {
  RecordAddedManually: <RecordAddedManually />,
  RecordCreated: <RecordCreated />,
  RecordUpdated: <RecordUpdated />,
  RecordMatchesCondition: <RecordMatchesCondition />,
  Webhook: <Webhook />,
  TriggersHub: <TriggersHub />,
};

const TriggerMenu = observer(() => {
  const { ui } = useStore();

  return <>{Commands[ui.flowCommandMenu.type]}</>;
});
