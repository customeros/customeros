import React, { useRef } from 'react';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { DateTimeUtils } from '@spaces/utils/date';
import { Box } from '@ui/layout/Box';
import { DatePicker as ReactDatePicker } from 'react-date-picker';
import { useUpdateLogEntryMutation } from '@organization/graphql/updateLogEntry.generated';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { FormInput } from '@ui/form/Input';
import {
  LogEntryUpdateFormDto,
  LogEntryUpdateFormDtoI,
} from './LogEntryUpdateFormDto';
import { useField, useForm } from 'react-inverted-form';
import { useQueryClient } from '@tanstack/react-query';

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
export const LogEntryDatePicker: React.FC<{ event: LogEntryWithAliases }> = ({
  event,
}) => {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const logEntryStartedAtValues = new LogEntryUpdateFormDto(event);
  const { state } = useForm<LogEntryUpdateFormDtoI>({
    formId: 'log-entry-update',
    defaultValues: logEntryStartedAtValues,

    stateReducer: (_, _a, next) => {
      return next;
    },
  });
  const { getInputProps } = useField('date', 'log-entry-update');
  const { id, onChange, value: dateValue } = getInputProps();
  const updateLogEntryMutation = useUpdateLogEntryMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => queryClient.invalidateQueries(['GetTimeline.infinite']),
        500,
      );
    },
  });

  const handleUpdateStartedAt = () => {
    updateLogEntryMutation.mutate({
      id: event.id,
      input: {
        content: event.content,
        contentType: event.contentType,
        ...LogEntryUpdateFormDto.toPayload({
          ...state.values,
          time: state.values.time,
        }),
      },
    });
  };

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
          onBlur={handleUpdateStartedAt}
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
        <Flex alignItems='center' mx={1}>
          •
        </Flex>
        <Box>
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
            _hover={{ borderColor: 'transparent', cursor: 'text' }}
            _focus={{ borderColor: 'transparent', cursor: 'text' }}
            _focusVisible={{ borderColor: 'transparent', cursor: 'text' }}
            formId='log-entry-update'
            name='time'
            type='time'
            list='hidden'
            onBlur={handleUpdateStartedAt}
            defaultValue={DateTimeUtils.formatTime(event.logEntryStartedAt)}
          />
        </Box>
      </Flex>
    </>
  );
};
