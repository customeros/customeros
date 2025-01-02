import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const DuplicateDomainInformationModal = observer(() => {
  const { ui } = useStore();

  const context = ui.commandMenu.context;

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  const handleConfirm = () => {
    if (!context.ids?.length || !context.property) return;

    // todo merge orgs

    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  const handleClose = () => {
    ui.commandMenu.toggle('DuplicateDomainInformationModal');
    ui.commandMenu.clearCallback();
  };

  // todo get domain, associatedOrgName, editedOrgName from context
  const domain = '';
  const associatedOrgName = '';
  const editedOrgName = '';

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100 cursor-default'>
        <div className='flex justify-between'>
          <h1 className='text-base font-semibold'>Domain already exists</h1>
          <div>
            <CommandCancelIconButton onClose={handleClose} />
          </div>
        </div>
        <p className='mt-1 text-sm'>
          <span className='font-medium mr-1'>{domain},</span>
          is already associated with another existing organization,
          <span className='font-medium mx-1'>{associatedOrgName},</span>. Would
          you like to merge {associatedOrgName} into {editedOrgName}?
        </p>

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton ref={closeButtonRef} onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            dataTest={'add-contact-to-flow-confirmation'}
            data-test='contact-actions-confirm-flow-change'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Merge orgs
          </Button>
        </div>
      </article>
    </Command>
  );
});
