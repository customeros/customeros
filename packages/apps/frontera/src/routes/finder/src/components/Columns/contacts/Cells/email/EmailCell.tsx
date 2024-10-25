import { useRef, useMemo, useState } from 'react';

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
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { EmailValidationMessage } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';

interface EmailCellProps {
  contactId: string;
  validationDetails: EmailValidationDetails | undefined;
}

export const EmailCell = observer(
  ({ validationDetails, contactId }: EmailCellProps) => {
    const store = useStore();
    const [isLoading, setIsLoading] = useState(false);

    const [isHovered, setIsHovered] = useState(false);

    const contactStore = store.contacts.value.get(contactId);
    const oldEmail = useMemo(
      () => contactStore?.value?.primaryEmail?.email,
      [],
    );

    const enrichedContact = contactStore?.value.enrichDetails;

    const enrichingStatus =
      !enrichedContact?.enrichedAt &&
      enrichedContact?.requestedAt &&
      !enrichedContact?.failedAt;

    const [isEdit, setIsEdit] = useState(false);
    const ref = useRef(null);

    useOutsideClick({
      ref: ref,
      handler: () => {
        setIsEdit(false);
      },
    });

    const orgActive =
      contactStore?.value.latestOrganizationWithJobRole?.organization.name;

    const email = contactStore?.value?.primaryEmail?.email;

    return (
      <div
        ref={ref}
        className='flex  cursor-pointer'
        onDoubleClick={() => setIsEdit(true)}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        style={{ marginRight: isHovered ? '-20px' : '0px' }}
      >
        <Menu>
          <MenuButton className='text-ellipsis overflow-hidden whitespace-nowrap'>
            <div className='flex items-center gap-2 '>
              {!isEdit && !email && (
                <p className='text-gray-400 '>
                  {enrichingStatus
                    ? 'Enriching...'
                    : isLoading
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
          <MenuList align='center' className='max-w-[600px] w-[250px]'>
            {orgActive && (
              <MenuItem
                onClick={() => {
                  setIsLoading(true);
                  contactStore?.findEmail().finally(() => setIsLoading(false));
                }}
              >
                <div className='overflow-hidden text-ellipsis'>
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
                  contactStore?.update(
                    (c) => {
                      c.emails.push({
                        id: crypto.randomUUID(),
                        email: '',
                        appSource: '',
                        contacts: [],
                        createdAt: new Date().toISOString(),
                        updatedAt: new Date().toISOString(),
                        // eslint-disable-next-line @typescript-eslint/no-explicit-any
                      } as any);

                      return c;
                    },
                    { mutate: false },
                  );
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
            {contactStore?.value.emails.map((email) => (
              <MenuItem
                key={email.email}
                onClick={() => {
                  contactStore?.setPrimaryEmail(email.id);
                }}
              >
                <div className='flex items-center overflow-hidden text-ellipsis justify-between w-full [&_svg]:size-4'>
                  <div className='flex items-center gap-2 max-w-[100px] w-[100px]'>
                    <EmailValidationMessage
                      email={email.email || ''}
                      validationDetails={email.emailValidationDetails}
                    />
                    {email.email}
                  </div>
                  {contactStore.value.primaryEmail?.email === email?.email && (
                    <Check className='text-primary-600' />
                  )}
                </div>
              </MenuItem>
            ))}
          </MenuList>
        </Menu>
        {isHovered &&
          orgActive &&
          (isLoading ? (
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
                  setIsLoading(true);
                  contactStore
                    ?.findEmail()
                    .finally(() =>
                      setTimeout(() => setIsLoading(false), 60000),
                    );
                }}
              />
            </Tooltip>
          ))}
        {(contactStore?.value.primaryEmail?.email ?? '').length > 0 &&
          isHovered && (
            <Menu>
              <MenuButton asChild>
                <IconButton
                  size='xxs'
                  variant='ghost'
                  aria-label='edit'
                  className='rounded-[5px] ml-[2px] '
                  icon={<DotsVertical className='text-gray-500' />}
                />
              </MenuButton>

              <MenuList align='start'>
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

                    contactStore?.update(
                      (c) => {
                        if (idx === 0) {
                          c.emails = [];
                        }

                        if (idx !== undefined && idx > -1) {
                          c.emails.splice(idx, 1);
                        }

                        return c;
                      },
                      { mutate: false },
                    );
                    contactStore?.updateEmail(oldEmail || '', idx);
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
