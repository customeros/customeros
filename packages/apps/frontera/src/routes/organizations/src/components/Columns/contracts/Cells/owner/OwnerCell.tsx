import { useMemo, useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import { Select, getContainerClassNames } from '@ui/form/Select/Select';

interface OwnerProps {
  id: string;
}

export const OwnerCell = observer(({ id }: OwnerProps) => {
  const store = useStore();
  const contract = store.contracts.value.get(id);
  const [isEditing, setIsEditing] = useState(false);
  const opportunityId = useMemo(() => {
    return (
      contract?.value?.opportunities?.find((e) => e.internalStage === 'OPEN')
        ?.id || contract?.value?.opportunities?.[0]?.id
    );
  }, []);

  const opportunity = store.opportunities.value.get(opportunityId ?? '');

  const owner = opportunity?.value;

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

  // const open = () => {
  //   setIsEditing(true);
  //   store.ui.setIsEditingTableCell(true);
  // };

  const close = () => {
    setIsEditing(false);
    store.ui.setIsEditingTableCell(false);
  };

  const handleSelect = (option: SelectOption) => {
    const targetOwner = store.users.value.get(option?.value);

    opportunity?.update((value) => {
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
          // onDoubleClick={open}
          className={cn(
            value ? 'text-gray-700' : 'text-gray-400',
            'cursor-default',
          )}
        >
          {value?.label ?? 'No owner'}
        </p>
        {/*<IconButton*/}
        {/*  size='xxs'*/}
        {/*  onClick={open}*/}
        {/*  variant='ghost'*/}
        {/*  id='edit-button'*/}
        {/*  aria-label='edit owner'*/}
        {/*  className='edit-button opacity-0'*/}
        {/*  icon={<Edit03 className='text-gray-500 size-3' />}*/}
        {/*/>*/}
      </div>
    );
  }

  return (
    <Select
      size='xs'
      autoFocus
      isDisabled
      isClearable
      value={value}
      onBlur={close}
      defaultMenuIsOpen
      options={options}
      placeholder='Owner'
      backspaceRemovesValue
      openMenuOnClick={false}
      onChange={handleSelect}
      menuPortalTarget={document.body}
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          close();
        }
      }}
      classNames={{
        container: ({ isFocused }) =>
          getContainerClassNames('border-0 w-[164px]', undefined, {
            isFocused,
            size: 'xs',
          }),
      }}
    />
  );
});
