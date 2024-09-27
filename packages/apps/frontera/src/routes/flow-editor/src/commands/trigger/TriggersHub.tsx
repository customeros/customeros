import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';

import { Code01 } from '@ui/media/icons/Code01';
import { useStore } from '@shared/hooks/useStore';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { PlusSquare } from '@ui/media/icons/PlusSquare';
import { RefreshCw01 } from '@ui/media/icons/RefreshCw01';
import { CheckCircleBroken } from '@ui/media/icons/CheckCircleBroken';

import { ContactAddedManuallySubItem } from './ContactTriggerSubItems.tsx';

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
    <>
      <CommandItem
        leftAccessory={<PlusCircle />}
        keywords={['record', 'added', 'manually']}
        onSelect={() => {
          updateSelectedNode('RecordAddedManually');
        }}
      >
        Record added manually...
      </CommandItem>

      <ContactAddedManuallySubItem />
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
    </>
  );
});
