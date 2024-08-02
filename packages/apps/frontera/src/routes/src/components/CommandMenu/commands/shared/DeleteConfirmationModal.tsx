import { observer } from 'mobx-react-lite';

import { XClose } from '@ui/media/icons/XClose';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { Command } from '@ui/overlay/CommandMenu';
import { useStore } from '@shared/hooks/useStore';

export const DeleteConfirmationModal = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const handleClose = () =>
    store.ui.commandMenu.toggle('DeleteConfirmationModal');

  return (
    <Command label='Change Stage'>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold '>The big red button</h1>
          <IconButton
            size='xs'
            variant='ghost'
            icon={<XClose />}
            aria-label='cancel'
            onClick={handleClose}
          />
        </div>
        <p className='text-sm mt-1'>
          Are you feeling lucky? This is it. ‘Delete’ to go all in, ‘Cancel’ to
          chicken out.
        </p>
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
            onClick={() => {
              store.organizations.hide([context?.id as string]);
              handleClose();
            }}
          >
            Delete
          </Button>
        </div>
      </article>
    </Command>
  );
});
