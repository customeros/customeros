import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { RemoveOrganizationDomainCase } from '@domain/usecases/command-menu/remove-organization-domain.usecase';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

const removeDomainCase = new RemoveOrganizationDomainCase();

export const RemoveDomain = observer(() => {
  const { ui, organizations } = useStore();
  const context = ui.commandMenu.context;
  const isPrimary = context.meta?.isPrimary;
  const organization = organizations.getById(context.ids?.[0] as string);

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });
  useEffect(() => {
    if (organization) {
      removeDomainCase.setEntity(organization);
      removeDomainCase.setDomain(context?.meta?.domain);
    }
  }, [organization?.id]);

  const handleConfirm = () => {
    removeDomainCase.submit();
    ui.commandMenu.setOpen(false);
  };

  const handleClose = () => {
    ui.commandMenu.toggle('ConfirmSingleFlowEdit');
    ui.commandMenu.clearCallback();
  };

  return (
    <Command shouldFilter={false}>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100 cursor-default'>
        <div className='flex justify-between'>
          <h1 className='text-base font-semibold'>
            Remove this {isPrimary ? 'domain' : 'subdomain'}
          </h1>
          <div>
            <CommandCancelIconButton onClose={handleClose} />
          </div>
        </div>
        <p className='mt-2 text-sm'>
          {isPrimary
            ? `Removing ${context?.meta?.domain}, will also remove any related
            subdomains`
            : `Removing ${context?.meta?.domain} will not remove its associated primary domain`}
        </p>

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='error'
            onClick={handleConfirm}
            dataTest={'remove-domain'}
            data-test='contact-actions-confirm-flow-change'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Remove
          </Button>
        </div>
      </article>
    </Command>
  );
});
