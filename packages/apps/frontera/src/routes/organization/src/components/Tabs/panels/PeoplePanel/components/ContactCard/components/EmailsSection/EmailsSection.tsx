import { Fragment } from 'react';

import { uniqBy } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Spinner } from '@ui/feedback/Spinner';
import { Star06 } from '@ui/media/icons/Star06';
import { Mail01 } from '@ui/media/icons/Mail01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';

import { EmailMenuActions } from './EmailMenuActions';
import { EmailValidationMessage } from '../../EmailValidationMessage';

interface EmailsSectionProps {
  contactId: string;
}

export const EmailsSection = observer(({ contactId }: EmailsSectionProps) => {
  const store = useStore();
  const [_, copyToClipboard] = useCopyToClipboard();

  const contactStore = store.contacts.value.get(String(contactId));

  const activeCompany =
    (contactStore?.value?.organizations?.content?.length ?? 1) - 1;
  const company = contactStore?.value.organizations?.content?.[activeCompany];

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
  ).sort((a, b) => (b.primary ? 1 : -1));
  const enrichedContact = contactStore?.value.enrichDetails;

  const isEnrichingEmail =
    !enrichedContact?.emailEnrichedAt &&
    enrichedContact?.emailRequestedAt &&
    !isPrimaryEmail;

  return (
    <div>
      <div className='flex items-start w-full'>
        <div className='flex mr-4'>
          <Mail01 className='mt-0.5 text-gray-500' />
        </div>

        <div className='flex flex-col flex-9 w-full'>
          {allEmails.length === 0 && (
            <div className='flex w-full gap-2'>
              <span className='text-gray-400'>
                {isEnrichingEmail
                  ? 'Finding email...'
                  : !enrichedContact?.emailFound
                  ? 'Work email not found'
                  : 'Work email'}
              </span>
              {isEnrichingEmail ? (
                <Tooltip label={`Finding email at ${company?.name}`}>
                  <Spinner
                    size='sm'
                    label='finding email'
                    className='text-gray-400 fill-gray-700 mr-2'
                  />
                </Tooltip>
              ) : (
                <IconButton
                  size='xxs'
                  variant='ghost'
                  icon={<Star06 />}
                  onClick={() => {}}
                  colorScheme='grayModern'
                  className='size-5 mt-0.5'
                  aria-label='enrich-work-email'
                />
              )}
            </div>
          )}
          {allEmails?.map((email, idx) => (
            <Fragment key={`${idx}-${email.id}`}>
              <div className=' flex items-center justify-between w-full'>
                <div key={email.id} className='flex items-center gap-1'>
                  <span
                    className='text-sm max-w-[170px] text-ellipsis overflow-hidden'
                    onClick={() =>
                      copyToClipboard(email?.email || '', 'Email copied')
                    }
                  >
                    {email.email || 'Not set'}
                  </span>
                  {contactStore?.value.emails.length !== 1 && email.primary && (
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
                  <EmailMenuActions
                    idx={idx}
                    id={email.id}
                    email={email.email || ''}
                    contactId={contactId || ''}
                  />
                </div>
              </div>
            </Fragment>
          ))}
        </div>
      </div>
    </div>
  );
});

{
  /* <>
<div className='flex items-center justify-between w-full text-sm group/menu'>
  <div className='flex items-center gap-2'>
    {allEmails.length === 0 && (
      <span className='text-gray-400 ml-1'>
        {isEnrichingEmail ? 'Finding email' : 'No emails yet'}
      </span>
    )}

    <div className='flex items-center gap-1'>
      {isEnrichingEmail ? (
        <Tooltip label={`Finding email at ${company?.name}`}>
          <Spinner
            size='sm'
            label='finding email'
            className='text-gray-400 fill-gray-700 mr-2'
          />
        </Tooltip>
      ) : company?.name &&
        domains?.length &&
        !contactStore?.value.primaryEmail?.email ? (
        <IconButton
          size='xxs'
          variant='ghost'
          icon={<Star06 />}
          aria-label='find-work-email'
          onClick={() => contactStore?.findEmail()}
        />
      ) : null}
    </div>
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
    {isEnrichingEmail && (
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
</> */
}
