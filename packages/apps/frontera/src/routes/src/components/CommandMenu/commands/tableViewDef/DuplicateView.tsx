import React, { useRef, useState, useEffect } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { TableViewType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { Button } from '@ui/form/Button/Button.tsx';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const DuplicateView = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const [value, setValue] = useState<'my-view' | 'team-view'>('team-view');
  const [name, setName] = useState<string>('');
  const tableViewDef = store.tableViewDefs.getById(context.ids?.[0]);

  const title = match(tableViewDef?.value?.tableType)
    .with(TableViewType.Invoices, () => `${tableViewDef?.value.name} Invoices`)
    .otherwise(() => tableViewDef?.value.name);

  const confirmButtonRef = useRef<HTMLButtonElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleClose = () => {
    store.ui.commandMenu.toggle('DuplicateView');
    store.ui.commandMenu.clearCallback();
  };

  const handleConfirm = () => {
    store.tableViewDefs.createFavorite({
      id: store.ui.commandMenu.context.ids[0],
      name,
      isShared: value === 'team-view',
    });
    store.ui.commandMenu.toggle('DuplicateView');
    store.ui.commandMenu.clearCallback();
  };

  useEffect(() => {
    setTimeout(() => {
      inputRef.current?.focus();
    }, 0);
  }, []);

  return (
    <Command>
      <article className='relative w-full p-6 flex flex-col border-b border-b-gray-100'>
        <div className='flex items-center justify-between'>
          <h1 className='text-base font-semibold'>Duplicate '{title}'</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <Input
          value={name}
          ref={inputRef}
          className='my-4'
          variant='unstyled'
          placeholder='View name'
          onChange={(e) => setName(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Escape') {
              handleClose();
            }
          }}
        />

        <div>
          <label className='text-sm'>
            <div className='mb-2  font-medium'>Duplicate to...</div>
            <RadioGroup
              value={value}
              name='last-touchpoint-date-before'
              onValueChange={(val: 'my-view' | 'team-view') => setValue(val)}
            >
              <Radio value={'team-view'}>
                <span className='text-sm'>Team views</span>
              </Radio>
              <Radio value={'my-view'}>
                <span className='text-sm'>My views</span>
              </Radio>
            </RadioGroup>
          </label>
        </div>
        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            ref={confirmButtonRef}
            onClick={handleConfirm}
            data-test='org-actions-confirm-archive'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Duplicate
          </Button>
        </div>
      </article>
    </Command>
  );
});
