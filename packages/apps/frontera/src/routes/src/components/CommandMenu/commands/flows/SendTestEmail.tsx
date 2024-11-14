import React, { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const SendTestEmail = observer(() => {
  const store = useStore();
  const [emailAddress, setEmailAddress] = useState('');
  const [error, setError] = useState('');
  // const { flows } = useStore();
  const flowId = store.ui.commandMenu.context.ids[0];
  const inputRef = useRef<HTMLInputElement>(null);
  const flow = store.flows.value.get(flowId);

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  const handleConfirm = () => {
    const isValid = isValidEmail(emailAddress);

    if (!isValid) {
      setError('Please enter a valid email address');
      inputRef.current?.focus();

      return;
    }

    flow?.sendTestEmail(
      {
        sendToEmailAddress: emailAddress,
        subject: store.ui.commandMenu.context.meta?.subject || '',
        bodyTemplate: store.ui.commandMenu.context.meta?.bodyTemplate || '',
      },
      {
        onSuccess: () => {
          store.ui.commandMenu.setOpen(false);
          store.ui.commandMenu.clearContext();
        },
      },
    );
  };

  return (
    <Command label={`Send test email`}>
      <div className='p-6 pb-0 flex flex-col gap-1'>
        <div className='flex items-center justify-between'>
          <h1
            className='text-base font-semibold'
            data-test='create-new-flow-modal-title'
          >
            Send a test email...
          </h1>

          <CommandCancelIconButton
            onClose={() => {
              store.ui.commandMenu.setOpen(false);
            }}
          />
        </div>
      </div>

      <div className='pr-6 pl-6 pb-6 flex flex-col gap-2 '>
        <p className='text-sm mt-2'>
          Immediately from {/* todo replace with the actual email*/} mailbox
          setup for your tenant using random variables if applicable
        </p>
        <Input
          autoFocus
          size='sm'
          ref={inputRef}
          id='emailAddress'
          variant='unstyled'
          value={emailAddress}
          placeholder='To email address'
          dataTest='send-test-flow-email'
          onChange={(e) => {
            setEmailAddress(e.target.value);
          }}
          className={cn({
            'border-b border-error-400': error.length > 0,
          })}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              handleConfirm();
            }
          }}
        />
        {error.length > 0 && <p className='text-xs text-error-600'>{error}</p>}
      </div>

      <div className='flex w-full gap-2 pl-6 pr-6 pb-6'>
        <CommandCancelButton
          dataTest='cancel-create-new-flow'
          onClose={() => {
            store.ui.commandMenu.setOpen(false);
          }}
        />

        <Button
          className='w-full'
          colorScheme='primary'
          onClick={handleConfirm}
          loadingText={'Sending…'}
          isLoading={flow?.isLoading}
          data-test='confirm-create-new-flow'
        >
          Send test email
        </Button>
      </div>
    </Command>
  );
});

const isValidEmail = (email: string) => {
  if (!email || typeof email !== 'string') return false;

  // Simple but effective email regex
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  return emailRegex.test(email);
};
