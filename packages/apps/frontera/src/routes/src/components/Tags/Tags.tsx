import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { Plus } from '@ui/media/icons/Plus';
import { Combobox } from '@ui/form/Combobox';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { SelectOption } from '@shared/types/SelectOptions';
import { Tag, TagLabel, TagRightButton } from '@ui/presentation/Tag';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover.tsx';

interface TagsProps {
  dataTest?: string;
  className?: string;
  inputValue: string;
  placeholder?: string;
  value: SelectOption[];
  options: SelectOption[];
  inputPlaceholder?: string;
  leftAccessory?: React.ReactNode;
  onCreate: (value: string) => void;
  setInputValue: (value: string) => void;
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
    inputValue,
    setInputValue,
  }: TagsProps) => {
    const store = useStore();

    const handleClear = (id: string) => {
      onChange?.(value.filter((o) => o.value !== id));
    };

    return (
      <Popover>
        <Tooltip align='start' label='Company tags'>
          <PopoverTrigger className={cn('flex items-center w-full', className)}>
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
                      <TagRightButton
                        onClick={(e) => {
                          e.stopPropagation();
                          handleClear(option.value);
                        }}
                      >
                        <Icon name={'x-close'} />
                      </TagRightButton>
                    </Tag>
                  );
                })
              ) : (
                <span className='text-gray-400 text-sm'>{placeholder}</span>
              )}
            </div>
          </PopoverTrigger>
        </Tooltip>
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
              if (e.key === 'Enter' && !options.length) {
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
