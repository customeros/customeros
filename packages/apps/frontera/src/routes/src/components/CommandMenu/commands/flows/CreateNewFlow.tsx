import React, { useState } from 'react';

import { useDidMount } from 'rooks';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const CreateNewFlow = observer(() => {
  const store = useStore();
  const [allowSubmit, setAllowSubmit] = useState(false);
  const { flows } = useStore();

  const [sequenceName, setSequenceName] = useState('');

  useDidMount(() => {
    setTimeout(() => {
      setAllowSubmit(true);
    }, 100);
  });

  const handleConfirm = () => {
    if (!allowSubmit) return;
    setAllowSubmit(false);

    flows.create({
      name: sequenceName,
      description: '',
    });

    store.ui.commandMenu.toggle('CreateNewFlow');
  };

  return (
    <Command label={`Rename `}>
      <div className='p-6 pb-4 flex flex-col gap-1 '>
        <div className='flex items-center justify-between'>
          <h1
            className='text-base font-semibold'
            data-test='create-new-flow-modal-title'
          >
            Create new flow
          </h1>

          <CommandCancelIconButton
            onClose={() => {
              store.ui.commandMenu.setOpen(false);
            }}
          />
        </div>
      </div>

      <div className='pr-6 pl-6 pb-6 flex flex-col gap-2 '>
        <Input
          autoFocus
          id='sequenceName'
          variant='unstyled'
          value={sequenceName}
          placeholder='Flow name'
          dataTest='create-new-flow-name'
          onChange={(e) => {
            setSequenceName(e.target.value);
          }}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              handleConfirm();
            }
          }}
        />
      </div>

      <div className='flex w-full gap-2 pl-6 pr-6 pb-6'>
        <CommandCancelButton
          dataTest='cancel-create-new-flow'
          onClose={() => {
            store.ui.commandMenu.setOpen(false);
          }}
        />

        <Button
          className='w-full'
          colorScheme='primary'
          onClick={handleConfirm}
          data-test='confirm-create-new-flow'
        >
          Create
        </Button>
      </div>
    </Command>
  );
});
