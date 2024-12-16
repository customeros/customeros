import { useCallback } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { Combobox } from '@ui/form/Combobox';
import { Check } from '@ui/media/icons/Check';
import { SelectOption } from '@ui/utils/types';
import { User01 } from '@ui/media/icons/User01';
import { useStore } from '@shared/hooks/useStore';
import { components, OptionProps } from '@ui/form/Select';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { Popover, PopoverTrigger, PopoverContent } from '@ui/overlay/Popover';

interface UserCellProps {
  id: string;
}

export const UserCell = observer(({ id }: UserCellProps) => {
  const store = useStore();
  const mailbox = store.mailboxes.value.get(id);
  const user = mailbox?.user;
  const { open, onToggle } = useDisclosure();

  if (mailbox?.value?.scheduledEmails > 0) {
    return (
      <p
        tabIndex={0}
        role={'button'}
        className={cn('hover:text-gray-500', !user && 'text-gray-400')}
        onClick={() => {
          store.ui.commandMenu.setType('ActiveFlowUpdateInfo');
          store.ui.commandMenu.setContext({
            ...store.ui.commandMenu.context,
            meta: {
              type: 'senders',
            },
          });
          store.ui.commandMenu.setOpen(true);
        }}
      >
        {user?.name || 'Not set yet'}
      </p>
    );
  }

  return (
    <Popover open={open} onOpenChange={onToggle}>
      <PopoverTrigger>
        <p
          onDoubleClick={onToggle}
          className={cn('hover:text-gray-500', !user && 'text-gray-400')}
        >
          {user?.name || 'Not set yet'}
        </p>
      </PopoverTrigger>
      {open && <UserSelector id={id} onOpenChange={onToggle} />}
    </Popover>
  );
});

interface UserSelectorProps {
  id: string;
  onOpenChange: (v: boolean) => void;
}

const UserSelector = observer(({ id, onOpenChange }: UserSelectorProps) => {
  const store = useStore();
  const mailbox = store.mailboxes.value.get(id);
  const user = mailbox?.user;
  const users = store.users.toComputedArray((arr) => {
    return arr.filter(
      (e) =>
        Boolean(e.value.firstName) ||
        Boolean(e.value.lastName) ||
        Boolean(e.value.name),
    );
  });

  const options = users
    ?.map((user) => ({
      value: user.id,
      label: user.name,
    }))
    ?.sort((a, b) => a.label.localeCompare(b.label));

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
            <span className='flex-1'>{children}</span>
            {user?.id === _user?.id && <Check />}
          </div>
        </components.Option>
      );
    },
    [user?.id],
  );

  const handleSelect = (option: SelectOption) => {
    mailbox!.value.userId = option.value;
    mailbox?.commit();
    onOpenChange(false);
  };

  return (
    <PopoverContent className='w-[264px]'>
      <Combobox
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
            onOpenChange(false);
          }
        }}
      />
    </PopoverContent>
  );
});
