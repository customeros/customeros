import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

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
  const [isSaving, setIsSaving] = useState(false);
  const { flows } = useStore();
  const navigate = useNavigate();

  const [flowName, setFlowName] = useState('');

  const handleConfirm = () => {
    setIsSaving(true);

    flows.create(flowName, {
      onSuccess: (id) => {
        setIsSaving(false);

        navigate(`/flow-editor/${id}`);
        store.ui.commandMenu.toggle('CreateNewFlow');
      },
    });
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
          id='flowName'
          value={flowName}
          variant='unstyled'
          placeholder='Flow name'
          dataTest='create-new-flow-name'
          onChange={(e) => {
            setFlowName(e.target.value);
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
          isLoading={isSaving}
          colorScheme='primary'
          onClick={handleConfirm}
          loadingText={'Creating…'}
          data-test='confirm-create-new-flow'
        >
          Create
        </Button>
      </div>
    </Command>
  );
});
