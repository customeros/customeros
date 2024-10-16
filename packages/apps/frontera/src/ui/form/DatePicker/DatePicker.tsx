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

      const normalizedDate = new Date(value as string | number | Date);

      onChange?.(normalizedDate, event);
    };

    return (
      <Calendar
        ref={ref}
        value={value}
        defaultValue={value}
        prevLabel={<ChevronLeft />}
        nextLabel={<ChevronRight />}
        onChange={handleDateInputChange}
        formatMonth={(locale, date) =>
          date.toLocaleDateString(locale, { month: 'short' })
        }
        {...props}
      />
    );
  },
);
