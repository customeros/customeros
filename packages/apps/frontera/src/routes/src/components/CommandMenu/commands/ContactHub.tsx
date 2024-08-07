import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

import { GlobalSharedCommands } from './GlobalHub';

export const ContactHub = observer(() => {
  const store = useStore();

  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.value.name}`;

  return (
    <Command>
      <CommandInput label={label} placeholder='Type a command or search' />
      <Command.List>
        <CommandItem
          leftAccessory={<Tag01 />}
          onSelect={() => {
            store.ui.commandMenu.setType('EditPersonaTag');
          }}
        >
          Edit persona tag...
        </CommandItem>

        <Command.Group heading='Navigate'>
          <GlobalSharedCommands />
        </Command.Group>
      </Command.List>
    </Command>
  );
});
