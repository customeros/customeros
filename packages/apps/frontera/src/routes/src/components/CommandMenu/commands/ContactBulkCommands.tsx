import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Clock } from '@ui/media/icons/Clock';
import { useStore } from '@shared/hooks/useStore';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import { GlobalSharedCommands } from '@shared/components/CommandMenu/commands/GlobalHub.tsx';

export const ContactBulkCommands = observer(() => {
  const store = useStore();
  const selectedIds = store.ui.commandMenu.context.ids;
  const label = `${selectedIds?.length} contacts`;

  return (
    <Command>
      <CommandInput label={label} placeholder='Type a command or search' />
      <Command.List>
        <Command.Group>
          <CommandItem
            leftAccessory={<Tag01 />}
            keywords={['change', 'eidt', 'update', 'tag', 'label', 'profile']}
            onSelect={() => {
              store.ui.commandMenu.setType('EditPersonaTag');
            }}
          >
            Edit persona tag...
          </CommandItem>

          <CommandItem
            leftAccessory={<Certificate02 />}
            onSelect={() => {
              store.ui.commandMenu.setType('EditJobTitle');
            }}
            keywords={[
              'change',
              'edit',
              'update',
              'job',
              'title',
              'position',
              'designation',
            ]}
          >
            Edit job title
          </CommandItem>
          <CommandItem
            leftAccessory={<Certificate02 />}
            onSelect={() => {
              store.ui.commandMenu.setType('ChangeOrAddJobRoles');
            }}
            keywords={[
              'change',
              'edit',
              'update',
              'job',
              'roles',
              'position',
              'function',
            ]}
          >
            Edit job roles...
          </CommandItem>
          <CommandItem
            leftAccessory={<Clock />}
            keywords={['change', 'edit', 'update', 'timezone', 'location']}
            onSelect={() => {
              store.ui.commandMenu.setType('EditTimeZone');
            }}
          >
            Edit time zone...
          </CommandItem>
        </Command.Group>

        <Command.Group heading='Navigate'>
          <GlobalSharedCommands />
        </Command.Group>
      </Command.List>
    </Command>
  );
});
