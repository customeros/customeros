import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { User } from '@graphql/types';
import { Combobox } from '@ui/form/Combobox';
import { Key01 } from '@ui/media/icons/Key01';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { SelectOption } from '@shared/types/SelectOptions';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OwnerProps {
  id: string;
  owner?: Owner;
  dataTest?: string;
}

export const OwnerInput = observer(({ id, owner, dataTest }: OwnerProps) => {
  const store = useStore();
  const [inputValue, setInputValue] = useState('');
  const [isOpen, setIsOpen] = useState(false);

  const users = store.users.tenantUsers.filter(
    (e) =>
      Boolean(e.value.firstName) ||
      Boolean(e.value.lastName) ||
      Boolean(e.value.name),
  );

  const options = users
    ?.map((user) => ({
      value: user.id,
      label: user.name,
    }))
    ?.sort((a, b) => a.label.localeCompare(b.label));

  const value = owner ? options?.find((o) => o.value === owner.id) : null;

  const handleSelect = (option: SelectOption[]) => {
    const organization = store.organizations.value.get(id);

    if (!organization) return;
    const filteredOptions = option.filter((e) => e.value !== value?.value);

    const targetOwner = store.users.value.get(filteredOptions[0]?.value);

    organization.value.owner = option[0]?.value ? targetOwner?.value : null;
    organization.commit();
  };

  return (
    <>
      <Popover open={isOpen} onOpenChange={(open) => setIsOpen(open)}>
        <Tooltip label='Owner' align='start' placement='top'>
          <PopoverTrigger className={cn('flex items-center')}>
            <Key01 className='mr-3 text-gray-500' />
            <div
              data-test={dataTest}
              className='flex flex-wrap  w-fit items-center'
            >
              {value ? (
                <div className='text-sm'>{value.label}</div>
              ) : (
                <span className='text-gray-400 text-sm'>Owner</span>
              )}
            </div>
          </PopoverTrigger>
        </Tooltip>
        <PopoverContent align='start' className='min-w-[264px] max-w-[320px]'>
          <Combobox
            isMulti
            value={value}
            options={options}
            placeholder='Owner'
            inputValue={inputValue}
            onInputChange={setInputValue}
            onChange={(newValue) => {
              handleSelect(newValue);
              setIsOpen(false);
            }}
            noOptionsMessage={({ inputValue }) => (
              <div className='text-gray-700 px-3 py-1 mt-0.5 rounded-md bg-grayModern-100 gap-1 flex items-center'>
                <span>{`No results matching "${inputValue}"`}</span>
              </div>
            )}
          />
        </PopoverContent>
      </Popover>
    </>
  );
});
