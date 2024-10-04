import { useState, useEffect, useCallback } from 'react';

import { debounce } from 'lodash';
import { format } from 'date-fns';

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
  filterValue: string | [string | null, string | null];
  onChangeFilterValue: (value: string | [string | null, string | null]) => void;
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
    debounce((value: string | [string | null, string | null]) => {
      onChangeFilterValue(value);
    }, 300),
    [onChangeFilterValue],
  );

  const handleDateChange = (date: Date | Date[] | null) => {
    if (!date) return;

    if (Array.isArray(date)) {
      const [start, end] = date;
      const formattedStartDate = start ? format(start, 'yyyy-MM-dd') : null;
      const formattedEndDate = end ? format(end, 'yyyy-MM-dd') : null;

      debouncedOnChangeFilterValue([formattedStartDate, formattedEndDate]);
    } else {
      const formattedDate = format(date, 'yyyy-MM-dd');

      if (operatorValue === ComparisonOperator.Lt) {
        debouncedOnChangeFilterValue([null, formattedDate]);
      } else if (operatorValue === ComparisonOperator.Gt) {
        debouncedOnChangeFilterValue([formattedDate, null]);
      } else {
        debouncedOnChangeFilterValue(formattedDate);
      }
    }
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

    return ComparisonOperator.Lt === operatorValue
      ? formatDate(filterValue?.[1])
      : ComparisonOperator.Gt === operatorValue
      ? formatDate(filterValue?.[0])
      : ComparisonOperator.Between === operatorValue
      ? `${formatDate(filterValue?.[0])} - ${formatDate(filterValue?.[1])}`
      : '';
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
          selectRange={ComparisonOperator.Between === operatorValue}
          onChange={(value) => handleDateChange(value as Date | Date[] | null)}
          value={
            ComparisonOperator.Lt === operatorValue
              ? filterValue?.[1]
              : ComparisonOperator.Gt === operatorValue
              ? filterValue?.[0]
              : [
                  new Date(filterValue?.[0] || Date.now()),
                  new Date(filterValue?.[1] || Date.now()),
                ]
          }
        />
      </PopoverContent>
    </Popover>
  );
};
