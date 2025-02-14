import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Market } from '@graphql/types';
import { Combobox } from '@ui/form/Combobox';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Briefcase02 } from '@ui/media/icons/Briefcase02';
import { SelectOption } from '@shared/types/SelectOptions';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';
import { businessTypeOptions } from '@organization/components/Tabs/panels/AboutPanel/util';

interface BusinessTypeInputProps {
  id: string;
  dataTest?: string;
}

export const BusinessTypeInput = observer(
  ({ id, dataTest }: BusinessTypeInputProps) => {
    const store = useStore();
    const [inputValue, setInputValue] = useState('');
    const [isOpen, setIsOpen] = useState(false);
    const organization = store.organizations.value.get(id);

    const options = businessTypeOptions
      ?.map((user) => ({
        value: user.value,
        label: user.label,
      }))
      ?.sort((a, b) => a.label.localeCompare(b.label));

    const value = organization?.value?.market
      ? businessTypeOptions?.find((o) => o.value === organization.value.market)
      : null;

    const handleSelect = (option: SelectOption<Market>[]) => {
      const organization = store.organizations.value.get(id);

      if (!organization) return;
      organization.draft();

      const filteredOptions = option.filter((e) => e.value !== value?.value);

      organization.value.market = filteredOptions[0]?.value;
      organization.commit();
    };

    return (
      <>
        <Popover open={isOpen} onOpenChange={(open) => setIsOpen(open)}>
          <Tooltip align='start' label='Business type'>
            <PopoverTrigger className={cn('flex items-center ')}>
              <Briefcase02 className='text-gray-500 mr-3' />
              <div
                data-test={dataTest}
                className='flex flex-wrap gap-1 w-fit items-center'
              >
                {value ? (
                  <div className='text-sm'>{value.label}</div>
                ) : (
                  <span className='text-gray-400 text-sm'>Business type</span>
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
  },
);
