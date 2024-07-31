import React, { useState } from 'react';
import { useField } from 'react-inverted-form';

import { DateTimeUtils } from '@utils/date';
import { FormDatePicker } from '@ui/form/DatePicker';
import { FormInput, FormInputProps } from '@ui/form/Input/FormInput';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover';

export const LogEntryDatePicker: React.FC<{
  formId: string;
  event: LogEntryWithAliases;
}> = ({ event, formId }) => {
  const { getInputProps } = useField('date', formId);
  const { id, onChange, value: dateValue, onBlur } = getInputProps();
  const { getInputProps: getTimeInputProps } = useField('time', formId);
  const { onBlur: onTimeBlur, value: timeValue } = getTimeInputProps();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <label htmlFor={id} className='text-sm font-semibold text-gray-700'>
        Date
      </label>
      <div className='flex'>
        <Popover open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
          <PopoverTrigger className='data-[state=open]:text-gray-700 data-[state=closed]:text-gray-500'>
            <span className=' cursor-pointer whitespace-pre pb-[1px] text-sm border-t-[1px] border-transparent hover:text-gray-700'>{`${DateTimeUtils.format(
              dateValue,
              'EEEE, dd MMM yyyy',
            )} • `}</span>
          </PopoverTrigger>
          <PopoverContent
            side='top'
            align='start'
            sticky='always'
            className='flex p-0 z-50'
          >
            <FormDatePicker
              name='date'
              label='Date'
              formId={formId}
              value={dateValue}
              onBlur={() => onBlur(dateValue)}
              labelProps={{ className: 'hidden' }}
              onChange={(value) => {
                onChange(value as Date);
                setIsOpen(false);
              }}
            />
          </PopoverContent>
        </Popover>
        <TimeInput
          name='time'
          formId='log-entry-update'
          onBlur={() => onTimeBlur(timeValue)}
          defaultValue={DateTimeUtils.formatTime(event.logEntryStartedAt)}
        />
      </div>
    </>
  );
};

interface TimeInputProps extends Omit<FormInputProps, 'value' | 'onChange'> {
  value?: string;
  onChange?: (value: string) => void;
}

const TimeInput = ({ onChange, value, ...rest }: TimeInputProps) => {
  return (
    <FormInput
      size='xs'
      type='time'
      list='hidden'
      value={value}
      className='text-gray-500 mb-[-3px] text-sm appearance-none leading-[1] [&::-webkit-calendar-picker-indicator]:hidden p-0 min-h-0 w-fit focus:text-gray-700 focus:border-primary-500 cursor-text list-none'
      {...rest}
    />
  );
};
