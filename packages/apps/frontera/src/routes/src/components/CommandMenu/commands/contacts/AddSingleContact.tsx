import { useRef, useMemo, useEffect, MouseEvent, KeyboardEvent } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { CreateContactUsecase } from '@domain/usecases/contact-details/create-contact.usecase';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
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

  const usecase = useMemo(() => new CreateContactUsecase(), []);

  const currentConfig =
    INPUT_CONFIGS.find((config) => config.type === usecase.type) ??
    INPUT_CONFIGS[0];

  const confirmButtonText =
    usecase.type === 'name' ? 'Add contact' : 'Add & enrich';

  const handleClose = (
    e?:
      | KeyboardEvent
      | KeyboardEvent<HTMLButtonElement | HTMLInputElement>
      | MouseEvent,
  ) => {
    e?.stopPropagation();
    e?.preventDefault();
    usecase.clearState();
    store.ui.commandMenu.toggle('AddSingleContact');
    store.ui.commandMenu.clearContext();
  };

  const handleSubmit = async () => {
    usecase.setOrganizationId(store.ui.commandMenu.context?.ids?.[0] as string);
    await usecase.submit();
  };

  useEffect(() => {
    const email = store.ui.commandMenu.context?.meta?.email;

    if (email?.length) {
      usecase.setType('email');
      usecase.setInputValue(email);
    }
  }, [store.ui.commandMenu.context?.meta]);

  useModKey('Enter', handleSubmit);
  useKey('Escape', () => handleClose());

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
                onClick={() => usecase.setType(config.type)}
                data-inactive={usecase.type !== config.type}
                className={cn('w-full', {
                  selected: usecase.type === config.type,
                })}
                dataTest={
                  config.type === 'name' ? 'org-people-add-by-name' : undefined
                }
              >
                {config.label}
              </Button>
            ))}
          </ButtonGroup>

          <Input
            autoFocus
            ref={inputRef}
            variant='unstyled'
            value={usecase.inputValue}
            dataTest='org-people-name-input'
            placeholder={currentConfig.placeholder}
            onChange={(e) => {
              usecase.setInputValue(e.target.value);
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                handleClose(e);
              }

              if (e.key === 'Enter') {
                handleSubmit();
              }
              e.stopPropagation();
            }}
          />

          {!CreateContactUsecase.isBrowserExtensionEnabled &&
            usecase.type === 'linkedin' && (
              <div className='flex justify-between bg-success-50 rounded-md px-2 py-1 mb-4'>
                <div className='flex items-center gap-2 text-success-700 text-sm'>
                  <Icon name='zap' />
                  <span>
                    <span>Prospect faster.</span>{' '}
                    <a
                      target='_blank'
                      className='underline hover:text-success-900'
                      href='https://chromewebstore.google.com/detail/customeros/khmdccjeodppdldkgifcnkndemjpfoml'
                    >
                      Get our Chrome extension
                    </a>
                    <span>.</span>
                  </span>
                </div>
              </div>
            )}
        </div>

        <p className={cn('text-error-500 text-[12px] mt-0')}>
          {usecase.type === 'name' &&
            usecase.currentError.isEmpty &&
            'Every hero needs a name'}

          {usecase.type !== 'name' && usecase.currentError.isEmpty
            ? 'Huston we have a blank...'
            : usecase.currentError.message
            ? usecase.currentError.message
            : usecase.currentError.isInvalid
            ? `Invalid ${usecase.type} format`
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
            isLoading={usecase.isLoading}
            loadingText={'Adding contact...'}
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
