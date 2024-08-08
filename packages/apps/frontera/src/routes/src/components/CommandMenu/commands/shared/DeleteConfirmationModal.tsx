import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { OpportunityStore } from '@store/Opportunities/Opportunity.store';
import { OrganizationStore } from '@store/Organizations/Organization.store.ts';

import { XClose } from '@ui/media/icons/XClose';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { Command } from '@ui/overlay/CommandMenu';
import { useStore } from '@shared/hooks/useStore';

export const DeleteConfirmationModal = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = match(context.entity)
    .returnType<OpportunityStore | OrganizationStore | undefined>()
    .with('Opportunity', () => store.opportunities.value.get(context.ids?.[0]))
    .with('Organization', () => store.organizations.value.get(context.ids?.[0]))
    .otherwise(() => undefined);

  const handleClose = () => {
    store.ui.commandMenu.toggle('DeleteConfirmationModal');
  };

  const handleConfirm = () => {
    match(context.entity)
      .with('Organization', () => {
        store.organizations.hide(context.ids as string[]);
      })
      .with('Organizations', () => {
        store.organizations.hide(context.ids as string[]);
      })
      .with('Opportunity', () => {
        store.opportunities.archive(context.ids?.[0]);
      })
      .otherwise(() => {});

    handleClose();
  };

  const title = match(context.entity)
    .with(
      'Organization',
      () => `Archive ${(entity as OrganizationStore)?.value.name}?`,
    )
    .with(
      'Organizations',
      () => `Archive ${context.ids?.length} organizations?`,
    )
    .with(
      'Opportunity',
      () => `Archive ${(entity as OpportunityStore)?.value.name}`,
    )
    .otherwise(() => `Archive selected ${context.entity?.toLowerCase()}`);

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>{title}</h1>
          <IconButton
            size='xs'
            variant='ghost'
            icon={<XClose />}
            aria-label='cancel'
            onClick={handleClose}
          />
        </div>

        <div className='flex justify-between gap-3 mt-6'>
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
            colorScheme='error'
            onClick={handleConfirm}
            data-test='org-actions-confirm-archive'
          >
            Archive
          </Button>
        </div>
      </article>
    </Command>
  );
});
