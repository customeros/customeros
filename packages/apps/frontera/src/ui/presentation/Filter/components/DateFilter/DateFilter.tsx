import { useState, useEffect } from 'react';

import { format } from 'date-fns';

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
  const [startDate, setStartDate] = useState<Date | null>(null);
  const [endDate, setEndDate] = useState<Date | null>(null);

  useEffect(() => {
    if (!filterValue) {
      if (filterName) {
        setTimeout(() => {
          setIsOpen(true);
        }, 100);
      }
    }
  }, [filterName]);

  const handleStartDateChange = (date: Date | null) => {
    if (!date) return;
    const formattedDate = format(date, 'yyyy-MM-dd');

    setStartDate(date);
    onChangeFilterValue([formattedDate, filterValue ? filterValue[1] : null]);
  };

  const handleEndDateChange = (date: Date | null) => {
    if (!date) return;
    const formattedDate = format(date, 'yyyy-MM-dd');

    setEndDate(date);
    onChangeFilterValue([filterValue ? filterValue[0] : null, formattedDate]);
    setIsOpen(false);
  };

  const handleSingleDateChange = (date: Date | null) => {
    if (!date) return;

    const formattedDate = format(date, 'yyyy-MM-dd');

    if (operatorValue === ComparisonOperator.Lt) {
      onChangeFilterValue([null, formattedDate]);
    } else if (operatorValue === ComparisonOperator.Gt) {
      onChangeFilterValue([formattedDate, null]);
    } else {
      onChangeFilterValue(formattedDate);
    }
    setIsOpen(false);
  };

  const selectedValue =
    filterValue &&
    (ComparisonOperator.Lt === operatorValue
      ? filterValue[1]
      : ComparisonOperator.Gt === operatorValue
      ? filterValue[0]
      : `${filterValue[0]} - ${filterValue[1]}`);

  return (
    <Popover
      modal={true}
      open={isOpen}
      onOpenChange={(value) => setIsOpen(value)}
    >
      <PopoverTrigger>
        <Button
          size='xs'
          colorScheme='grayModern'
          onClick={() => setIsOpen(!isOpen)}
          className='border-l-0 rounded-none text-gray-700 bg-white font-normal'
        >
          <span className=' max-w-[160px] text-ellipsis whitespace-nowrap overflow-hidden'>
            {selectedValue || '...'}
          </span>
        </Button>
      </PopoverTrigger>
      <PopoverContent
        side='bottom'
        align='start'
        className='py-1 min-w-[254px]'
      >
        {operatorValue === ComparisonOperator.Between ? (
          <div className='flex  space-y-2'>
            {/* Start Date Picker */}
            <div>
              <span>Start Date</span>
              <DatePicker
                value={startDate || new Date(Date.now())}
                onChange={(date) =>
                  handleStartDateChange(date instanceof Date ? date : null)
                }
              />
            </div>

            {/* End Date Picker */}
            <div>
              <span>End Date</span>
              <DatePicker
                value={endDate ? new Date(endDate) : new Date(Date.now())}
                onChange={(date) =>
                  handleEndDateChange(date instanceof Date ? date : null)
                }
              />
            </div>
          </div>
        ) : (
          <DatePicker
            defaultValue={new Date(Date.now())}
            onChange={(date) => {
              if (date instanceof Date) {
                handleSingleDateChange(date);
              }
            }}
            value={
              filterValue !== undefined
                ? new Date(
                    filterValue?.[0] === null
                      ? (filterValue?.[1] as string)
                      : (filterValue?.[0] as string),
                  )
                : new Date(Date.now())
            }
          />
        )}
      </PopoverContent>
    </Popover>
  );
};
