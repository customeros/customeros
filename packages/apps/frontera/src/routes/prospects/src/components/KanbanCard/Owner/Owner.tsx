import { useRef, useState, useCallback } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Combobox } from '@ui/form/Combobox';
import { Check } from '@ui/media/icons/Check';
import { SelectOption } from '@ui/utils/types';
import { User01 } from '@ui/media/icons/User01';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { components, OptionProps } from '@ui/form/Select/Select';
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from '@ui/overlay/Popover/Popover';

interface OwnerProps {
  opportunityId: string;
  ownerId?: string | null;
}

export const Owner = observer(({ ownerId, opportunityId }: OwnerProps) => {
  const store = useStore();
  const contentRef = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const user = store.users.value.get(ownerId ?? '');
  const users = store.users.toArray();

  const options = users.map((user) => ({
    label: user?.name,
    value: user?.id,
  }));

  const handleSelect = (value: SelectOption) => {
    const user = store.users.value.get(value.value);
    const opportunity = store.opportunities.value.get(opportunityId);

    if (!user || !opportunity) return;

    opportunity?.update((value) => {
      if (!value.owner) {
        Object.assign(value, { owner: user.value });

        return value;
      }

      Object.assign(value.owner, user.value);

      return value;
    });
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
              className={cn('w-5 h-5 min-w-5 mr-2', 'border border-gray-200')}
            />
            <span className='flex-1'>{children}</span>
            {user?.id === _user?.id && <Check />}
          </div>
        </components.Option>
      );
    },
    [user?.id],
  );

  return (
    <Tooltip label={!open ? user?.name : undefined}>
      <span>
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger>
            <Avatar
              size='xs'
              textSize='xxs'
              name={user?.name ?? 'Unnamed'}
              src={user?.value?.profilePhotoUrl ?? ''}
              icon={<User01 className='text-gray-500 size-3' />}
              className={'w-5 h-5 min-w-5 border border-gray-200'}
            />
          </PopoverTrigger>

          <PopoverContent ref={contentRef} className='w-[264px]'>
            <Combobox
              size='xs'
              options={options}
              escapeClearsValue
              onChange={handleSelect}
              closeMenuOnSelect={false}
              placeholder='Assign owner...'
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
    </Tooltip>
  );
});
