import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const DuplicateFlow = observer(() => {
  const [isSaving, setIsSaving] = useState(false);
  const { flows, ui } = useStore();

  const [flowName, setFlowName] = useState('');
  const context = ui.commandMenu.context;

  const flow = flows.value.get(context.ids?.[0]);

  const handleConfirm = () => {
    setIsSaving(true);

    if (!flow?.id) {
      ui.toastError('Flow not found', 'flow-duplication-error');
      ui.commandMenu.setOpen(false);

      return;
    }

    const duplicateName = !flowName.trim()?.length
      ? `${flow?.value?.name} - Take two`
      : flowName;

    flows.duplicate(duplicateName, flow.id, {
      onSuccess: (id) => {
        setIsSaving(false);
        ui.toastSuccess(
          `${duplicateName} created`,
          `flow-duplication-success-${id}`,
        );
        ui.commandMenu.toggle('DuplicateFlow');
      },
    });
  };

  return (
    <Command label={`Duplicate flow`}>
      <div className='p-6 pb-4 flex flex-col gap-1 '>
        <div className='flex items-center justify-between'>
          <h1
            className='text-base font-semibold'
            data-test='create-new-flow-modal-title'
          >
            Duplicate {flow?.value.name}
          </h1>

          <CommandCancelIconButton
            onClose={() => {
              ui.commandMenu.setOpen(false);
              ui.commandMenu.clearContext();
            }}
          />
        </div>

        <p className='text-sm'>
          This will create a new flow with the same blocks and content. Existing
          contacts from the original flow will not be included in the new one.
        </p>
      </div>

      <div className='pr-6 pl-6 pb-6 flex flex-col gap-2 '>
        <Input
          autoFocus
          size='sm'
          id='flowName'
          value={flowName}
          variant='unstyled'
          placeholder='Flow name'
          dataTest='duplicated-flow-name'
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
            ui.commandMenu.setOpen(false);
            ui.commandMenu.clearContext();
          }}
        />

        <Button
          className='w-full'
          isLoading={isSaving}
          colorScheme='primary'
          onClick={handleConfirm}
          loadingText={'Duplicating flow…'}
          data-test='confirm-duplicate-flow'
        >
          Duplicate flow
        </Button>
      </div>
    </Command>
  );
});
