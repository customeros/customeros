import React, { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { FlowStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const StopFlow = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const flow = store.flows.value.get(context.ids?.[0]);

  const confirmButtonRef = useRef<HTMLButtonElement>(null);

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    flow?.update((f) => {
      f.status = FlowStatus.Paused;

      return f;
    });
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>
            Stop flow '{flow?.value.name}'?
          </h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        {/* todo remove when bottom part is integrated with BE */}
        <p className='text-sm mt-2'>
          You currently have {flow?.value?.contacts?.length}{' '}
          {flow?.value?.contacts?.length === 1 ? 'contact' : 'contacts'} active
          in this flow.
        </p>
        {/* todo uncomment when be is ready */}
        {/*<div>*/}
        {/*  <label className='text-sm'>*/}
        {/*    <div className='mb-2'>*/}
        {/*      You currently have {} contacts active in this flow. How would you*/}
        {/*      like to handle them?*/}
        {/*    </div>*/}
        {/*    <RadioGroup*/}
        {/*      value={value}*/}
        {/*      name='last-touchpoint-date-before'*/}
        {/*      onValueChange={(val: 'my-view' | 'team-view') => setValue(val)}*/}
        {/*    >*/}
        {/*      <Radio value={'pause'}>*/}
        {/*        <span className='text-sm'>Pause and resume them later</span>*/}
        {/*      </Radio>*/}
        {/*      <Radio value={'end'}>*/}
        {/*        <span className='text-sm'>End them early</span>*/}
        {/*      </Radio>*/}
        {/*    </RadioGroup>*/}
        {/*  </label>*/}
        {/*</div>*/}
        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='error'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='flow-actions-confirm-stop'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Stop flow
          </Button>
        </div>
      </article>
    </Command>
  );
});
