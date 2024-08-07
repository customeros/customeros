import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Phone } from '@ui/media/icons/Phone';
import { Clock } from '@ui/media/icons/Clock';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import { GlobalSharedCommands } from '@shared/components/CommandMenu/commands/GlobalHub.tsx';

export const ContactCommands = observer(() => {
  const store = useStore();
  const id = (store.ui.commandMenu.context.ids as string[])?.[0];
  const contact = store.contacts.value.get(id);
  const label = `Contact - ${contact?.value.name}`;

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
            Edit Persona tag...
          </CommandItem>

          <CommandItem
            leftAccessory={<Mail01 />}
            keywords={['change', 'edit', 'update', 'email', 'address', '@']}
            onSelect={() => {
              store.ui.commandMenu.setType('EditEmail');
            }}
          >
            Edit email
          </CommandItem>

          <CommandItem
            leftAccessory={<Edit03 />}
            keywords={['change', 'edit', 'update', 'name', 'rename', 'contact']}
            onSelect={() => {
              store.ui.commandMenu.setType('EditName');
            }}
          >
            Edit name
          </CommandItem>
          <CommandItem
            leftAccessory={<Phone />}
            onSelect={() => {
              store.ui.commandMenu.setType('EditPhoneNumber');
            }}
            keywords={[
              'change',
              'edit',
              'update',
              'phone',
              'number',
              'mobile',
              'telephone',
            ]}
          >
            Edit phone number
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
            Change or add job roles...
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
