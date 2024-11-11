import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { Spinner } from '@ui/feedback/Spinner';
import { Star06 } from '@ui/media/icons/Star06';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { TextInput } from '@ui/media/icons/TextInput';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { EmailValidationDetails } from '@graphql/types';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { EmailValidationMessage } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';

interface EmailCellProps {
  contactId: string;
  validationDetails: EmailValidationDetails | undefined;
}

export const EmailCell = observer(
  ({ validationDetails, contactId }: EmailCellProps) => {
    const store = useStore();

    const [isHovered, setIsHovered] = useState(false);

    const contactStore = store.contacts.value.get(contactId);

    const enrichedContact = contactStore?.value.enrichDetails;

    const enrichingStatus =
      !enrichedContact?.enrichedAt &&
      enrichedContact?.requestedAt &&
      !enrichedContact?.failedAt;

    const ref = useRef(null);

    const orgActive =
      contactStore?.value.latestOrganizationWithJobRole?.organization.name;

    const email = contactStore?.value?.primaryEmail?.email;

    const enrichedEmailStatus =
      !enrichedContact?.emailEnrichedAt &&
      enrichedContact?.emailRequestedAt &&
      !email;

    return (
      <div
        ref={ref}
        className='flex cursor-pointer'
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <Menu>
          <MenuButton className='text-ellipsis overflow-hidden whitespace-nowrap'>
            <div className='flex items-center gap-2'>
              {!email && (
                <p className='text-gray-400 '>
                  {enrichingStatus
                    ? 'Enriching...'
                    : enrichedEmailStatus
                    ? 'Finding email...'
                    : 'Not set'}
                </p>
              )}
              {email && (
                <EmailValidationMessage
                  email={email}
                  validationDetails={validationDetails}
                />
              )}
              <p>{email}</p>
            </div>
          </MenuButton>
          <MenuList align='start' className='max-w-[600px] w-[250px]'>
            {orgActive && (
              <MenuItem
                onClick={() => {
                  contactStore?.findEmail();
                }}
              >
                <div className='overflow-hidden text-ellipsis w-[200px]'>
                  <Star06 className='mr-2 text-gray-500' />
                  {`Find email at ${orgActive}`}
                </div>
              </MenuItem>
            )}

            <MenuItem
              onClick={() => {
                if (contactStore?.value.primaryEmail?.email) {
                  store.ui.setSelectionId(
                    contactStore?.value.emails.length || 0 + 1,
                  );

                  contactStore.value.emails.push({
                    id: crypto.randomUUID(),
                    email: '',
                    appSource: '',
                    contacts: [],
                    createdAt: new Date().toISOString(),
                    updatedAt: new Date().toISOString(),
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                  } as any);
                }
                store.ui.commandMenu.setContext({
                  ids: [contactStore?.value.id || ''],
                  entity: 'Contact',
                  property: 'email',
                });
                store.ui.commandMenu.setType('EditEmail');
                store.ui.commandMenu.setOpen(true);
              }}
            >
              <div className='overflow-hidden text-ellipsis'>
                <PlusCircle className='mr-2 text-gray-500' />
                Add new email
              </div>
            </MenuItem>
            {contactStore?.value.emails
              .filter((e) => e.email !== '')
              .map((email) => (
                <MenuItem
                  key={email.email}
                  onClick={() => {
                    contactStore?.setPrimaryEmail(email.id);
                  }}
                >
                  <div className='flex items-center overflow-hidden text-ellipsis justify-between w-full [&_svg]:size-4'>
                    <div className='flex items-center gap-2 max-w-[200px]'>
                      <EmailValidationMessage
                        email={email.email || ''}
                        validationDetails={email.emailValidationDetails}
                      />
                      <p className='truncate'>{email.email}</p>
                    </div>
                    {contactStore.value.primaryEmail?.email ===
                      email?.email && <Check className='text-primary-600' />}
                  </div>
                </MenuItem>
              ))}
          </MenuList>
        </Menu>
        {isHovered &&
          orgActive &&
          (enrichedEmailStatus ? (
            <Tooltip label={`Finding email at ${orgActive} `}>
              <Spinner
                size='sm'
                label='finding email'
                className='text-gray-400 fill-gray-700 ml-2'
              />
            </Tooltip>
          ) : (
            <Tooltip asChild label={`Find email at ${orgActive}`}>
              <IconButton
                size='xxs'
                variant='ghost'
                icon={<Star06 />}
                className={'ml-2'}
                aria-label='Find work email'
                onClick={() => {
                  contactStore?.findEmail();
                }}
              />
            </Tooltip>
          ))}
        {(contactStore?.value.primaryEmail?.email ?? '').length > 0 && (
          <Menu>
            <MenuButton asChild>
              {isHovered && (
                <IconButton
                  size='xxs'
                  variant='ghost'
                  aria-label='edit'
                  className='rounded-[5px] ml-[2px] '
                  icon={<DotsVertical className='text-gray-500' />}
                />
              )}
            </MenuButton>

            <MenuList align='start' side='bottom'>
              <MenuItem
                className='group/edit-email'
                onClick={() => {
                  store.ui.commandMenu.setType('EditEmail');
                  store.ui.commandMenu.setOpen(true);
                }}
              >
                <div className='overflow-hidden text-ellipsis'>
                  <TextInput className='mr-2 group-hover/edit-email:text-gray-700 text-gray-500 ' />
                  Edit email
                </div>
              </MenuItem>
              <MenuItem
                className='group/archive-email'
                onClick={() => {
                  const idx = contactStore?.value.emails.findIndex(
                    (e) => e.email === email,
                  );

                  if (idx !== -1) {
                    contactStore?.value.emails.splice(idx || 0, 1);
                  }
                  contactStore?.commit();
                }}
              >
                <div className='overflow-hidden text-ellipsis'>
                  <Archive className='mr-2 group-hover/archive-email:text-gray-700 text-gray-500' />
                  Archive email
                </div>
              </MenuItem>
            </MenuList>
          </Menu>
        )}
      </div>
    );
  },
);
