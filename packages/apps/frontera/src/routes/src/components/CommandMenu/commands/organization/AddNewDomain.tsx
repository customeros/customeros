import React, { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { AddOrganizationDomainCase } from '@domain/usecases/command-menu/add-organization-domain.usecase';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

const addNewDomainCase = new AddOrganizationDomainCase();

export const AddNewDomain = observer(() => {
  const { ui, organizations } = useStore();
  const context = ui.commandMenu.context;
  const [error, setError] = useState('');
  const organization = organizations.value.get(context.ids?.[0] as string);

  const inputRef = useRef<HTMLInputElement>(null);

  const handleConfirm = () => {
    if (addNewDomainCase.inputValue === '') {
      setError('Houston, we have a blank...');

      return;
    }
    addNewDomainCase.submit();
    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });
  useEffect(() => {
    if (organization) {
      addNewDomainCase.setEntity(organization);
    }
  }, [organization?.id]);

  useEffect(() => {
    if (inputRef.current) {
      inputRef.current?.focus();
    }
  }, [inputRef.current]);

  const handleClose = () => {
    ui.commandMenu.toggle('ConfirmSingleFlowEdit');
    ui.commandMenu.clearCallback();
  };

  return (
    <Command shouldFilter={false}>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100 cursor-default'>
        <div className='flex justify-between'>
          <h1 className='text-base font-semibold'>Add new domain</h1>
          <div>
            <CommandCancelIconButton onClose={handleClose} />
          </div>
        </div>
        <p className='mt-1 text-sm'>
          Domains that redirect will appear as subdomains of their primary
          domain
        </p>

        <div className={'mt-4'}>
          <Input
            autoFocus
            size={'sm'}
            ref={inputRef}
            placeholder='Organization’s domain'
            value={addNewDomainCase.inputValue}
            onChange={(e) => {
              addNewDomainCase.setInputValue(e.target.value);
            }}
            onKeyDownCapture={(e) => {
              if (e.key === '') {
                e.stopPropagation();
              }
            }}
            className={cn({
              'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                error,
            })}
          />
          {error && <p className='text-xs text-error-600'> {error}</p>}
        </div>

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            onClick={handleConfirm}
            dataTest={'add-contact-to-flow-confirmation'}
            data-test='contact-actions-confirm-flow-change'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Add domain
          </Button>
        </div>
      </article>
    </Command>
  );
});
