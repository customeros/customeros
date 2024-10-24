import { useRef, useState } from 'react';

import { set } from 'date-fns/set';
import { addDays } from 'date-fns/addDays';
import { getHours } from 'date-fns/getHours';
import { getMinutes } from 'date-fns/getMinutes';
import { toZonedTime, fromZonedTime } from 'date-fns-tz';

import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { DatePicker } from '@ui/form/DatePicker';
import { Input, InputProps } from '@ui/form/Input/Input';
import { Divider } from '@ui/presentation/Divider/Divider';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

// localDate and utcDate namings are misused.
// It should be the other way around.
// TODO: swap names of localDate and utcDate;

interface DueDatePickerProps {
  value: string;
  onChange: (value: string) => void;
}

export const ReminderDueDatePicker = ({
  value,
  onChange,
}: DueDatePickerProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;

  const localDate = fromZonedTime(new Date(value), timeZone);
  const time = (() => {
    const hours = String(getHours(localDate)).padStart(2, '0');
    const minutes = String(getMinutes(localDate)).padStart(2, '0');

    return `${hours}:${minutes}`;
  })();

  const handleChange = (date: Date | null) => {
    if (!date) return;
    const [hours, minutes] = time.split(':').map(Number);
    const updatedLocalDate = set(date, {
      hours,
      minutes,
      seconds: 0,
      milliseconds: 0,
    });
    const utcDate = toZonedTime(updatedLocalDate, timeZone);

    onChange(utcDate.toISOString());
  };

  const handleClickTomorrow = () => {
    const date = set(addDays(new Date(), 1), {
      hours: 9,
      minutes: 0,
      seconds: 0,
      milliseconds: 0,
    });
    const utcDate = toZonedTime(date, timeZone);

    onChange(utcDate.toISOString());
  };

  return (
    <div ref={containerRef} className='flex flex-start items-center'>
      <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
        <PopoverTrigger className='data-[state=open]:text-gray-700 data-[state=closed]:text-gray-500'>
          <span className='cursor-pointer whitespace-pre pb-[1px] text-sm border-t-[1px] border-transparent hover:text-gray-700'>
            {`${DateTimeUtils.format(value, DateTimeUtils.date)} • `}
          </span>
        </PopoverTrigger>
        <PopoverContent
          side='top'
          align='start'
          sticky='always'
          className='items-end'
          onClick={(e) => e.stopPropagation()}
          onOpenAutoFocus={(el) => el.preventDefault()}
        >
          <DatePicker
            value={localDate}
            minDate={new Date()}
            onChange={(date) => {
              handleChange(date as Date);
              setIsOpen(false);
            }}
          />
          <Divider className='my-2' />
          <Button
            variant='outline'
            className='rounded-full mr-3'
            onClick={() => {
              handleClickTomorrow();
              setIsOpen(false);
            }}
          >
            Tomorrow
          </Button>
        </PopoverContent>
      </Popover>
      <TimeInput
        value={time}
        onChange={(v) => {
          const [hours, minutes] = v.split(':').map(Number);
          const date = set(localDate, { hours, minutes });
          const utcDate = toZonedTime(date, timeZone);

          onChange(utcDate.toISOString());
        }}
      />
    </div>
  );
};

interface TimeInputProps extends Omit<InputProps, 'value' | 'onChange'> {
  value?: string;
  onChange?: (value: string) => void;
}

const TimeInput = ({ onChange, value, ...rest }: TimeInputProps) => {
  return (
    <Input
      size='xs'
      type='time'
      list='hidden'
      value={value}
      onChange={(e) => {
        const val = e.target.value;

        onChange?.(val);
      }}
      className='text-gray-500 mb-[-4px] text-sm appearance-none leading-[1] [&::-webkit-calendar-picker-indicator]:hidden p-0 min-h-0 w-fit focus:text-gray-700 focus:border-primary-500 cursor-text list-none'
      {...rest}
    />
  );
};
