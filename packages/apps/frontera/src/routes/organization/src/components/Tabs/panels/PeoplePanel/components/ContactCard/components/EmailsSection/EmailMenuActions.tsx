import { observer } from 'mobx-react-lite';
import { SetEmailCase } from '@domain/Contacts/SetEmail.usecase';

import { Copy01 } from '@ui/media/icons/Copy01';
import { Star01 } from '@ui/media/icons/Star01';
import { Send03 } from '@ui/media/icons/Send03';
import { IconButton } from '@ui/form/IconButton';
import { useEvent } from '@shared/hooks/useEvent';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { Stars02 } from '@ui/media/icons/Stars02';
import { TextInput } from '@ui/media/icons/TextInput';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface EmailMenuActionsProps {
  id: string;
  idx: number;
  email: string;
  contactId: string;
}

export const EmailMenuActions = observer(
  ({ contactId, idx, email }: EmailMenuActionsProps) => {
    const store = useStore();
    const [_, copyToClipboard] = useCopyToClipboard();

    const { dispatchEvent } = useEvent<{ email: string; openEditor: string }>(
      'openEmailEditor',
    );

    const contactStore = store.contacts.value.get(contactId);

    const isPrimaryEmail = contactStore?.value.emails[idx]?.primary;
    const useCase = new SetEmailCase();

    if (!contactStore) return;
    useCase.setEntity(contactStore);

    return (
      <Menu>
        <MenuButton asChild>
          <IconButton
            size='xxs'
            variant='ghost'
            icon={<DotsVertical />}
            aria-label='add new email'
          />
        </MenuButton>
        <MenuList>
          <MenuItem
            className='group/edit-email'
            onClick={() => {
              store.ui.setSelectionId(idx);
              store.ui.commandMenu.setType('EditEmail');
              store.ui.commandMenu.setContext({
                ids: [contactStore?.value.metadata.id ?? ''],
                entity: 'Contact',
                property: 'email',
              });
              store.ui.commandMenu.setOpen(true);
            }}
          >
            <div className='flex items-center gap-2'>
              <TextInput className='group-hover/edit-email:text-gray-700 text-gray-500' />
              <span>Edit email</span>
            </div>
          </MenuItem>
          {!isPrimaryEmail && (
            <MenuItem
              className='group/edit-email'
              onClick={() => {
                useCase.setEmail(email || '');
                useCase.setPrimaryEmailForContact();
              }}
            >
              <div className='flex items-center gap-2'>
                <Star01 className='group-hover/edit-email:text-gray-700 text-gray-500' />
                <span>Make primary</span>
              </div>
            </MenuItem>
          )}
          <MenuItem
            className='group/send-email'
            onClick={() => {
              dispatchEvent({ email: email, openEditor: 'email' });
            }}
          >
            <Send03 className='text-gray-500 group-hover/send-email:text-gray-700' />
            Send email to contact
          </MenuItem>
          <MenuItem
            className='group/add-email'
            onClick={() => {
              store.ui.setSelectionId(contactStore?.value.emails.length || 1);

              contactStore?.value.emails.push({
                id: crypto.randomUUID(),
                email: '',
                appSource: '',
                contacts: [],
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString(),
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
              } as any);

              store.ui.commandMenu.setContext({
                ids: [contactStore?.id || ''],
                entity: 'Contact',
                property: 'email',
              });
              store.ui.commandMenu.setType('EditEmail');
              store.ui.commandMenu.setOpen(true);
            }}
          >
            <PlusCircle className='text-gray-500 group-hover/add-email:text-gray-700' />
            Add email
          </MenuItem>
          <MenuItem
            className='group/work-email'
            onClick={() => {
              contactStore.findEmail();
            }}
          >
            <Stars02 className='text-gray-500 group-hover/work-email:text-gray-700' />
            Find work email
          </MenuItem>
          <MenuItem
            className='group/archive-email'
            onClick={() => {
              contactStore.draft();
              contactStore?.value.emails.splice(idx, 1);

              contactStore?.commit();
            }}
          >
            <Archive className='text-gray-500 group-hover/archive-email:text-gray-700' />
            Archive email
          </MenuItem>

          <MenuItem
            className='group/copy-email'
            onClick={() => copyToClipboard(email || '', 'Email copied')}
          >
            <Copy01 className='group-hover/copy-email:text-gray-700 text-gray-500' />
            Copy email
          </MenuItem>
        </MenuList>
      </Menu>
    );
  },
);
