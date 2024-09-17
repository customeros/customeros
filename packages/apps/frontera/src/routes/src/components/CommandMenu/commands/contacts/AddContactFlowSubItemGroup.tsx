import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store.ts';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Shuffle01 } from '@ui/media/icons/Shuffle01';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const AddContactFlowSubItemGroup = observer(() => {
  const { contacts, ui, flows } = useStore();

  const context = ui.commandMenu.context;

  const contact = contacts.value.get(context.ids?.[0] as string);
  const selectedIds = context.ids;

  const handleSelect = (opt: FlowStore) => {
    if (!context.ids?.[0] || !contact) return;

    if (selectedIds?.length === 1) {
      opt.linkContact(contact.id);
    }

    if (selectedIds?.length > 1) {
      opt.linkContacts(selectedIds);
    }

    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const sequenceOptions = flows.toComputedArray((arr) => arr);

  return (
    <>
      {sequenceOptions.map((flowSequence) => {
        const isSelected =
          context.ids?.length === 1 && contact?.flow?.id === flowSequence.id;

        return (
          <CommandSubItem
            icon={<Shuffle01 />}
            key={flowSequence.id}
            leftLabel='Move to flow'
            rightLabel={flowSequence.value.name ?? 'Unnamed'}
            rightAccessory={isSelected ? <Check /> : undefined}
            onSelectAction={() => {
              handleSelect(flowSequence as FlowStore);
            }}
          />
        );
      })}
    </>
  );
});
