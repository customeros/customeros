import { observer } from 'mobx-react-lite';

import { XClose } from '@ui/media/icons/XClose';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { Command } from '@ui/overlay/CommandMenu';
import { useStore } from '@shared/hooks/useStore';

export const MergeConfirmationModal = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const handleClose = () => {
    store.ui.commandMenu.toggle('MergeConfirmationModal');
    store.ui.commandMenu.clearContext();
  };

  const handleConfirm = () => {
    const [primary, ...rest] = context.ids as string[];

    store.organizations.merge(
      primary,
      rest,
      store.ui.commandMenu?.context?.callback,
    );
    handleClose();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>
            Merge {context.ids?.length} organizations?
          </h1>
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
          >
            Merge
          </Button>
        </div>
      </article>
    </Command>
  );
});
