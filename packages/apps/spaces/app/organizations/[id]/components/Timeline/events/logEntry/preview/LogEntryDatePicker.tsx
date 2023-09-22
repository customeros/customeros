import React from 'react';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { DateTimeUtils } from '@spaces/utils/date';
import { Box } from '@ui/layout/Box';
import { DatePicker as ReactDatePicker } from 'react-date-picker';
import { FormInput } from '@ui/form/Input';
import { useField } from 'react-inverted-form';

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
  event: LogEntryWithAliases;
  formId: string;
}> = ({ event, formId }) => {
  const { getInputProps } = useField('date', formId);
  const { id, onChange, value: dateValue, onBlur } = getInputProps();

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
            inset: `116px auto auto auto !important`,
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
        <Text as='span' mx={1}>
          •
        </Text>
        <Box lineHeight={1}>
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
            onBlur={() => onBlur(dateValue)}
            defaultValue={DateTimeUtils.formatTime(event.logEntryStartedAt)}
          />
        </Box>
      </Flex>
    </>
  );
};
