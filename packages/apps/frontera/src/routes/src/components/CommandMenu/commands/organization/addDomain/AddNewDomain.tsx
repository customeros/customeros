import React, { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { AddOrganizationDomainCase } from '@domain/usecases/command-menu/add-organization-domain.usecase.ts';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';
import { DuplicateDomainInformation } from '@shared/components/CommandMenu/commands/organization/addDomain/DuplicateDomainInformationModal.tsx';

const addNewDomainCase = new AddOrganizationDomainCase();

export const AddNewDomain = observer(() => {
  const { ui, organizations } = useStore();
  const context = ui.commandMenu.context;
  const organization = organizations.value.get(context.ids?.[0] as string);
  const [showDuplicateInfo, setShowDuplicateInfo] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleConfirm = async () => {
    addNewDomainCase.resetValidation();
    addNewDomainCase.checkIfEmpty();

    await addNewDomainCase.validateDomain();

    if (addNewDomainCase.error) {
      if (addNewDomainCase.associatedOrgId) {
        setShowDuplicateInfo(true);
      }

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
    addNewDomainCase.reset();
    ui.commandMenu.toggle('AddNewDomain');
    ui.commandMenu.clearCallback();
  };

  return (
    <Command shouldFilter={false} className={'!w-auto'}>
      <article
        className={cn(
          'relative w-full p-6 flex flex-col border-b border-b-gray-100 cursor-default bg-white',
          'transition transform duration-500 ease-in-out origin-top translate-y-0 opacity-100 scale-100 z-50',
          {
            'absolute bg-white rounded-md scale-[0.80] opacity-90 -translate-y-2 transition transform duration-500 z-1':
              showDuplicateInfo,
          },
        )}
      >
        <>
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
              size={'sm'}
              ref={inputRef}
              placeholder='Organization’s domain'
              value={addNewDomainCase.inputValue}
              onChange={(e) => {
                addNewDomainCase.setInputValue(e.target.value);
              }}
              onKeyDownCapture={(e) => {
                e.stopPropagation();

                if (e.key === 'Enter') {
                  handleConfirm();
                }
              }}
              className={cn({
                'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                  addNewDomainCase.error && !addNewDomainCase.associatedOrgId,
              })}
            />
            {addNewDomainCase.error && !addNewDomainCase.associatedOrgId && (
              <p className='text-xs text-error-600'>
                {' '}
                {addNewDomainCase.error}
              </p>
            )}
          </div>

          <div className='flex justify-between gap-3 mt-6'>
            <CommandCancelButton onClose={handleClose} />

            <Button
              size='sm'
              variant='outline'
              className='w-full'
              colorScheme='primary'
              onClick={handleConfirm}
              loadingText={'Adding domain...'}
              isLoading={addNewDomainCase.isValidating}
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
        </>
      </article>

      <DuplicateDomainInformation
        onClose={handleClose}
        isOpen={showDuplicateInfo}
        domain={addNewDomainCase.inputValue}
        associatedOrgId={addNewDomainCase.associatedOrgId || ''}
      />
    </Command>
  );
});
