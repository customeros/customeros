import { Fragment } from 'react';

import { uniqBy } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Plus } from '@ui/media/icons/Plus';
import { Spinner } from '@ui/feedback/Spinner';
import { Mail02 } from '@ui/media/icons/Mail02';
import { Star06 } from '@ui/media/icons/Star06';
import { Star01 } from '@ui/media/icons/Star01';
import { Copy01 } from '@ui/media/icons/Copy01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { TextInput } from '@ui/media/icons/TextInput';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { EmailValidationMessage } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';

interface EmailsSectionProps {
  contactId: string | number;
}

export const EmailsSection = observer(({ contactId }: EmailsSectionProps) => {
  const store = useStore();
  const [_, copyToClipboard] = useCopyToClipboard();

  const contactStore = store.contacts.value.get(String(contactId));

  const activeCompany =
    (contactStore?.value?.organizations?.content?.length ?? 1) - 1;
  const company =
    contactStore?.value.organizations?.content?.[activeCompany]?.name;

  const isPrimaryEmail = contactStore?.value.primaryEmail;

  const allEmails = uniqBy(
    contactStore
      ? [
          ...contactStore.value.emails,
          ...(contactStore.value.primaryEmail
            ? [contactStore.value.primaryEmail]
            : []),
        ]
      : [],
    'id',
  );
  const enrichedContact = contactStore?.value.enrichDetails;
  const enrichedEmailStatus =
    !enrichedContact?.emailEnrichedAt &&
    enrichedContact?.emailRequestedAt &&
    !isPrimaryEmail;

  return (
    <>
      <div className='flex items-center justify-between w-full text-sm group/menu'>
        <div className='flex items-center gap-2'>
          <Mail02 className='mt-[1px] text-gray-500' />
          <span className='text-gray-500'>Emails</span>
          {allEmails.length === 0 && (
            <span className='text-gray-400 ml-[57px]'>No emails yet</span>
          )}
        </div>

        <div className='flex items-center gap-2'>
          <Menu>
            <MenuButton asChild>
              <div>
                <Tooltip align='end' side='bottom' label={'Add new email'}>
                  <div>
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      icon={<Plus />}
                      aria-label='add new email'
                    />
                  </div>
                </Tooltip>
              </div>
            </MenuButton>
            <MenuList>
              {company && (
                <MenuItem
                  className='group/find-email '
                  onClick={() => {
                    contactStore?.findEmail();
                  }}
                >
                  <div className='flex items-center gap-1'>
                    <Star06 className='group-hover/find-email:text-gray-700 text-gray-500' />
                    <span className='max-w-[150px] text-ellipsis overflow-hidden whitespace-nowrap'>{`Find email at ${company}`}</span>
                  </div>
                </MenuItem>
              )}
              <MenuItem
                className='group/add-email'
                onClick={() => {
                  store.ui.setSelectionId(
                    contactStore?.value.emails.length || 1,
                  );

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
                    ids: [contactStore?.value.id || ''],
                    entity: 'Contact',
                    property: 'email',
                  });
                  store.ui.commandMenu.setType('EditEmail');
                  store.ui.commandMenu.setOpen(true);
                }}
              >
                <PlusCircle className='text-gray-500 group-hover/add-email:text-gray-700' />
                Add new email
              </MenuItem>
            </MenuList>
          </Menu>
          {enrichedEmailStatus && (
            <Tooltip label={`Finding email at ${company} `}>
              <div>
                <Spinner
                  size='sm'
                  label='finding email'
                  className='text-gray-400 fill-gray-700'
                />
              </div>
            </Tooltip>
          )}
        </div>
      </div>
      <div className='ml-6'>
        {allEmails?.map((email, idx) => (
          <Fragment key={`${idx}-${email.id}`}>
            <div className=' flex items-center justify-between '>
              <div key={email.id} className='flex items-center gap-1 '>
                <span
                  className='text-sm max-w-[170px] text-ellipsis overflow-hidden'
                  onClick={() =>
                    copyToClipboard(email?.email || '', 'Email copied')
                  }
                >
                  {email.email || 'Not set'}
                </span>
                {isPrimaryEmail?.id === email.id && (
                  <span className='text-gray-500 text-sm'> • Primary</span>
                )}
              </div>
              <div className='flex items-center gap-2'>
                {email && (
                  <EmailValidationMessage
                    email={email?.email || ''}
                    validationDetails={email.emailValidationDetails}
                  />
                )}
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
                    {isPrimaryEmail?.id !== email.id && (
                      <MenuItem
                        className='group/edit-email'
                        onClick={() => {
                          contactStore?.setPrimaryEmail(email.id);
                        }}
                      >
                        <div className='flex items-center gap-2'>
                          <Star01 className='group-hover/edit-email:text-gray-700 text-gray-500' />
                          <span>Make primary</span>
                        </div>
                      </MenuItem>
                    )}

                    <MenuItem
                      className='group/edit-email'
                      onClick={() => {
                        store.ui.setSelectionId(idx);
                        store.ui.commandMenu.setType('EditEmail');
                        store.ui.commandMenu.setContext({
                          ids: [contactStore?.value.id ?? ''],
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
                    <MenuItem
                      className='group/archive-email'
                      onClick={() => {
                        contactStore?.value.emails.splice(idx, 1);

                        contactStore?.commit();
                      }}
                    >
                      <Archive className='text-gray-500 group-hover/archive-email:text-gray-700' />
                      Archive email
                    </MenuItem>
                    <MenuItem
                      className='group/copy-email'
                      onClick={() =>
                        copyToClipboard(email.email || '', 'Email copied')
                      }
                    >
                      <Copy01 className='group-hover/copy-email:text-gray-700 text-gray-500' />
                      Copy email
                    </MenuItem>
                  </MenuList>
                </Menu>
              </div>
            </div>
          </Fragment>
        ))}
      </div>
    </>
  );
});
