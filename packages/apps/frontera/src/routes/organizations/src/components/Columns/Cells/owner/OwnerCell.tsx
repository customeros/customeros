import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { User } from '@graphql/types';
import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Select, getContainerClassNames } from '@ui/form/Select/Select';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OwnerProps {
  id: string;
  owner?: Owner;
}

export const OwnerCell = observer(({ id, owner }: OwnerProps) => {
  const store = useStore();
  const organization = store.organizations.value.get(id);
  const [isEditing, setIsEditing] = useState(false);

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

  const value = owner ? options?.find((o) => o.value === owner.id) : null;

  const open = () => {
    setIsEditing(true);
    store.ui.setIsEditingTableCell(true);
  };

  const close = () => {
    setIsEditing(false);
    store.ui.setIsEditingTableCell(false);
  };

  const handleSelect = (option: SelectOption) => {
    const targetOwner = store.users.value.get(option?.value);

    organization?.update((value) => {
      if (!option || !option?.value) {
        value.owner = null;
      } else {
        value.owner = targetOwner?.value;
      }

      return value;
    });
  };

  if (!isEditing) {
    return (
      <div className='flex w-full gap-1 items-center [&_.edit-button]:hover:opacity-100'>
        <p
          className={cn(
            value ? 'text-gray-700' : 'text-gray-400',
            'cursor-default',
          )}
          onDoubleClick={open}
        >
          {value?.label ?? 'Owner'}
        </p>
        <IconButton
          className='edit-button opacity-0'
          aria-label='edit owner'
          size='xxs'
          variant='ghost'
          id='edit-button'
          onClick={open}
          icon={<Edit03 className='text-gray-500 size-3' />}
        />
      </div>
    );
  }

  return (
    <Select
      size='xs'
      isClearable
      value={value}
      placeholder='Owner'
      autoFocus
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          close();
        }
      }}
      onBlur={close}
      defaultMenuIsOpen
      menuPortalTarget={document.body}
      backspaceRemovesValue
      openMenuOnClick={false}
      onChange={handleSelect}
      options={options}
      classNames={{
        container: ({ isFocused }) =>
          getContainerClassNames('border-0 w-[164px]', {
            isFocused,
            size: 'xs',
          }),
      }}
    />
  );
});
