import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Phone } from '@ui/media/icons/Phone';
import { Clock } from '@ui/media/icons/Clock';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const ContactCommands = observer(() => {
  const store = useStore();
  const id = (store.ui.commandMenu.context.ids as string[])?.[0];
  const contact = store.contacts.value.get(id);
  const label = `Contact - ${contact?.value.name}`;

  return (
    <CommandsContainer label={label}>
      <>
        <CommandItem
          leftAccessory={<Tag01 />}
          onSelect={() => {
            store.ui.commandMenu.setType('EditPersonaTag');
          }}
          keywords={[
            'edit',
            'persona',
            'tag',
            'change',
            'update',
            'label',
            'profile',
          ]}
        >
          Edit persona tag...
        </CommandItem>

        {!!contact?.value?.tags?.length && (
          <CommandItem
            leftAccessory={<Tag01 />}
            onSelect={() => {
              contact?.removeAllTagsFromContact();
              store.ui.commandMenu.setOpen(false);
            }}
          >
            Remove tags
          </CommandItem>
        )}

        <CommandItem
          leftAccessory={<Mail01 />}
          keywords={['edit', 'email', 'change', 'update', 'address', '@']}
          onSelect={() => {
            store.ui.commandMenu.setType('EditEmail');
          }}
        >
          Edit email
        </CommandItem>

        <CommandItem
          leftAccessory={<Edit03 />}
          keywords={['edit', 'name', 'change', 'update', 'rename', 'contact']}
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
            'edit',
            'phone',
            'number',
            'change',
            'update',
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
            'edit',
            'job',
            'title',
            'change',
            'update',
            'position',
            'designation',
          ]}
        >
          Edit job title
        </CommandItem>
        <CommandItem
          leftAccessory={<Certificate02 />}
          keywords={['edit', 'job', 'roles', 'update', 'position', 'function']}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeOrAddJobRoles');
          }}
        >
          Edit job roles...
        </CommandItem>
        <CommandItem
          leftAccessory={<Clock />}
          keywords={['edit', 'timezone', 'change', 'update', 'location']}
          onSelect={() => {
            store.ui.commandMenu.setType('EditTimeZone');
          }}
        >
          Edit time zone...
        </CommandItem>
      </>
    </CommandsContainer>
  );
});
