import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';
import { User01 } from '@ui/media/icons/User01.tsx';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const ContactAddedManuallySubItem = observer(() => {
  const { ui } = useStore();
  const { setNodes } = useReactFlow();

  const updateSelectedNode = () => {
    setNodes((nodes) =>
      nodes.map((node) => {
        if (node.id === ui.flowCommandMenu.context.id) {
          return {
            ...node,
            data: {
              ...node.data,
              triggerType: 'RecordAddedManually',
              entity: 'CONTACT',
            },
          };
        }

        return node;
      }),
    );
  };

  return (
    <>
      <CommandSubItem
        icon={<User01 />}
        leftLabel={'Contact'}
        rightLabel={'added manually'}
        onSelectAction={() => {
          updateSelectedNode();
          ui.flowCommandMenu.setOpen(false);
          ui.flowCommandMenu.setType('TriggersHub');
        }}
      />
    </>
  );
});
