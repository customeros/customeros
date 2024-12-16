import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandCancelIconButton } from '@ui/overlay/CommandMenu';

export const ActiveFlowUpdateInfo = observer(() => {
  const store = useStore();

  const type = store.ui.commandMenu.context.meta?.type;

  const handleClose = () => {
    store.ui.commandMenu.clearContext();
    store.ui.commandMenu.setOpen(false);
  };
  const title = match(type)
    .with('steps', () => 'Steps can’t be added or removed')
    .with('senders', () => 'Senders can’t be changed in live flows')
    .otherwise(() => '');

  const descriptions = match(type)
    .with('steps', () => (
      <>
        Once a flow has gone live the first time, you can only edit content and
        basic settings.
        <p className='mt-2'>
          To add or remove a step, you should duplicate the flow or create a new
          one.
        </p>
      </>
    ))
    .with('senders', () => (
      <>
        This sender has scheduled actions. To change or remove the senders,
        you’ll first have to stop the flow.
      </>
    ))
    .otherwise(() => '');

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>{title}</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <p className='text-sm mt-2'>{descriptions}</p>
        <div className='flex justify-between gap-3 mt-6'>
          <Button
            size='sm'
            variant='outline'
            className='w-full'
            onClick={handleClose}
            data-test='flow-actions-confirm-stop'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleClose();
              }
            }}
          >
            Got it
          </Button>
        </div>
      </article>
    </Command>
  );
});
