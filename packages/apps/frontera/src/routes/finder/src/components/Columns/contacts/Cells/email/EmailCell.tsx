import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';
import { EditEmailCase } from '@domain/usecases/command-menu/edit-email.usecase';

import { Check } from '@ui/media/icons/Check';
import { Spinner } from '@ui/feedback/Spinner';
import { Star06 } from '@ui/media/icons/Star06';
import { Copy01 } from '@ui/media/icons/Copy01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { TextInput } from '@ui/media/icons/TextInput';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { EmailValidationMessage } from '@shared/components/ContactDetails/components/EmailsSection/EmailValidationMessage';

interface EmailCellProps {
  contactId: string;
}

const editEmailUseCase = EditEmailCase.getInstance();

export const EmailCell = observer(({ contactId }: EmailCellProps) => {
  const store = useStore();

  const [isHovered, setIsHovered] = useState(false);
  const [isOpened, setIsOpened] = useState(false);
  const [_, copyToClipboard] = useCopyToClipboard();

  const contactStore = store.contacts.value.get(contactId);

  if (!contactStore) return;

  const isEnrichingContact = contactStore?.isEnriching;

  const ref = useRef(null);

  const activeOrgId = contactStore?.value.primaryOrganizationId;
  const domains =
    activeOrgId && store.organizations.value.get(activeOrgId)?.value?.domains;
  const orgActive = contactStore?.value.primaryOrganizationName;

  const email = contactStore.value.emails.find((e) => e.primary)?.email;

  const isEnrichingEmail = contactStore?.emailEnriching;
  const enrichedEmailNotFound =
    !contactStore.value.enrichedEmailFound &&
    !email &&
    contactStore.value.enrichedEmailEnrichedAt;

  const enrichingStatus = isEnrichingContact
    ? 'Enriching...'
    : isEnrichingEmail
    ? 'Finding email...'
    : enrichedEmailNotFound
    ? 'Not found'
    : 'Not set';

  const validationDetails = contactStore?.value.emails.find(
    (e) => e.primary === true,
  )?.emailValidationDetails;

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
            <MenuItem>
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
              store.ui.commandMenu.setContext({
                ids: [contactStore?.id || ''],
                entity: 'Contact',
                property: 'email',
              });
              store.ui.commandMenu.setType('AddEmail');
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
                  contactStore.value.emails.forEach((e) => {
                    e.primary = false;
                  });
                  contactStore?.draft();

                  const findIdx = contactStore?.value.emails.findIndex(
                    (e) => e.email === email.email,
                  );

                  contactStore.value.emails[findIdx].primary = true;
                  contactStore?.commit();
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
                  {email?.primary === true && (
                    <Check className='text-primary-600' />
                  )}
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
              className='ml-2'
              icon={<Star06 />}
              aria-label='Find work email'
              onClick={() => {
                contactStore.findEmail();
              }}
            />
          </Tooltip>
        )}
      {(email ?? '').length > 0 && (
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
                editEmailUseCase.setEmail(email!);
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
                  contactStore?.draft();
                  contactStore?.value.emails.splice(idx, 1);
                  contactStore?.commit();
                }
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
});
