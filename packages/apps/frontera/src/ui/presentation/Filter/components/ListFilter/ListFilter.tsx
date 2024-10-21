import {
  useRef,
  useMemo,
  useState,
  useEffect,
  useCallback,
  ComponentType,
} from 'react';

import { flags } from '@ui/media/flags';
import { Avatar } from '@ui/media/Avatar';
import { Combobox } from '@ui/form/Combobox';
import { Check } from '@ui/media/icons/Check';
import { Button } from '@ui/form/Button/Button';
import { User01 } from '@ui/media/icons/User01';
import { components, OptionProps } from '@ui/form/Select/Select';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';
import {
  FlowStatus,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

import { handleOperatorName, handlePropertyPlural } from '../../utils/utils';
interface GroupedOption {
  readonly label: string;
  readonly options: { id: string; label: string }[];
}

interface ListFilterProps {
  filterName: string;
  operatorName: string;
  filterValue: string[];
  onMultiSelectChange: (ids: string[]) => void;
  groupOptions?: { label: string; options: { id: string; label: string }[] }[];
  options: {
    id: string;
    label: string;
    avatar?: string;
    isArchived?: FlowStatus;
  }[];
}

export const ListFilter = ({
  onMultiSelectChange,
  options: _options,
  filterValue,
  filterName,
  groupOptions,
  operatorName,
}: ListFilterProps) => {
  const [selectedIds, setSelectedIds] = useState<string[]>(filterValue ?? []);
  const [isOpen, setIsOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const debouncedOnMultiSelectChange = useCallback(
    (ids: string[]) => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      timeoutRef.current = setTimeout(() => {
        onMultiSelectChange(ids);
      }, 300);
    },
    [onMultiSelectChange],
  );

  useEffect(() => {
    if (!filterValue) {
      if (filterName) {
        setTimeout(() => {
          setIsOpen(true);
        }, 50);
      }
    }
  }, [filterName]);

  const handleItemClick = (id: string) => {
    let newSelectedIds;

    if (selectedIds.includes(id)) {
      newSelectedIds = selectedIds.filter((selectedId) => selectedId !== id);
    } else {
      newSelectedIds = [...selectedIds, id];
    }

    setSelectedIds(newSelectedIds);
    debouncedOnMultiSelectChange(newSelectedIds);
  };

  const filterValueLabels =
    filterName !== 'Primary email status'
      ? _options
          .filter((option) => selectedIds?.includes(option.id))
          .map((option) => option.label)
      : groupOptions
          ?.flatMap((group) => group.options)
          .filter((option) => selectedIds?.includes(option.id))
          .map((option) => option.label);

  useEffect(() => {
    setTimeout(() => {
      if (isOpen && inputRef.current) {
        inputRef.current.focus();
      }
    }, 100);
  }, [isOpen]);

  const options = useMemo(
    () => [
      ..._options.filter(
        (o) =>
          o.label !== undefined &&
          o.label !== '' &&
          o.isArchived !== FlowStatus.Archived,
      ),
    ],
    [_options.length],
  );

  const Option = useCallback(
    ({ children, ...props }: OptionProps) => {
      const data = props?.data as {
        id: string;
        label: string;
        avatar?: string;
      };

      return (
        <components.Option {...props}>
          <div className='flex items-center gap-2'>
            {filterName === 'Owner' && (
              <Avatar
                size='xxs'
                textSize='xxs'
                variant='outlineCircle'
                src={data?.avatar ?? ''}
                name={data?.label ?? ''}
                icon={<User01 className='text-gray-500 size-3' />}
              />
            )}
            <span
              className='flex-1'
              style={{
                marginLeft: filterName === 'Primary email status' ? '8px' : '0',
              }}
            >
              {children}
            </span>
            {selectedIds.includes(data?.id) && (
              <Check className='text-primary-600' />
            )}
          </div>
        </components.Option>
      );
    },
    [selectedIds.length],
  );

  const CountryOption = useCallback(
    ({ children, ...props }: OptionProps) => {
      const country = props?.data as {
        id: string;
        label: string;
        avatar?: string;
      };

      return (
        <components.Option {...props}>
          <div className='flex items-center gap-2'>
            <span className='mb-[2px]'>{flags[country.id]}</span>
            <span>{children}</span>
            {selectedIds.includes(country.id) && (
              <Check className='text-primary-600' />
            )}
          </div>
        </components.Option>
      );
    },
    [selectedIds.length],
  );

  const getOptions = useCallback(() => {
    return filterName !== 'Primary email status' ? options : groupOptions;
  }, [filterName]);

  if (
    operatorName === ComparisonOperator.IsEmpty ||
    operatorName === ComparisonOperator.IsNotEmpty
  )
    return null;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const Options: ComponentType<OptionProps<any, any, any>> =
    filterName === 'Country' ? CountryOption : Option;

  const formatGroupLabel = (groupOption: GroupedOption) => (
    <div className='flex justify-between items-center'>
      <span className='font-medium text-gray-700'>{groupOption.label}</span>
    </div>
  );

  return (
    <Popover open={isOpen} onOpenChange={(open) => setIsOpen(open)}>
      <PopoverTrigger asChild>
        <Button
          size='xs'
          colorScheme='grayModern'
          className='border-l-0 rounded-none text-gray-700 bg-white font-normal'
        >
          {filterValueLabels?.length === 1
            ? filterValueLabels?.[0]
            : (filterValueLabels?.length ?? 0) > 1
            ? `${filterValueLabels?.length} ${handlePropertyPlural(
                filterName,
                selectedIds,
              )}`
            : '...'}
        </Button>
      </PopoverTrigger>
      <PopoverContent
        side='bottom'
        align='start'
        className='py-1 min-w-[254px]'
      >
        <Combobox
          escapeClearsValue
          options={getOptions()}
          closeMenuOnSelect={false}
          formatGroupLabel={formatGroupLabel}
          noOptionsMessage={() => 'Nothing in sight...'}
          onChange={(value) => handleItemClick(value.id)}
          components={{
            Option: Options,
          }}
          onKeyDown={(e) => {
            if (e.key === 'Escape') setIsOpen(false);
          }}
          placeholder={`${filterName} ${handleOperatorName(
            operatorName as ComparisonOperator,
            'list',
            selectedIds.length > 1,
          )}...`}
        />
      </PopoverContent>
    </Popover>
  );
};
