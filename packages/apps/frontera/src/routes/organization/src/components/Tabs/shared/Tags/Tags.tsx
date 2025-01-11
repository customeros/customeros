import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { Combobox } from '@ui/form/Combobox';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { Tag01 } from '@ui/media/icons/Tag01.tsx';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import { CreatableSelectProps } from '@ui/form/CreatableSelect';
import { Tag, TagLabel, TagCloseButton } from '@ui/presentation/Tag';
import { EntityType } from '@shared/types/__generated__/graphql.types';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';

interface TagsProps {
  dataTest?: string;
  placeholder: string;
  value: SelectOption[];
  size?: CreatableSelectProps['size'];
  onCreateOption?: (value: string) => void;
  onChange: (value: SelectOption[]) => void;
}

export const Tags = observer(
  ({ dataTest, placeholder, onCreateOption, value, onChange }: TagsProps) => {
    const store = useStore();
    const [inputValue, setInputValue] = useState('');

    const options = store.tags
      ? store.tags
          .getByEntityType(EntityType.Contact)
          .filter((t) => t.value.name !== '')
          .map(
            (tag) =>
              ({
                value: tag.value.metadata.id,
                label: tag.value.name,
              } as SelectOption),
          )
      : [];

    const foundOption = options.some((o) =>
      o.label.toLowerCase().includes(inputValue.toLowerCase()),
    );

    const handleClear = (id: string) => {
      onChange?.(value.filter((o) => o.value !== id));
    };

    return (
      <>
        <Popover>
          <PopoverTrigger className={cn('flex items-center')}>
            <div className='flex items-center gap-2 text-sm mr-[75px] text-gray-500'>
              <Tag01 className='mt-[1px] text-gray-500' />
              Tags
            </div>
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
                <span className='text-gray-400 text-sm'>{placeholder}</span>
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
              placeholder={placeholder}
              onInputChange={setInputValue}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !foundOption) {
                  onCreateOption?.(inputValue);
                  setInputValue('');
                }
              }}
              noOptionsMessage={({ inputValue }) => (
                <div
                  className='text-gray-700 px-3 py-1 mt-0.5 rounded-md bg-grayModern-100 gap-1 flex items-center'
                  onClick={() => {
                    onCreateOption?.(inputValue);
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
      </>
    );
  },
);
