import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const SequenceHub = observer(() => {
  const store = useStore();

  return (
    <CommandsContainer label={'Sequences'}>
      <CommandItem
        leftAccessory={<PlusCircle />}
        onSelect={() => {
          store.ui.commandMenu.setType('CreateNewSequence');
        }}
      >
        Add new sequence...
      </CommandItem>
    </CommandsContainer>
  );
});
