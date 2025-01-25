import { Fragment } from 'react';

import { observer } from 'mobx-react-lite';

import { Spinner } from '@ui/feedback/Spinner';
import { Star06 } from '@ui/media/icons/Star06';
import { Mail01 } from '@ui/media/icons/Mail01';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';

import { EmailMenuActions } from './EmailMenuActions';
import { EmailValidationMessage } from './EmailValidationMessage';

interface EmailsSectionProps {
  contactId: string;
}

export const EmailsSection = observer(({ contactId }: EmailsSectionProps) => {
  const store = useStore();
  const [_, copyToClipboard] = useCopyToClipboard();

  const contactStore = store.contacts.value.get(contactId);

  const orgName = contactStore?.value.primaryOrganizationName;

  const allEmails = contactStore?.value.emails || [];
  const enrichedEmail = contactStore?.emailEnriching;
  const isEnrichingEmail = enrichedEmail;

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
                data-test='add-work-email'
                className='text-gray-400 cursor-pointer text-sm'
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
                {isEnrichingEmail
                  ? 'Finding email...'
                  : contactStore?.value.enrichedEmailEnrichedAt
                  ? 'Work email not found'
                  : 'Work email'}
              </p>
              {isEnrichingEmail ? (
                <Tooltip label={`Finding email at ${orgName}`}>
                  <Spinner
                    size='sm'
                    label='finding email'
                    className='text-gray-400 fill-gray-700 mr-2'
                  />
                </Tooltip>
              ) : (
                <Tooltip label={`Finding email at ${orgName}`}>
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    icon={<Star06 />}
                    className='mt-0.5'
                    colorScheme='grayModern'
                    aria-label='enrich-work-email'
                    onClick={() => contactStore?.findEmail()}
                  />
                </Tooltip>
              )}
            </div>
          )}
          {allEmails?.map((email, idx) => (
            <Fragment key={`${idx}-${email.id}`}>
              <div className=' flex items-center justify-between w-full'>
                <div key={email.id} className='flex items-center'>
                  <p
                    className='text-sm max-w-[230px] text-ellipsis overflow-hidden cursor-default'
                    onClick={() =>
                      copyToClipboard(email?.email || '', 'Email copied')
                    }
                  >
                    {email?.email!.length > 0 ? email.email : 'Not set'}
                  </p>

                  {contactStore?.value.emails.length !== 1 && email.primary && (
                    <span className='text-gray-500 text-sm ml-1'>
                      {' '}
                      • Primary
                    </span>
                  )}
                </div>
                <div className='flex items-center gap-2'>
                  {isEnrichingEmail && idx === 0 && (
                    <Spinner
                      size='sm'
                      label='finding email'
                      className='text-gray-400 fill-gray-700 mr-2'
                    />
                  )}
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
