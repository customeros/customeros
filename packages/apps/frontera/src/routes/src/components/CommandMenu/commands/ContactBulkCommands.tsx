import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Clock } from '@ui/media/icons/Clock';
import { Delete } from '@ui/media/icons/Delete';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { Shuffle01 } from '@ui/media/icons/Shuffle01.tsx';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp.tsx';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';

import { CommandsContainer } from './shared';

export const ContactBulkCommands = observer(() => {
  const store = useStore();
  const selectedIds = store.ui.commandMenu.context.ids;
  const label = `${selectedIds?.length} contacts`;
  const contactStores = selectedIds.map((e) => store.contacts.value.get(e));

  return (
    <CommandsContainer label={label}>
      <>
        <CommandItem
          leftAccessory={<Tag01 />}
          keywords={contactKeywords.edit_persona_tag}
          onSelect={() => {
            store.ui.commandMenu.setType('EditPersonaTag');
          }}
        >
          Edit persona tag...
        </CommandItem>

        <CommandItem
          leftAccessory={<Shuffle01 />}
          keywords={contactKeywords.move_to_flow}
          onSelect={() => {
            store.ui.commandMenu.setType('EditContactFlow');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
              </Kbd>
              <Kbd>Q</Kbd>
            </>
          }
        >
          Add to flow...
        </CommandItem>
        <CommandItem
          leftAccessory={<Certificate02 />}
          keywords={contactKeywords.edit_job_title}
          onSelect={() => {
            store.ui.commandMenu.setType('EditJobTitle');
          }}
        >
          Edit job title
        </CommandItem>
        <CommandItem
          leftAccessory={<Certificate02 />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeOrAddJobRoles');
          }}
        >
          Edit job roles...
        </CommandItem>
        <CommandItem
          leftAccessory={<Clock />}
          keywords={contactKeywords.edit_time_zone}
          onSelect={() => {
            store.ui.commandMenu.setType('EditTimeZone');
          }}
        >
          Edit time zone...
        </CommandItem>

        {contactStores.some((contact) => !!contact?.flows?.length) && (
          <CommandItem
            leftAccessory={<Shuffle01 />}
            keywords={contactKeywords.remove_from_flow}
            onSelect={() => {
              store.ui.commandMenu.setType('UnlinkContactFromFlow');
            }}
          >
            Remove from flow
          </CommandItem>
        )}

        <CommandItem
          leftAccessory={<Archive />}
          keywords={contactKeywords.archive_contact}
          onSelect={() => {
            store.ui.commandMenu.setType('DeleteConfirmationModal');
          }}
          rightAccessory={
            <>
              <CommandKbd />
              <Kbd>
                <Delete className='size-3' />
              </Kbd>
            </>
          }
        >
          Archive contacts
        </CommandItem>
      </>
    </CommandsContainer>
  );
});

const contactKeywords = {
  archive_contact: ['archive', 'contact', 'delete', 'remove', 'hide'],
  edit_persona_tag: [
    'edit',
    'persona',
    'tag',
    'change',
    'update',
    'label',
    'profile',
  ],
  edit_job_title: [
    'edit',
    'job',
    'title',
    'change',
    'update',
    'position',
    'designation',
  ],
  edit_job_roles: ['edit', 'job', 'roles', 'update', 'position', 'function'],
  edit_time_zone: ['edit', 'timezone', 'change', 'update', 'location'],
  edit_email: ['edit', 'email', 'change', 'update', 'address', '@'],
  edit_name: ['edit', 'name', 'change', 'update', 'rename', 'contact'],
  edit_phone_number: [
    'edit',
    'phone',
    'number',
    'change',
    'update',
    'mobile',
    'telephone',
  ],
  move_to_flow: ['move', 'to', 'flow', 'edit', 'change'],
  remove_from_flow: ['remove', 'flow', 'delete'],
};
