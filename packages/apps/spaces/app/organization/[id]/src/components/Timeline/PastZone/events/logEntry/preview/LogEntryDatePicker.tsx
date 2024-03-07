import React from 'react';
import { useField } from 'react-inverted-form';
import { DatePicker as ReactDatePicker } from 'react-date-picker';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FormInput } from '@ui/form/Input';
import { DateTimeUtils } from '@spaces/utils/date';
import { LogEntryWithAliases } from '@organization/src/components/Timeline/types';

const calendarStyles = {
  '& .react-date-picker--open .react-date-picker__calendar-button, .react-date-picker:focus-within .react-date-picker__calendar-button, .react-date-picker:focus .react-date-picker__calendar-button, .react-date-picker:focus-visible .react-date-picker__calendar-button':
    {
      borderColor: 'transparent !important',
    },
  '& .react-date-picker__calendar-button:hover': {
    borderColor: 'transparent !important',
  },

  '& .react-date-picker': {
    height: 'min-content',
    position: 'initial !important',
  },
  '& .react-date-picker__wrapper': {
    height: 'min-content',
    position: 'initial',
  },
  '& .react-date-picker__button': {
    p: 0,
    height: 'min-content',
  },
};
export const LogEntryDatePicker: React.FC<{
  formId: string;
  event: LogEntryWithAliases;
}> = ({ event, formId }) => {
  const { getInputProps } = useField('date', formId);
  const { id, onChange, value: dateValue, onBlur } = getInputProps();
  const { getInputProps: getTimeInputProps } = useField('time', formId);
  const { onBlur: onTimeBlur, value: timeValue } = getTimeInputProps();

  return (
    <>
      <Text
        size='sm'
        fontSize='sm'
        fontWeight='semibold'
        as='label'
        htmlFor={id}
      >
        Date
      </Text>
      <Flex
        alignItems='center'
        sx={{
          ...calendarStyles,
          '& .react-date-picker__calendar': {
            inset: `95px auto auto auto !important`,
          },
        }}
      >
        <ReactDatePicker
          id={id}
          clearIcon={null}
          onChange={onChange}
          onBlur={() => onBlur(dateValue)}
          defaultValue={new Date(event.logEntryStartedAt).toISOString()}
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
            <Text
              color={event.logEntryStartedAt ? 'gray.700' : 'gray.400'}
              role='button'
              fontSize='sm'
            >
              {DateTimeUtils.format(dateValue, 'EEEE, dd MMM yyyy')}
            </Text>
          }
        />
        <Text as='span' alignSelf='baseline' mx={1} fontSize='sm'>
          •
        </Text>
        <Box lineHeight={0}>
          <FormInput
            p={0}
            h='min-content'
            sx={{
              '&::-webkit-calendar-picker-indicator': {
                display: 'none',
              },
            }}
            fontSize='sm'
            size='xs'
            lineHeight='1'
            cursor='text'
            formId='log-entry-update'
            name='time'
            type='time'
            list='hidden'
            onBlur={() => onTimeBlur(timeValue)}
            defaultValue={DateTimeUtils.formatTime(event.logEntryStartedAt)}
          />
        </Box>
      </Flex>
    </>
  );
};
