import { useState, useEffect, useCallback } from 'react';

import { debounce } from 'lodash';

import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { DatePicker } from '@ui/form/DatePicker';
import { ComparisonOperator } from '@shared/types/__generated__/graphql.types';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

interface DateValueFilterProps {
  filterName: string;
  operatorValue: string;
  filterValue: string | null;
  onChangeFilterValue: (value: string | null) => void;
}

export const DateFilter = ({
  operatorValue,
  filterValue,
  onChangeFilterValue,
  filterName,
}: DateValueFilterProps) => {
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    if (!filterValue) {
      if (filterName) {
        setTimeout(() => {
          setIsOpen(true);
        }, 100);
      }
    }
  }, [filterName]);

  const debouncedOnChangeFilterValue = useCallback(
    debounce((value: string | null) => {
      onChangeFilterValue(value);
    }, 300),
    [onChangeFilterValue],
  );

  const handleDateChange = (date: Date | null) => {
    if (!date) return;

    const formattedDate = DateTimeUtils.getUTCDateAtMidnight(date);

    debouncedOnChangeFilterValue(formattedDate);
    setIsOpen(false);
  };

  const selectedValue = () => {
    const currentYear = new Date().getFullYear();

    const formatDate = (date: string | null) => {
      if (!date) return '...';
      const dateObj = new Date(date);
      const year = dateObj.getFullYear();

      return DateTimeUtils.format(
        dateObj.toString(),
        year === currentYear
          ? DateTimeUtils.dateDayAndMonth
          : DateTimeUtils.date,
      );
    };

    return ComparisonOperator.Lt === operatorValue ||
      ComparisonOperator.Gt === operatorValue
      ? formatDate(filterValue)
      : '...';
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
          className='border-l-0 rounded-none text-gray-700 bg-white font-normal'
        >
          <span className=' max-w-[160px] text-ellipsis whitespace-nowrap overflow-hidden'>
            {selectedValue()}
          </span>
        </Button>
      </PopoverTrigger>
      <PopoverContent
        side='bottom'
        align='start'
        className='py-1 min-w-[254px]'
      >
        <DatePicker
          value={new Date(filterValue || new Date())}
          onChange={(value) => handleDateChange(value as Date | null)}
        />
      </PopoverContent>
    </Popover>
  );
};
