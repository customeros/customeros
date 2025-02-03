import { useRef, useEffect } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { CreateContactUsecase } from '@domain/usecases/contact-details/create-contact.usecase';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { Mail01 } from '@ui/media/icons/Mail01';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { useModKey } from '@shared/hooks/useModKey';
import { Signature } from '@ui/media/icons/Signature';
import { LinkedinOutline } from '@ui/media/icons/LinkedinOutline';
import { Command, CommandCancelIconButton } from '@ui/overlay/CommandMenu';

const contactCreate = new CreateContactUsecase();

interface InputConfig {
  label: string;
  icon: JSX.Element;
  placeholder: string;
  type: 'linkedin' | 'email' | 'name';
}

const INPUT_CONFIGS: InputConfig[] = [
  {
    type: 'linkedin',
    placeholder: 'linkedin.com/in/johnlemon',
    icon: <LinkedinOutline className='text-inherit' />,
    label: 'LinkedIn',
  },
  {
    type: 'email',
    placeholder: 'john@heyjude.band',
    icon: <Mail01 className='text-inherit' />,
    label: 'Email',
  },
  {
    type: 'name',
    placeholder: 'First and last name',
    icon: <Signature className='text-inherit' />,
    label: 'Name',
  },
];

export const AddSingleContact = observer(() => {
  const store = useStore();
  const inputRef = useRef<HTMLInputElement>(null);

  const currentConfig =
    INPUT_CONFIGS.find((config) => config.type === contactCreate.type) ??
    INPUT_CONFIGS[0];

  const confirmButtonText =
    contactCreate.type === 'name' ? 'Add contact' : 'Add & enrich';

  const handleClose = (e?: KeyboardEvent) => {
    e?.stopPropagation();
    e?.preventDefault();
    contactCreate.clearState();
    store.ui.commandMenu.toggle('AddSingleContact');
    store.ui.commandMenu.clearContext();
  };

  const handleSubmit = async () => {
    contactCreate.setOrganizationId(
      store.ui.commandMenu.context?.ids?.[0] as string,
    );
    await contactCreate.submit();
  };

  useEffect(() => {
    const email = store.ui.commandMenu.context?.meta?.email;

    if (email?.length) {
      contactCreate.setType('email');
      contactCreate.setInputValue(email);
    }
  }, [store.ui.commandMenu.context?.meta]);

  useModKey('Enter', handleSubmit);
  useKey('Escape', handleClose);

  return (
    <Command shouldFilter={false} label='Add contacts'>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100 max-h-[580px]'>
        <div className='flex items-center justify-between mb-2'>
          <h1 className='text-base font-medium'>
            Add a contact using their...
          </h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <div className='text-sm flex flex-col gap-4'>
          <ButtonGroup className='flex items-center w-full'>
            {INPUT_CONFIGS.map((config) => (
              <Button
                size='xs'
                key={config.type}
                leftIcon={config.icon}
                onClick={() => contactCreate.setType(config.type)}
                data-inactive={contactCreate.type !== config.type}
                dataTest={
                  config.type === 'name' ? 'org-people-add-by-name' : undefined
                }
                className={cn('w-full', {
                  selected: contactCreate.type === config.type,
                })}
              >
                {config.label}
              </Button>
            ))}
          </ButtonGroup>

          <Input
            autoFocus
            ref={inputRef}
            variant='unstyled'
            dataTest='org-people-name-input'
            value={contactCreate.inputValue}
            placeholder={currentConfig.placeholder}
            onChange={(e) => {
              contactCreate.setInputValue(e.target.value);
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                handleClose(e);
              }
              e.stopPropagation();
            }}
          />
        </div>

        <p className={cn('text-error-500 text-[12px] mt-0')}>
          {contactCreate.type === 'name' &&
            contactCreate.currentError.isEmpty &&
            'Every hero needs a name'}

          {contactCreate.type !== 'name' && contactCreate.currentError.isEmpty
            ? 'Huston we have a blank...'
            : contactCreate.currentError.message
            ? contactCreate.currentError.message
            : contactCreate.currentError.isInvalid
            ? `Invalid ${contactCreate.type} format`
            : ''}
        </p>

        <div className='flex justify-between gap-3 mt-2'>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            onClick={handleClose}
          >
            Cancel
          </Button>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            onClick={handleSubmit}
            loadingText={'Adding contact...'}
            isLoading={contactCreate.isLoading}
            dataTest='confirm-contact-creation'
            leftSpinner={<Spinner size='sm' label='creating contacts' />}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleSubmit();
              }
            }}
          >
            {confirmButtonText}
          </Button>
        </div>
      </article>
    </Command>
  );
});
