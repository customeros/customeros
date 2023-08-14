import React from 'react';
import { FormControl, FormLabel } from '@chakra-ui/react';
import { Flex } from '@ui/layout/Flex';
import { Box } from '@ui/layout/Box';
import Calendar from '@spaces/atoms/icons/Calendar';
import { Text } from '@ui/typography/Text';
import {
  DatePicker as ReactDatePicker,
  DatePickerProps as ReactDatePickerProps,
} from 'react-date-picker';
import { DateTimeUtils } from '@spaces/utils/date';
import Delete from '@spaces/atoms/icons/Delete';
import { useField } from 'react-inverted-form';

interface DatePickerProps extends ReactDatePickerProps {
  label: string;
  name: string;
  formId: string;
}

type DateInputValue = null | string | number | Date;

export const DatePicker: React.FC<DatePickerProps> = ({
  label,
  name,
  formId,
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
      <Flex>
        <ReactDatePicker
          id={id}
          clearIcon={
            value && (
              <Delete color='var(--chakra-colors-gray-500)' height='1rem' />
            )
          }
          onChange={(val) => handleDateInputChange(val as DateInputValue)}
          defaultValue={value}
          calendarIcon={
            <Flex alignItems='center'>
              <Box mr={3} color='gray.500'>
                <Calendar height={16} />
              </Box>
              <Text color={value ? 'gray.700' : 'gray.400'} role='button'>
                {value
                  ? DateTimeUtils.format(
                      value.toISOString(),
                      DateTimeUtils.dateWithFullMonth,
                    )
                  : 'Start date'}
              </Text>
            </Flex>
          }
        />
      </Flex>
    </FormControl>
  );
};
