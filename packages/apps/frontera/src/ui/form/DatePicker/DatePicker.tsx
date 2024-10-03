import { forwardRef } from 'react';
import Calendar, { CalendarProps } from 'react-calendar';

import { ChevronLeft } from '@ui/media/icons/ChevronLeft';
import { ChevronRight } from '@ui/media/icons/ChevronRight';

export const DatePicker = forwardRef(
  ({ value, onChange, ...props }: CalendarProps, ref) => {
    const handleDateInputChange = (
      value: CalendarProps['value'],
      event: React.MouseEvent<HTMLButtonElement, MouseEvent>,
    ) => {
      if (!value) return onChange?.(null, event);

      if (Array.isArray(value)) {
        const [startDate, endDate] = value;
        const normalizedStartDate = new Date(
          Date.UTC(
            (startDate instanceof Date
              ? startDate
              : new Date(startDate || new Date())
            ).getFullYear(),
            (startDate instanceof Date
              ? startDate
              : new Date(startDate || new Date())
            ).getMonth(),
            (startDate instanceof Date
              ? startDate
              : new Date(startDate || new Date())
            ).getDate(),
          ),
        );
        const normalizedEndDate = new Date(
          Date.UTC(
            (endDate instanceof Date
              ? endDate
              : new Date(endDate || new Date())
            ).getFullYear(),
            (endDate instanceof Date
              ? endDate
              : new Date(endDate || new Date())
            ).getMonth(),
            (endDate instanceof Date
              ? endDate
              : new Date(endDate || new Date())
            ).getDate(),
          ),
        );

        onChange?.([normalizedStartDate, normalizedEndDate], event);
      } else {
        const date = new Date(value as string | number | Date);
        const normalizedDate = new Date(
          Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()),
        );

        onChange?.(normalizedDate, event);
      }
    };

    return (
      <Calendar
        ref={ref}
        value={value}
        defaultValue={value}
        prevLabel={<ChevronLeft />}
        nextLabel={<ChevronRight />}
        onChange={handleDateInputChange}
        {...props}
      />
    );
  },
);
