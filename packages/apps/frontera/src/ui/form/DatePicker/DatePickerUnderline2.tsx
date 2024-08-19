import Calendar from 'react-calendar';
import React, { useRef, useState } from 'react';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft.tsx';
import { ChevronRight } from '@ui/media/icons/ChevronRight.tsx';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

interface DatePickerProps {
  minDate?: Date;
  classNames?: string;
  value?: Date | string;
  onChange: (date: Date | null) => void;
}

export const DatePickerUnderline2 = ({
  onChange,
  value,
  classNames,
  minDate,
}: DatePickerProps) => {
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
    <div ref={containerRef} className='flex flex-start items-center'>
      <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
        <PopoverTrigger className='data-[state=open]:text-gray-700 data-[state=closed]:text-gray-500'>
          <span
            className={cn(
              'underline cursor-pointer whitespace-pre pb-[1px] text-inherit border-t-[1px] border-transparent hover:text-gray-700 ',
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
          sticky='always'
          className='items-end z-[999]'
          onClick={(e) => e.stopPropagation()}
          onOpenAutoFocus={(el) => el.preventDefault()}
        >
          <div>
            <Calendar
              minDate={minDate}
              defaultValue={value}
              prevLabel={<ChevronLeft />}
              nextLabel={<ChevronRight />}
              onChange={(value) => handleDateInputChange(value as Date)}
            />
          </div>
        </PopoverContent>
      </Popover>
    </div>
  );
};
