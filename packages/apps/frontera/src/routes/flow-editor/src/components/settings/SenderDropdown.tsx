import React, { useRef, useState, useCallback } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowStore } from '@store/Flows/Flow.store.ts';

import { Combobox } from '@ui/form/Combobox';
import { SelectOption } from '@ui/utils/types';
import { User01 } from '@ui/media/icons/User01';
import { Button } from '@ui/form/Button/Button';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { useStore } from '@shared/hooks/useStore';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { components, OptionProps } from '@ui/form/Select/Select';
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from '@ui/overlay/Popover/Popover';

interface SenderDropdownProps {
  flowId: string;
}

export const SenderDropdown = observer(({ flowId }: SenderDropdownProps) => {
  const store = useStore();
  const contentRef = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const users = store.users.toArray();
  const flow = store.flows.value.get(flowId) as FlowStore;
  const selectedSenders =
    flow?.value.senders.map((e) => e.user?.id ?? '') ?? [];

  const options = users
    .filter(
      (e) =>
        !e.name.toLowerCase().includes('customeros') &&
        !selectedSenders.includes(e.id),
    )
    .sort((a, b) => {
      const aMailboxesCount = a?.value?.mailboxes.length ?? 0;
      const bMailboxesCount = b?.value?.mailboxes.length ?? 0;

      return bMailboxesCount - aMailboxesCount;
    })
    .map((user) => ({
      label: user?.name,
      value: user?.id,
    }));

  const handleSelect = (value: SelectOption) => {
    store.flowSenders.createFlowSender(flowId, value.value);
    setOpen(false);
  };

  const Option = useCallback(
    ({ children, ...props }: OptionProps) => {
      const id = (props?.data as SelectOption)?.value;
      const _user = store.users.value.get(id);

      return (
        <components.Option {...props}>
          <div className='flex w-full items-center'>
            <Avatar
              size='xs'
              textSize='xxs'
              name={_user?.name ?? 'Unnamed'}
              src={_user?.value?.profilePhotoUrl ?? ''}
              icon={<User01 className='text-gray-500 size-3' />}
              className={'w-5 h-5 min-w-5 mr-2 border border-gray-200'}
            />
            <div className='flex-1 text-sm'>{children}</div>
          </div>
        </components.Option>
      );
    },
    [selectedSenders],
  );

  return (
    <span>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger>
          <Button
            size='xs'
            variant='ghost'
            className='-ml-2'
            colorScheme='primary'
            dataTest='flow-add-senders'
            onClick={() => setOpen(true)}
            leftIcon={<PlusCircle className='text-inherit' />}
          >
            Add senders
          </Button>
        </PopoverTrigger>

        <PopoverContent align='start' ref={contentRef} className='w-[330px]'>
          <Combobox
            options={options}
            escapeClearsValue
            onChange={handleSelect}
            closeMenuOnSelect={false}
            placeholder='Search a sender...'
            components={{
              Option,
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                setOpen(false);
              }
            }}
          />
        </PopoverContent>
      </Popover>
    </span>
  );
});
