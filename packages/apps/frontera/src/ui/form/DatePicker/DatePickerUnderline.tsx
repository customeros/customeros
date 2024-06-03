import { useField } from 'react-inverted-form';
import React, { useRef, useState } from 'react';

import { DateTimeUtils } from '@utils/date';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { ChevronLeft } from '@ui/media/icons/ChevronLeft.tsx';
import { ChevronRight } from '@ui/media/icons/ChevronRight.tsx';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

interface DatePickerProps {
  name: string;
  formId: string;
}
export const DatePickerUnderline: React.FC<DatePickerProps> = ({
  name,
  formId,
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);

  const { getInputProps } = useField(name, formId);
  const { onChange, value } = getInputProps();
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
    <div className='inline-flex flex-start items-center' ref={containerRef}>
      <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
        <PopoverTrigger className='data-[state=open]:text-gray-700 data-[state=closed]:text-gray-500'>
          <span className='underline cursor-pointer whitespace-pre pb-[1px] text-inherit border-t-[1px] border-transparent hover:text-gray-700'>{`${
            value
              ? DateTimeUtils.format(value, DateTimeUtils.date)
              : 'Select date'
          }`}</span>
        </PopoverTrigger>
        <PopoverContent
          align='start'
          side='top'
          className='items-end z-[999]'
          sticky='always'
          onOpenAutoFocus={(el) => el.preventDefault()}
          onClick={(e) => e.stopPropagation()}
        >
          <DatePicker
            name={name}
            formId={formId}
            defaultValue={new Date(value)}
            nextLabel={<ChevronRight />}
            prevLabel={<ChevronLeft />}
            onChange={(date) => {
              handleDateInputChange(date as Date);
            }}
          />
        </PopoverContent>
      </Popover>
    </div>
  );
};
