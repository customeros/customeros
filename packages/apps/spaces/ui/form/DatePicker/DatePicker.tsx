import React from 'react';
import { useField } from 'react-inverted-form';
import {
  DatePicker as ReactDatePicker,
  DatePickerProps as ReactDatePickerProps,
} from 'react-date-picker';

import { FormLabel, FormControl } from '@chakra-ui/react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Delete } from '@ui/media/icons/Delete';
import { DateTimeUtils } from '@spaces/utils/date';
import { Calendar } from '@ui/media/icons/Calendar';

interface DatePickerProps extends ReactDatePickerProps {
  name: string;
  label: string;
  inset?: string;
  formId: string;
  placeholder?: string;
  calendarIconHidden?: boolean;
}

type DateInputValue = null | string | number | Date;

export const DatePicker: React.FC<DatePickerProps> = ({
  label,
  name,
  formId,
  placeholder,
  calendarIconHidden,
  inset,
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
    <FormControl>
      <FormLabel fontWeight={600} color='gray.700' fontSize='sm' mb={-1}>
        {label}
      </FormLabel>
      <Flex
        sx={{
          '& .react-date-picker__calendar-button': {
            pl: 0,
          },
          '& .react-date-picker__calendar': {
            inset: `${inset ?? '120% 0px auto auto'} !important`,
          },
          '& .react-date-picker__clear-button': {
            top: '7px',
          },
          '& .react-calendar__month-view__weekdays__weekday': {
            textTransform: 'capitalize',
          },
        }}
      >
        <ReactDatePicker
          id={id}
          clearIcon={value && <Delete color='gray.500' height='1rem' />}
          onChange={(val) => handleDateInputChange(val as DateInputValue)}
          defaultValue={value}
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
              {!calendarIconHidden && (
                <Box mr={3} color='gray.500'>
                  <Calendar />
                </Box>
              )}

              <Text color={value ? 'gray.700' : 'gray.400'} role='button'>
                {value
                  ? DateTimeUtils.format(
                      value.toISOString(),
                      DateTimeUtils.dateWithAbreviatedMonth,
                    )
                  : `${placeholder ? placeholder : 'Start date'}`}
              </Text>
            </Flex>
          }
        />
      </Flex>
    </FormControl>
  );
};
