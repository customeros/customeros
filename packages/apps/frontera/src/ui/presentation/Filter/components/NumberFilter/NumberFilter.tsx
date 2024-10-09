import { useState, useEffect } from 'react';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

import { handleOperatorName } from '../../utils/utils';

interface NumberFilterProps {
  filterName: string;
  filterValue: string;
  operatorValue: string;
  onChangeFilterValue: (value: string) => void;
}

const formatNumberWithCommas = (value: string | number | undefined): string => {
  if (value === undefined || value === '' || value === 0) return '';
  const numString = value?.toString();

  return numString?.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

export const NumberFilter = ({
  filterName,
  operatorValue,
  onChangeFilterValue,
  filterValue,
}: NumberFilterProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [inputValue, setInputValue] = useState(filterValue);
  const [timeoutId, setTimeoutId] = useState<NodeJS.Timeout | null>(null);

  useEffect(() => {
    if (!filterValue && filterName) {
      const timer = setTimeout(() => {
        setIsOpen(true);
      }, 100);

      return () => clearTimeout(timer);
    }
  }, [filterName, filterValue]);

  if (
    operatorValue === ComparisonOperator.IsEmpty ||
    operatorValue === ComparisonOperator.IsNotEmpty
  )
    return null;

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;

    setInputValue(newValue);

    if (timeoutId) clearTimeout(timeoutId);

    const newTimeoutId = setTimeout(() => {
      onChangeFilterValue(newValue);
    }, 300);

    setTimeoutId(newTimeoutId);
  };

  return (
    <Popover
      modal={true}
      open={isOpen}
      onOpenChange={(value) => setIsOpen(value)}
    >
      <PopoverTrigger asChild>
        <Button
          size='xs'
          colorScheme='grayModern'
          onClick={() => setIsOpen(!isOpen)}
          className='rounded-none text-gray-700 bg-white font-normal border-l-0'
        >
          <span className='max-w-[160px] text-ellipsis whitespace-nowrap overflow-hidden'>
            {(filterName !== 'Founded'
              ? formatNumberWithCommas(filterValue)
              : filterValue) || '...'}
          </span>
        </Button>
      </PopoverTrigger>
      <PopoverContent
        side='bottom'
        align='start'
        className='py-1 min-w-[254px]'
      >
        <Input
          size='sm'
          type='number'
          variant='unstyled'
          value={inputValue}
          onChange={handleInputChange}
          placeholder={`${filterName} ${handleOperatorName(
            operatorValue as ComparisonOperator,
          )}...`}
          onKeyDown={(e) => {
            if (e.key === 'e') {
              e.preventDefault();
            }

            if (e.key === 'Escape') {
              setIsOpen(false);
            }
          }}
        />
      </PopoverContent>
    </Popover>
  );
};
