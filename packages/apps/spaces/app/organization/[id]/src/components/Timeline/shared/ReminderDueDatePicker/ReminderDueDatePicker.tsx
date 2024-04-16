import { useRef, useState } from 'react';
import { useField } from 'react-inverted-form';

import set from 'date-fns/set';
import addDays from 'date-fns/addDays';
import getHours from 'date-fns/getHours';
import getMinutes from 'date-fns/getMinutes';

import { Portal } from '@ui/utils/';
import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { InlineDatePicker } from '@ui/form/DatePicker';
import { Input, InputProps } from '@ui/form/Input/Input2';
import { Divider } from '@ui/presentation/Divider/Divider';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

interface DueDatePickerProps {
  name: string;
  formId: string;
}

export const ReminderDueDatePicker = ({ name, formId }: DueDatePickerProps) => {
  const { getInputProps } = useField(name, formId);
  const { onChange, ...inputProps } = getInputProps();
  const containerRef = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);

  const time = (() => {
    const dateStr = inputProps.value;
    const date = dateStr ? new Date(dateStr) : new Date();

    const hours = (() => {
      const h = String(getHours(date));

      return h.length === 1 ? `0${h}` : h;
    })();
    const minutes = (() => {
      const h = String(getMinutes(date));

      return h.length === 1 ? `0${h}` : h;
    })();

    return `${hours}:${minutes}`;
  })();

  const handleChange = (date: Date | null) => {
    if (!date) return;
    const [hours, minutes] = time.split(':').map(Number);
    const _date = set(date, { hours, minutes, seconds: 0, milliseconds: 0 });

    onChange(_date.toISOString());
  };

  const handleClickTomorrow = () => {
    const date = set(addDays(new Date(), 1), {
      hours: 9,
      minutes: 0,
      seconds: 0,
      milliseconds: 0,
    });
    onChange(date.toISOString());
  };

  return (
    <div className='flex flex-start items-center' ref={containerRef}>
      <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
        <PopoverTrigger className='data-[state=open]:text-gray-700 data-[state=closed]:text-gray-500'>
          <span className=' cursor-pointer whitespace-pre pb-[1px] text-sm border-t-[1px] border-transparent hover:text-gray-700'>{`${DateTimeUtils.format(
            inputProps.value,
            DateTimeUtils.date,
          )} • `}</span>
        </PopoverTrigger>
        <Portal>
          <PopoverContent
            align='start'
            side='top'
            className='items-end'
            sticky='always'
          >
            <InlineDatePicker
              {...inputProps}
              onChange={(date) => {
                handleChange(date);
                setIsOpen(false);
              }}
              minDate={new Date()}
            />

            <Divider className='my-2' />
            <Button
              className='rounded-full mr-3'
              variant='outline'
              onClick={() => {
                handleClickTomorrow();
                setIsOpen(false);
              }}
            >
              Tomorrow
            </Button>
          </PopoverContent>
        </Portal>
      </Popover>
      <TimeInput
        value={time}
        onChange={(v) => {
          const [hours, minutes] = v.split(':').map(Number);
          const date = set(new Date(inputProps.value), { hours, minutes });

          onChange(date.toISOString());
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
      className='text-gray-500 mb-[-2px] text-sm appearance-none leading-[1] [&::-webkit-calendar-picker-indicator]:hidden p-0 min-h-0 w-fit focus:text-gray-700 focus:border-primary-500 cursor-text list-none'
      type='time'
      list='hidden'
      size='xs'
      value={value}
      onChange={(e) => {
        const val = e.target.value;
        onChange?.(val);
      }}
      {...rest}
    />
  );
};
