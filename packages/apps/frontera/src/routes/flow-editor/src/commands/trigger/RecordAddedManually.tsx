import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';

import { User01 } from '@ui/media/icons/User01';
import { useStore } from '@shared/hooks/useStore';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { Building07 } from '@ui/media/icons/Building07';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';

export const RecordAddedManually = observer(() => {
  const { ui } = useStore();
  const { setNodes } = useReactFlow();

  const updateSelectedNode = (entity: 'CONTACT') => {
    setNodes((nodes) =>
      nodes.map((node) => {
        if (node.id === ui.flowCommandMenu.context.id) {
          return {
            ...node,
            data: {
              ...node.data,
              entity,
            },
          };
        }

        return node;
      }),
    );
  };

  return (
    <>
      <CommandItem
        leftAccessory={<User01 />}
        keywords={keywords.contact}
        onSelect={() => {
          updateSelectedNode('CONTACT');
          ui.flowCommandMenu.setOpen(false);
          ui.flowCommandMenu.setType('TriggersHub');
        }}
      >
        Contact
      </CommandItem>
      <CommandItem
        disabled
        leftAccessory={<Building07 />}
        keywords={keywords.organization}
      >
        <span className='text-gray-700'>Organization</span>{' '}
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>{' '}
      <CommandItem
        disabled
        keywords={keywords.opportunity}
        leftAccessory={<CoinsStacked01 />}
      >
        <span className='text-gray-700'>Opportunity</span>{' '}
        <span className='text-gray-500'>(Coming soon)</span>
      </CommandItem>
    </>
  );
});
const keywords = {
  contact: ['contact', 'people'],
  organization: ['organization', 'company'],
  opportunity: ['opportunity', 'deal'],
};
