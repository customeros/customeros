import React from 'react';
import { useField } from 'react-inverted-form';
import {
  DatePicker as ReactDatePicker,
  DatePickerProps as ReactDatePickerProps,
} from 'react-date-picker';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
type DateInputValue = null | string | number | Date;

interface DatePickerProps extends ReactDatePickerProps {
  name: string;
  inset?: string;
  formId: string;
  placeholder?: string;
  calendarIconHidden?: boolean;
}

export const DatePickerUnderline: React.FC<DatePickerProps> = ({
  name,
  formId,
  placeholder,
  value,
}) => {
  const { getInputProps } = useField(name, formId);
  const { id, onChange } = getInputProps();
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
    <Flex
      sx={{
        '& .react-date-picker': {
          height: 'unset',
        },
        '& .react-date-picker__button': {
          padding: 0,
        },
        '& .react-date-picker__calendar-button': {
          width: 'fit-content',
          borderBottom: 'none !important',
          '&:after': {
            content: "''",

            height: '1px',
            width: '100%',
            background: 'gray.700',
            borderBottom: 'none !important',
            display: 'block',
            mt: '-2px',
          },
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
        defaultValue={new Date()}
        formatShortWeekday={(_, date) =>
          DateTimeUtils.format(date.toISOString(), DateTimeUtils.shortWeekday)
        }
        formatMonth={(_, date) =>
          DateTimeUtils.format(
            date.toISOString(),
            DateTimeUtils.abreviatedMonth,
          )
        }
        calendarIcon={
          <Flex alignItems='center'>
            <Text color={value ? 'gray.700' : 'gray.400'} role='button'>
              {value
                ? DateTimeUtils.format(
                    (value as Date)?.toISOString(),
                    DateTimeUtils.dateWithAbreviatedMonth,
                  )
                : `${placeholder ? placeholder : 'Start date'}`}
            </Text>
          </Flex>
        }
      />
    </Flex>
  );
};
