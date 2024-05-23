import Calendar from 'react-calendar';
import React, { useRef, useState } from 'react';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@spaces/utils/date';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft.tsx';
import { ChevronRight } from '@ui/media/icons/ChevronRight.tsx';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

interface DatePickerProps {
  value?: Date;
  minDate?: Date;
  classNames?: string;
  onChange: (date: Date | null) => void;
}
export const DatePickerUnderline2: React.FC<DatePickerProps> = ({
  onChange,
  value,
  classNames,
  minDate,
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);

  const handleDateInputChange = (data?: Date) => {
    if (!data) return onChange(null);
    const date = new Date(data);

    const normalizedDate = new Date(
      Date.UTC(
        date.getFullYear(),
        date.getMonth(),
        date.getDate(),
        date.getHours(),
        date.getMinutes(),
        date.getSeconds(),
      ),
    );
    onChange(normalizedDate);
    setIsOpen(false);
  };

  return (
    <div className='flex flex-start items-center' ref={containerRef}>
      <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
        <PopoverTrigger className='data-[state=open]:text-gray-700 data-[state=closed]:text-gray-500'>
          <span
            className={cn(
              'underline cursor-pointer whitespace-pre pb-[1px] text-inherit border-t-[1px] border-transparent hover:text-gray-700',
              classNames,
            )}
          >
            {value
              ? `${DateTimeUtils.format(
                  `${value?.toString()}`,
                  DateTimeUtils.dateWithShortYear,
                )}`
              : 'Select date'}
          </span>
        </PopoverTrigger>
        <PopoverContent
          align='start'
          side='bottom'
          className='items-end z-[999]'
          sticky='always'
          onOpenAutoFocus={(el) => el.preventDefault()}
          onClick={(e) => e.stopPropagation()}
        >
          <div>
            <Calendar
              onChange={(value) => handleDateInputChange(value as Date)}
              defaultValue={value}
              nextLabel={<ChevronRight />}
              prevLabel={<ChevronLeft />}
              minDate={minDate}
            />
          </div>
        </PopoverContent>
      </Popover>
    </div>
  );
};
