import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { Spinner } from '@ui/feedback/Spinner';
import { Star06 } from '@ui/media/icons/Star06';
import { Copy01 } from '@ui/media/icons/Copy01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { TextInput } from '@ui/media/icons/TextInput';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { EmailValidationDetails } from '@graphql/types';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
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
    const [isOpened, setIsOpened] = useState(false);
    const [_, copyToClipboard] = useCopyToClipboard();

    const contactStore = store.contacts.value.get(contactId);

    const enrichedContact = contactStore?.value.enrichDetails;

    const isEnrichingContact = contactStore?.isEnriching;

    const ref = useRef(null);

    const activeOrgId =
      contactStore?.value.latestOrganizationWithJobRole?.organization?.metadata
        ?.id;

    const domains =
      activeOrgId && store.organizations.value.get(activeOrgId)?.value?.domains;
    const orgActive =
      contactStore?.value.latestOrganizationWithJobRole?.organization?.name;

    const email = contactStore?.value?.primaryEmail?.email;

    const isEnrichingEmail =
      !enrichedContact?.emailEnrichedAt &&
      enrichedContact?.emailRequestedAt &&
      !email;

    const enrichedEmailNotFound =
      !enrichedContact?.emailFound &&
      !email &&
      enrichedContact?.emailEnrichedAt;

    const enrichingStatus = isEnrichingContact
      ? 'Enriching...'
      : isEnrichingEmail
      ? 'Finding email...'
      : enrichedEmailNotFound
      ? 'Not found'
      : 'Not set';

    return (
      <div
        ref={ref}
        className='flex cursor-pointer'
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <Menu onOpenChange={(newStatus) => setIsOpened(newStatus)}>
          <MenuButton
            asChild
            className='text-ellipsis overflow-hidden whitespace-nowrap'
          >
            <div className='flex items-center gap-2'>
              {!email && <p className='text-gray-400 '>{enrichingStatus}</p>}
              {email && (
                <EmailValidationMessage
                  email={email}
                  validationDetails={validationDetails}
                />
              )}
              <div className='inline-flex gap-2'>
                {email}
                {isEnrichingEmail && (
                  <Tooltip label={`Finding email at ${orgActive}`}>
                    <div>
                      <Spinner
                        size='sm'
                        label='finding email'
                        className='text-gray-400 fill-gray-700 ml-2'
                      />
                    </div>
                  </Tooltip>
                )}
              </div>
            </div>
          </MenuButton>
          <MenuList align='start' className='max-w-[600px] w-[250px]'>
            {orgActive && !!domains?.length && (
              <MenuItem onClick={() => contactStore?.findEmail()}>
                <div className=' flex overflow-hidden items-center text-ellipsis w-[200px]'>
                  {isEnrichingEmail ? (
                    <Tooltip label={`Finding email at ${orgActive}`}>
                      <Spinner
                        size='sm'
                        label='finding email'
                        className='text-gray-400 fill-gray-700 mr-2'
                      />
                    </Tooltip>
                  ) : (
                    <Star06 className='mr-2 text-gray-500 size-4' />
                  )}
                  <p className='w-[190px] truncate'>
                    {isEnrichingEmail
                      ? `Finding email at ${orgActive}`
                      : `Find email at ${orgActive}`}
                  </p>
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
        {(isHovered || isOpened) &&
          orgActive &&
          !isEnrichingEmail &&
          !!domains?.length && (
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
          )}
        {(contactStore?.value.primaryEmail?.email ?? '').length > 0 && (
          <Menu onOpenChange={(newStatus) => setIsOpened(newStatus)}>
            <MenuButton asChild>
              {(isHovered || isOpened) && (
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
              {email && (
                <MenuItem
                  className='group/copy-email'
                  onClick={() => {
                    copyToClipboard(email, 'Email copied');
                  }}
                >
                  <div className='overflow-hidden text-ellipsis'>
                    <Copy01 className='group-hover/copy-email:text-gray-700 text-gray-500 mr-2' />
                    Copy email
                  </div>
                </MenuItem>
              )}
            </MenuList>
          </Menu>
        )}
      </div>
    );
  },
);
