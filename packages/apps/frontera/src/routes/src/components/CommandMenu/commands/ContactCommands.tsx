import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Eye } from '@ui/media/icons/Eye.tsx';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { Delete } from '@ui/media/icons/Delete';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { User03 } from '@ui/media/icons/User03.tsx';
import { EyeOff } from '@ui/media/icons/EyeOff.tsx';
import { Shuffle01 } from '@ui/media/icons/Shuffle01.tsx';
import { Building07 } from '@ui/media/icons/Building07.tsx';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp.tsx';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

import { AddPersonaTagSubItemGroup } from './contacts/AddPersonaTagSubItemGroup';
import { AddContactFlowSubItemGroup } from './contacts/AddContactFlowSubItemGroup.tsx';

export const ContactCommands = observer(() => {
  const store = useStore();
  const id = (store.ui.commandMenu.context.ids as string[])?.[0];
  const contact = store.contacts.value.get(id);
  const label = `Contact - ${contact?.name}`;
  const [searchParams] = useSearchParams();

  return (
    <CommandsContainer label={label}>
      <>
        <CommandItem
          leftAccessory={<User03 />}
          rightAccessory={<Kbd>C</Kbd>}
          keywords={contactKeywords.add_contact}
          onSelect={() => {
            if (
              searchParams?.get('preset') ===
              store.tableViewDefs.flowContactsPreset
            ) {
              return store.ui.commandMenu.setType('AddContactsToFlow');
            }

            store.ui.commandMenu.setType('AddContactsBulk');
          }}
        >
          Add contacts
        </CommandItem>

        <CommandItem
          leftAccessory={<Tag01 />}
          keywords={contactKeywords.edit_persona_tag}
          onSelect={() => {
            store.ui.commandMenu.setType('EditPersonaTag');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
              </Kbd>
              <Kbd>T</Kbd>
            </>
          }
        >
          Change or add tags...
        </CommandItem>
        <AddPersonaTagSubItemGroup />

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
              <Kbd>A</Kbd>
            </>
          }
        >
          Add to flow...
        </CommandItem>
        <AddContactFlowSubItemGroup />

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

        {contact && contact?.value.emails.length > 0 && (
          <CommandItem
            leftAccessory={<Mail01 />}
            keywords={contactKeywords.edit_email}
            onSelect={() => {
              store.ui.commandMenu.setType('EditEmail');
            }}
            rightAccessory={
              <>
                <Kbd>
                  <ArrowBlockUp className='text-inherit size-3' />
                </Kbd>
                <Kbd>E</Kbd>
              </>
            }
          >
            Edit email
          </CommandItem>
        )}

        <CommandItem
          leftAccessory={<Mail01 />}
          keywords={contactKeywords.edit_email}
          onSelect={() => {
            store.ui.commandMenu.setType('AddEmail');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
              </Kbd>
              <Kbd>E</Kbd>
            </>
          }
        >
          Add new email
        </CommandItem>

        <CommandItem
          leftAccessory={store.ui.showPreviewCard ? <EyeOff /> : <Eye />}
          rightAccessory={<Kbd className='size-auto h-5 px-1.5'>Space</Kbd>}
          onSelect={() => {
            store.ui.setShowPreviewCard(!store.ui.showPreviewCard);
            store.ui.commandMenu.setOpen(false);
          }}
        >
          {store.ui.showPreviewCard
            ? 'Hide contact preview'
            : 'Preview contact'}
        </CommandItem>

        <CommandItem
          leftAccessory={<Edit03 />}
          keywords={contactKeywords.edit_name}
          onSelect={() => {
            store.ui.commandMenu.setType('EditName');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
              </Kbd>
              <Kbd>R</Kbd>
            </>
          }
        >
          Edit name
        </CommandItem>

        <CommandItem
          leftAccessory={<Building07 />}
          keywords={contactKeywords.change_latest_org}
          onSelect={() => {
            store.ui.commandMenu.setType('EditLatestOrgActive');
          }}
        >
          Change company
        </CommandItem>
        {/* <CommandItem
          leftAccessory={<Phone />}
          keywords={contactKeywords.edit_phone_number}
          onSelect={() => {
            store.ui.commandMenu.setType('EditPhoneNumber');
          }}
        >
          Edit phone number
        </CommandItem> */}
        <CommandItem
          leftAccessory={<Certificate02 />}
          keywords={contactKeywords.edit_job_title}
          onSelect={() => {
            store.ui.commandMenu.setType('EditJobTitle');
          }}
        >
          Edit job title
        </CommandItem>

        {contact?.hasFlows && (
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
                <Delete className='text-inherit size-3' />
              </Kbd>
            </>
          }
        >
          Archive contact
        </CommandItem>
      </>
    </CommandsContainer>
  );
});

const contactKeywords = {
  add_contact: [
    'add',
    'contact',
    'upload',
    'import',
    'new',
    'create',
    'person',
    'individual',
    'email',
    'linkedin',
  ],
  archive_contact: ['archive', 'contact', 'delete', 'remove', 'hide'],
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
  edit_persona_tag: [
    'edit',
    'persona',
    'tag',
    'change',
    'update',
    'label',
    'profile',
  ],
  move_to_flow: ['move', 'to', 'flow'],
  change_latest_org: [
    'edit',
    'organization',
    'company',
    'change',
    'update',
    'latest',
    'active',
  ],
  remove_from_flow: ['remove', 'flow', 'delete'],
};
