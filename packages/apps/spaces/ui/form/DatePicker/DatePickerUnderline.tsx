import React from 'react';
import { useField } from 'react-inverted-form';
import {
  DatePicker as ReactDatePicker,
  DatePickerProps as ReactDatePickerProps,
} from 'react-date-picker';

import { Box } from '@chakra-ui/layout';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@spaces/utils/date';
type DateInputValue = null | string | number | Date;

interface DatePickerProps extends ReactDatePickerProps {
  name: string;
  inset?: string;
  formId: string;
  placeholder?: string;
  defaultOpen?: boolean;
  calendarIconHidden?: boolean;
}

export const DatePickerUnderline: React.FC<DatePickerProps> = ({
  name,
  formId,
  placeholder,
  defaultOpen,
}) => {
  const { getInputProps } = useField(name, formId);
  const { id, onChange, value } = getInputProps();
  const handleDateInputChange = (data?: DateInputValue) => {
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
  };

  return (
    <Box
      display='inline'
      sx={{
        '& .react-date-picker__wrapper': {
          display: 'inline-flex',
        },

        '& .react-date-picker__calendar': {
          inset: `20px auto auto auto !important`,
        },
        '& .react-date-picker': {
          height: 'unset',
          display: 'inline',
        },
        '& .react-date-picker__button': {
          padding: 0,
        },
        '& .react-date-picker__calendar-button': {
          width: 'fit-content',
          borderBottom: 'none !important',
        },
        '& .react-date-picker--open .react-date-picker__calendar-button, .react-date-picker:focus-within .react-date-picker__calendar-button, .react-date-picker:focus .react-date-picker__calendar-button, .react-date-picker:focus-visible .react-date-picker__calendar-button':
          {
            borderBottom: 'none !important',
          },
        '& .react-date-picker:hover .react-date-picker__calendar-button': {
          borderBottom: 'none !important',
        },
      }}
    >
      <ReactDatePicker
        id={id}
        name={name}
        clearIcon={() => null}
        onChange={(val) => handleDateInputChange(val as DateInputValue)}
        formatShortWeekday={(_, date) =>
          DateTimeUtils.format(date.toISOString(), DateTimeUtils.shortWeekday)
        }
        formatMonth={(_, date) =>
          DateTimeUtils.format(
            date.toISOString(),
            DateTimeUtils.abreviatedMonth,
          )
        }
        value={value as Date}
        calendarIcon={
          <p
            className={cn(
              'underline text-gray-500 hover:text-gray-700 focus:text-gray-700',
            )}
            role='button'
          >
            {value
              ? DateTimeUtils.format(
                  (value as Date)?.toISOString(),
                  DateTimeUtils.dateWithAbreviatedMonth,
                )
              : `${placeholder ? placeholder : 'Start date'}`}
          </p>
        }
      />
    </Box>
  );
};
