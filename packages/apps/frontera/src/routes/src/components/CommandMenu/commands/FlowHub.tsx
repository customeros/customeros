import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const FlowHub = observer(() => {
  const store = useStore();

  return (
    <CommandsContainer label={'Flows'}>
      <CommandItem
        leftAccessory={<PlusCircle />}
        onSelect={() => {
          store.ui.commandMenu.setType('CreateNewFlow');
        }}
      >
        Add new flow...
      </CommandItem>
    </CommandsContainer>
  );
});
