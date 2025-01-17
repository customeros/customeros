import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { Combobox } from '@ui/form/Combobox';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import { Tag, TagLabel, TagCloseButton } from '@ui/presentation/Tag';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

interface TagsProps {
  dataTest?: string;
  className?: string;
  placeholder?: string;
  value: SelectOption[];
  options: SelectOption[];
  inputPlaceholder?: string;
  leftAccessory?: React.ReactNode;
  onCreate: (value: string) => void;
  onChange: (selection: SelectOption[]) => void;
}

export const Tags = observer(
  ({
    value,
    options,
    onChange,
    onCreate,
    className,
    placeholder,
    leftAccessory,
    inputPlaceholder,
    dataTest,
  }: TagsProps) => {
    const [inputValue, setInputValue] = useState('');
    const store = useStore();

    const handleClear = (id: string) => {
      onChange?.(value.filter((o) => o.value !== id));
    };

    const foundOption = options.some((o) =>
      o.label.toLowerCase().includes(inputValue.toLowerCase()),
    );

    return (
      <Popover>
        <PopoverTrigger className={cn('flex items-center', className)}>
          {leftAccessory}
          <div
            data-test={dataTest}
            className='flex flex-wrap gap-1 w-fit items-center'
          >
            {value.length ? (
              value.map((option) => {
                const tag = store.tags.getById(option.value)?.value;

                return (
                  <Tag
                    size={'md'}
                    variant='subtle'
                    key={option.value}
                    colorScheme={
                      // eslint-disable-next-line @typescript-eslint/no-explicit-any
                      (tag?.colorCode as unknown as any) ?? 'grayModern'
                    }
                  >
                    <TagLabel>{option.label}</TagLabel>
                    <TagCloseButton
                      onClick={(e) => {
                        e.stopPropagation();
                        handleClear(option.value);
                      }}
                    />
                  </Tag>
                );
              })
            ) : (
              <span className='text-gray-400'>{placeholder}</span>
            )}
          </div>
        </PopoverTrigger>
        <PopoverContent align='start' className='min-w-[264px] max-w-[320px]'>
          <Combobox
            isMulti
            value={value}
            options={options}
            onChange={onChange}
            inputValue={inputValue}
            onInputChange={setInputValue}
            placeholder={inputPlaceholder}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !foundOption) {
                onCreate?.(inputValue);
                setInputValue('');
              }
            }}
            noOptionsMessage={({ inputValue }) => (
              <div
                className='text-gray-700 px-3 py-1 mt-0.5 rounded-md bg-grayModern-100 gap-1 flex items-center'
                onClick={() => {
                  onCreate?.(inputValue);
                  setInputValue('');
                }}
              >
                <Plus />
                <span>{`Create "${inputValue}"`}</span>
              </div>
            )}
          />
        </PopoverContent>
      </Popover>
    );
  },
);
