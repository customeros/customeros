import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { useStore } from '@shared/hooks/useStore';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const ContactHub = observer(() => {
  const store = useStore();

  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.value.name}`;

  return (
    <CommandsContainer label={label}>
      <CommandItem
        leftAccessory={<Tag01 />}
        onSelect={() => {
          store.ui.commandMenu.setType('EditPersonaTag');
        }}
      >
        Edit persona tag...
      </CommandItem>
    </CommandsContainer>
  );
});
