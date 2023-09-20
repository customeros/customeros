import React, { useRef } from 'react';
import { CardHeader, CardBody } from '@ui/presentation/Card';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/presentation/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { useTimelineEventPreviewContext } from '../../preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import CopyLink from '@spaces/atoms/icons/CopyLink';
import Times from '@spaces/atoms/icons/Times';
import copy from 'copy-to-clipboard';
import { VStack } from '@ui/layout/Stack';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { User } from '@graphql/types';
import { DateTimeUtils } from '@spaces/utils/date';
import { Box } from '@ui/layout/Box';
import Image from 'next/image';
import noteIcon from 'public/images/event-ill-log-preview.png';
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

const getAuthor = (user: User) => {
  if (!user?.firstName && !user.lastName) {
    return 'Unknown';
  }

  return `${user.firstName} ${user.lastName}`.trim();
};

export const LogEntryPreviewModal: React.FC = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const event = modalContent as LogEntryWithAliases;
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
  const author = getAuthor(event?.logEntryCreatedBy);
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
        ...LogEntryUpdateFormDto.toPayload({
          ...state.values,
        }),
      },
    });
  };
  return (
    <>
      <CardHeader pb={1} position='sticky' top={0} borderRadius='xl'>
        <Flex
          direction='row'
          justifyContent='space-between'
          alignItems='center'
        >
          <Flex mb={2} alignItems='center'>
            <Heading size='sm' fontSize='lg'>
              Log entry
            </Heading>
          </Flex>
          <Flex direction='row' justifyContent='flex-end' alignItems='center'>
            <Tooltip label='Copy link to this entry' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Copy link to this entry'
                color='gray.500'
                size='sm'
                mr={1}
                icon={<CopyLink color='gray.500' height='18px' />}
                onClick={() => copy(window.location.href)}
              />
            </Tooltip>
            <Tooltip label='Close' aria-label='close' placement='bottom'>
              <IconButton
                variant='ghost'
                aria-label='Close preview'
                color='gray.500'
                size='sm'
                icon={<Times color='gray.500' height='24px' />}
                onClick={closeModal}
              />
            </Tooltip>
          </Flex>
        </Flex>
      </CardHeader>
      <CardBody mt={0} maxHeight='50%' pb={6}>
        <VStack gap={2} alignItems='flex-start' position='relative'>
          <Box position='absolute' top={-2} right={-3}>
            <Image src={noteIcon} alt='' height={123} width={174} />
          </Box>
          <Flex direction='column'>
            <Text size='sm' fontWeight='semibold'>
              Date
            </Text>
            <Flex
              sx={{
                '& .react-date-picker--open .react-date-picker__calendar-button, .react-date-picker:focus-within .react-date-picker__calendar-button, .react-date-picker:focus .react-date-picker__calendar-button, .react-date-picker:focus-visible .react-date-picker__calendar-button':
                  {
                    borderColor: 'transparent !important',
                  },
                '& .react-date-picker__calendar-button:hover': {
                  borderColor: 'transparent !important',
                },
                '& .react-date-picker__calendar': {
                  inset: '120% auto auto !important',
                },
                '& .react-date-picker': {
                  height: 'min-content',
                },
                '& .react-date-picker__wrapper': {
                  height: 'min-content',
                },
                '& .react-date-picker__button': {
                  p: 0,
                  height: 'min-content',
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
                  DateTimeUtils.format(
                    date.toISOString(),
                    DateTimeUtils.shortWeekday,
                  )
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
                    size='sm'
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
                  _hover={{ borderColor: 'transparent', cursor: 'text' }}
                  _focus={{ borderColor: 'transparent', cursor: 'text' }}
                  _focusVisible={{ borderColor: 'transparent', cursor: 'text' }}
                  formId='log-entry-update'
                  name='time'
                  type='time'
                  list='hidden'
                  onBlur={handleUpdateStartedAt}
                  defaultValue={DateTimeUtils.formatTime(
                    event.logEntryStartedAt,
                  )}
                />
              </Box>
            </Flex>
          </Flex>
          <Flex direction='column'>
            <Text size='sm' fontWeight='semibold'>
              Author
            </Text>
            <Text size='sm'>{author}</Text>
          </Flex>

          <Flex direction='column'>
            <Text size='sm' fontWeight='semibold'>
              Entry
            </Text>
            <Text
              className='slack-container'
              dangerouslySetInnerHTML={{ __html: `${event?.content}` }}
            />
          </Flex>

          <Text size='sm' fontWeight='medium'>
            {event.tags.map(({ name }) => `#${name}`).join(' ')}
          </Text>
        </VStack>
      </CardBody>
    </>
  );
};
