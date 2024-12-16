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

  const isPrimaryEmail = contactStore?.value?.primaryEmail;

  const allEmails = uniqBy(
    contactStore
      ? [
          ...contactStore.value.emails,
          ...(contactStore.value?.primaryEmail
            ? [contactStore.value?.primaryEmail]
            : []),
        ]
      : [],
    'id',
  ).sort((_a, b) => (b.primary ? 1 : -1));
  const enrichedContact = contactStore?.value.enrichDetails;

  const isEnrichingEmail =
    !enrichedContact?.emailEnrichedAt &&
    enrichedContact?.emailRequestedAt &&
    !isPrimaryEmail;

  return (
    <div>
      <div className='flex justify-center items-center w-full'>
        <div className='flex mr-4 items-center justify-between'>
          <Mail01 className='text-gray-500 mt-0.5' />
        </div>

        <div className='flex flex-col flex-9 w-full'>
          {allEmails.length === 0 && (
            <div className='flex w-full gap-2 items-center'>
              <p
                className='text-gray-400 cursor-pointer text-sm'
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
                    ids: [contactStore?.id || ''],
                    entity: 'Contact',
                    property: 'email',
                  });
                  store.ui.commandMenu.setType('EditEmail');
                  store.ui.commandMenu.setOpen(true);
                }}
              >
                {isEnrichingEmail
                  ? 'Finding email...'
                  : enrichedContact?.emailFound
                  ? 'Work email not found'
                  : 'Work email'}
              </p>
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
                  className='mt-0.5'
                  colorScheme='grayModern'
                  aria-label='enrich-work-email'
                />
              )}
            </div>
          )}
          {allEmails?.map((email, idx) => (
            <Fragment key={`${idx}-${email.id}`}>
              <div className=' flex items-center justify-between w-full'>
                <div key={email.id} className='flex items-center'>
                  <p
                    className='text-sm max-w-[230px] text-ellipsis overflow-hidden'
                    onClick={() =>
                      copyToClipboard(email?.email || '', 'Email copied')
                    }
                  >
                    {email.email || 'Not set'}
                  </p>
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
