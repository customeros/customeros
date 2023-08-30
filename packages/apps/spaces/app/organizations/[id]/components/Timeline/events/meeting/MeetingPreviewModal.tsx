import { useState, useEffect } from 'react';
import { useForm } from 'react-inverted-form';
import { convert } from 'html-to-text';

import { Icons } from '@ui/media/Icon';
import { Flex } from '@ui/layout/Flex';
import { Meeting } from '@graphql/types';
import { VStack } from '@ui/layout/Stack';
import { Center } from '@ui/layout/Center';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { DateTimeUtils } from '@spaces/utils/date';
import { ScaleFade } from '@ui/transitions/ScaleFade';
import { Card, CardBody, CardHeader, CardFooter } from '@ui/presentation/Card';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';

import { useTimelineEventPreviewContext } from '../../preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { MeetingIcon } from './meetingIcon';

export const MeetingPreviewModal = () => {
  const [isMounted, setIsMounted] = useState(false);
  const { isModalOpen, closeModal, modalContent } =
    useTimelineEventPreviewContext();

  const event = modalContent as Meeting;

  const when = `${DateTimeUtils.format(
    event?.startedAt,
    DateTimeUtils.longWeekday,
  )}, ${DateTimeUtils.format(
    event?.startedAt,
    DateTimeUtils.dateWithAbreviatedMonth,
  )} • ${DateTimeUtils.format(
    event.startedAt,
    DateTimeUtils.usaTimeFormatString,
  )} - ${DateTimeUtils.format(
    event.endedAt,
    DateTimeUtils.usaTimeFormatString,
  )}`;

  const participants = event?.attendedBy
    ?.map(
      (c) =>
        c.__typename === 'ContactParticipant' &&
        (c.contactParticipant?.firstName ||
          c.contactParticipant?.emails?.[0]?.email),
    )
    .filter(Boolean);

  const owner = (() => {
    if (event?.createdBy?.[0]?.__typename === 'ContactParticipant') {
      const participant = event?.createdBy?.[0]?.contactParticipant;
      return participant?.firstName ?? participant.emails?.[0]?.email;
    }
    return '';
  })();

  useForm<{ note: string }>({
    formId: 'meeting-notes',
    defaultValues: {
      note: convert(event?.note?.[0]?.html ?? ''),
    },
  });

  useEffect(() => {
    setIsMounted(true);
  }, []);

  return (
    <Flex
      position='absolute'
      top='0'
      bottom='0'
      left='0'
      right='0'
      zIndex={1}
      onClick={closeModal}
      cursor='pointer'
      backdropFilter='blur(3px)'
      justify='center'
      background={isMounted ? 'rgba(16, 24, 40, 0.45)' : 'rgba(16, 24, 40, 0)'}
      align='center'
      transition='all 0.1s linear'
    >
      <ScaleFade
        in={isModalOpen}
        style={{
          position: 'absolute',
          marginInline: 'auto',
          top: '1rem',
          width: '544px',
          minWidth: '544px',
        }}
      >
        <Card
          size='lg'
          position='absolute'
          mx='auto'
          top='4'
          w='544px'
          minW='544px'
          cursor='default'
          onClick={(e) => e.stopPropagation()}
        >
          <CardHeader
            as={Flex}
            justify='space-between'
            gap='4'
            align='center'
            px='6'
            pt='4'
            pb='2'
          >
            <Text
              fontSize='lg'
              color='gray.700'
              fontWeight='semibold'
              noOfLines={1}
            >
              {event.name}
            </Text>
            <Flex gap='1'>
              <IconButton
                size='sm'
                variant='ghost'
                aria-label='timeline link'
                icon={<Icons.Link3 color='gray.500' />}
              />
              <IconButton
                size='sm'
                variant='ghost'
                aria-label='close'
                onClick={closeModal}
                icon={<Icons.XClose color='gray.500' />}
              />
            </Flex>
          </CardHeader>
          <CardBody as={Flex} justify='space-between' px='6' py='0'>
            <VStack w='full' align='flex-start' spacing='4' mb='8'>
              <Flex flexDir='column'>
                <Text fontSize='sm' fontWeight='semibold' color='gray.700'>
                  When
                </Text>
                <Text fontSize='sm' color='gray.700'>
                  {when}
                </Text>
              </Flex>

              <Flex flexDir='column'>
                <Text fontSize='sm' fontWeight='semibold' color='gray.700'>
                  With
                </Text>
                <VStack spacing='0' align='flex-start'>
                  {owner && (
                    <Text fontSize='sm' color='gray.700'>
                      {owner}
                      <Text as='span' color='gray.500'>
                        {` • Organizer`}
                      </Text>
                    </Text>
                  )}
                  {participants?.map((participant, i) => (
                    <Text
                      fontSize='sm'
                      color='gray.700'
                      key={`${i}-${participant}`}
                    >
                      {participant}
                    </Text>
                  ))}
                </VStack>
              </Flex>

              <Flex flexDir='column'>
                <Text fontSize='sm' fontWeight='semibold' color='gray.700'>
                  Agenda
                </Text>
                <Text fontSize='sm' color='gray.700'>
                  {event?.agenda}
                </Text>
              </Flex>

              <FormAutoresizeTextarea
                size='sm'
                name='note'
                label='Notes'
                isLabelVisible
                formId='meeting-notes'
                placeholder='Write some notes from this meeting'
              />
            </VStack>

            <Center minW='12' h='10'>
              <MeetingIcon />
              <Text
                position='absolute'
                fontSize='xl'
                fontWeight='semibold'
                mt='4px'
                color='gray.700'
              >
                {new Date(event?.startedAt).getDate()}
              </Text>
            </Center>
          </CardBody>
          <CardFooter p='6' pt='4'>
            {/* <Flex>
              <Icons.GOOGLE_CALENDAR boxSize='5' />
              <Text fontSize='sm' color='primary.700' ml='2'>
                View in Google Calendar
              </Text>
            </Flex> */}
          </CardFooter>
        </Card>
      </ScaleFade>
    </Flex>
  );
};
