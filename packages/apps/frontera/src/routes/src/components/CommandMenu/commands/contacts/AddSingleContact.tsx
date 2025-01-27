import { useRef, useEffect, MouseEvent, KeyboardEvent } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { CreateContact } from '@domain/usecases/contact-details/create-contact.usecase.ts';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { useModKey } from '@shared/hooks/useModKey';
import { Signature } from '@ui/media/icons/Signature.tsx';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';
import { Command, CommandCancelIconButton } from '@ui/overlay/CommandMenu';

const contactCreate = new CreateContact();
export const AddSingleContact = observer(() => {
  const store = useStore();
  const inputRef = useRef<HTMLInputElement>(null);
  const inputPlaceholder =
    contactCreate.getType === 'linkedin'
      ? 'linkedin.com/in/johnlemon'
      : contactCreate.getType === 'email'
      ? 'john@heyjude.band'
      : 'First and last name';
  const confirmButtonPlaceholder =
    contactCreate.getType === 'name' ? 'Add contact' : 'Add & enrich';

  const handleClose = (
    e?:
      | MouseEvent<HTMLButtonElement>
      | KeyboardEvent<HTMLInputElement>
      | KeyboardEvent<HTMLButtonElement>,
  ) => {
    e?.stopPropagation();
    e?.preventDefault();
    contactCreate.clearState();
    store.ui.commandMenu.toggle('AddSingleContact');
    store.ui.commandMenu.clearContext();
  };

  useEffect(() => {
    if (store.ui.commandMenu.context?.meta?.email.length) {
      contactCreate.setType('email');
      contactCreate.setInputValue(store.ui.commandMenu.context.meta.email);
    }
  }, [store.ui.commandMenu.context?.meta]);

  const handleSubmit = async () => {
    contactCreate.setOrganizationId(
      store.ui.commandMenu.context?.ids?.[0] as string,
    );

    if (contactCreate.getType === 'name') {
      await contactCreate.submit();

      !contactCreate.invalidName && handleClose();
    }

    if (contactCreate.getType === 'linkedin') {
      await contactCreate.submit();

      if (contactCreate.emptyLinkedInUrl || contactCreate.invalidLinkedInUrl)
        return;
      if (contactCreate.errorLinkedIn) return;
      handleClose();
    }

    if (contactCreate.getType === 'email') {
      await contactCreate.submit();

      if (contactCreate.emptyEmail || contactCreate.invalidEmail) return;
      if (contactCreate.errorEmail) return;

      store.ui.commandMenu.context?.meta?.callback();
      handleClose();
    }
  };

  useModKey('Enter', () => {
    handleSubmit();
  });
  useKey('Escape', () => {
    handleSubmit();
  });

  return (
    <Command shouldFilter={false} label='Add contacts'>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100 max-h-[580px]'>
        <div className='flex items-center justify-between mb-2'>
          <h1 className='text-base font-medium'>Add a contact using their..</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <div className='text-sm flex flex-col gap-4'>
          <ButtonGroup className='flex items-center w-full'>
            <Button
              size='xs'
              leftIcon={<LinkedinOutline />}
              onClick={() => contactCreate.setType('linkedin')}
              data-inactive={contactCreate.getType !== 'linkedin'}
              className={cn('w-full', {
                selected: contactCreate.getType === 'linkedin',
              })}
            >
              LinkedIn
            </Button>
            <Button
              size='xs'
              leftIcon={<Signature />}
              onClick={() => contactCreate.setType('email')}
              data-inactive={contactCreate.getType !== 'email'}
              className={cn('w-full', {
                selected: contactCreate.getType === 'email',
              })}
            >
              Email
            </Button>
            <Button
              size='xs'
              leftIcon={<Signature />}
              dataTest='org-people-add-by-name'
              onClick={() => contactCreate.setType('name')}
              data-inactive={contactCreate.getType !== 'name'}
              className={cn('w-full', {
                selected: contactCreate.getType === 'name',
              })}
            >
              Name
            </Button>
          </ButtonGroup>
          <Input
            autoFocus
            ref={inputRef}
            variant='unstyled'
            placeholder={inputPlaceholder}
            dataTest='org-people-name-input'
            value={contactCreate.inputValue}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                handleClose(e);
              }
              e.stopPropagation();
            }}
            onChange={(e) => {
              contactCreate.setInputValue(e.target.value);

              if (contactCreate.inputValue) {
                contactCreate.getType === 'name' &&
                  contactCreate.validateName();
              }
              contactCreate.clearErrors();
            }}
          />
        </div>

        {contactCreate.getType === 'name' && (
          <p
            className={cn(
              'text-error-500 text-[12px] mt-0 opacity-0',
              contactCreate.invalidName && 'opacity-100',
            )}
          >
            Every hero needs a name
          </p>
        )}

        {contactCreate.getType === 'linkedin' && (
          <>
            <p
              className={cn(
                'text-error-500 text-[12px] mt-0 opacity-0',
                (contactCreate.emptyLinkedInUrl ||
                  contactCreate.errorLinkedIn ||
                  contactCreate.invalidLinkedInUrl) &&
                  'opacity-100',
              )}
            >
              {contactCreate.inputValue.length === 0
                ? 'Huston we have a blank...'
                : contactCreate.errorLinkedIn
                ? contactCreate.errorLinkedIn
                : 'Invalid LinkedIn URL'}
            </p>
          </>
        )}

        {contactCreate.getType === 'email' && (
          <>
            <p
              className={cn(
                'text-error-500 text-[12px] mt-0 opacity-0',
                (contactCreate.emptyEmail ||
                  contactCreate.errorEmail ||
                  contactCreate.invalidEmail) &&
                  'opacity-100',
              )}
            >
              {contactCreate.inputValue.length === 0
                ? 'Huston we have a blank...'
                : contactCreate.errorEmail
                ? contactCreate.errorEmail
                : 'Invalid email format'}
            </p>
          </>
        )}

        <div className='flex justify-between gap-3 mt-2'>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            onClick={handleClose}
            // onFocus={(e) => e.preventDefault()}
          >
            Cancel
          </Button>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            data-test='contact-actions-confirm-flow-change'
            leftSpinner={<Spinner size='sm' label='creating contacts' />}
            onClick={() => {
              handleSubmit();
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleSubmit();
              }
            }}
          >
            {confirmButtonPlaceholder}
          </Button>
        </div>
      </article>
    </Command>
  );
});
