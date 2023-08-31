'use client';

import { useState, useEffect } from 'react';
import { convert } from 'html-to-text';
import { useForm } from 'react-inverted-form';
import { utcToZonedTime, format } from 'date-fns-tz';

import { Icons } from '@ui/media/Icon';
import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Center } from '@ui/layout/Center';
import { Text } from '@ui/typography/Text';
import { Link } from '@ui/navigation/Link';
import { Tooltip } from '@ui/overlay/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { DateTimeUtils } from '@spaces/utils/date';
import { ScaleFade } from '@ui/transitions/ScaleFade';
import { ExternalSystemType, Meeting } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@spaces/hooks/useCopyToClipboard';
import { Card, CardBody, CardHeader, CardFooter } from '@ui/presentation/Card';
// import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';
import { useUpdateMeetingMutation } from '@organization/graphql/updateMeeting.generated';

import { useTimelineEventPreviewContext } from '../../preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';
import { MeetingIcon, HubspotIcon, CalcomIcon } from './icons';

export const MeetingPreviewModal = () => {
  const [_, copy] = useCopyToClipboard();
  const client = getGraphQLClient();
  const [isMounted, setIsMounted] = useState(false);
  const updateMeeting = useUpdateMeetingMutation(client);
  const { isModalOpen, closeModal, modalContent } =
    useTimelineEventPreviewContext();

  const event = modalContent as Meeting;
  const creatorTimeZone =
    event?.createdBy?.[0]?.__typename === 'ContactParticipant'
      ? event?.createdBy?.[0]?.contactParticipant?.timezone
      : '';

  const creatorTimeZoneCode = creatorTimeZone
    ? creatorTimeZone?.split(' ')?.[0]
    : '';

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

  const zonedStartDateStr = creatorTimeZone
    ? utcToZonedTime(event?.startedAt, creatorTimeZoneCode)
    : null;
  const zonedEndDateStr = creatorTimeZone
    ? utcToZonedTime(event?.endedAt, creatorTimeZoneCode)
    : null;

  const zoned = (() => {
    if (!zonedStartDateStr) return null;
    const start = format(
      zonedStartDateStr,
      !zonedEndDateStr
        ? DateTimeUtils.timeWithGMT
        : DateTimeUtils.usaTimeFormatString,
      {
        timeZone: creatorTimeZoneCode || undefined,
      },
    );

    if (!zonedEndDateStr) return start;

    const end = format(zonedEndDateStr, DateTimeUtils.timeWithGMT, {
      timeZone: creatorTimeZoneCode || undefined,
    });

    return `${start} - ${end}`;
  })();

  const owner = (() => {
    if (event?.createdBy?.[0]?.__typename === 'ContactParticipant') {
      const participant = event?.createdBy?.[0]?.contactParticipant;
      return participant.emails?.[0]?.email ?? participant?.firstName;
    }
    return '';
  })();

  const participants = event?.attendedBy
    ?.map(
      (c) =>
        c.__typename === 'ContactParticipant' &&
        (c.contactParticipant?.emails?.[0]?.email ||
          c.contactParticipant?.firstName ||
          'unknown'),
    )
    .filter((c) => c !== owner)
    .sort((a, b) => {
      if (a === 'unknown') return 1;
      if (b === 'unknown') return -1;
      return 0;
    });

  const externalSystem = event?.externalSystem?.[0]?.type;
  const externalUrl = event?.externalSystem?.[0]?.externalUrl;
  const [externalSystemLabel, ExternalSystemIcon] =
    getExternalSystemAssets(externalSystem);

  useForm<{ note: string }>({
    formId: 'meeting-notes',
    defaultValues: {
      note: convert(event?.note?.[0]?.html ?? ''),
    },
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        const { name, value } = action.payload;
        if (name === 'note') {
          updateMeeting.mutate({
            meetingId: event?.id,
            meeting: {
              appSource: event?.appSource,
              note: { id: event?.note?.[0]?.id, html: value },
            },
          });
        }
      }
      return next;
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
              <Tooltip label='Copy link'>
                <IconButton
                  size='sm'
                  variant='ghost'
                  aria-label='copy link'
                  icon={<Icons.Link3 color='gray.500' />}
                  onClick={() => copy(window.location.href)}
                />
              </Tooltip>
              <Tooltip label='Close'>
                <IconButton
                  size='sm'
                  variant='ghost'
                  aria-label='close'
                  onClick={closeModal}
                  icon={<Icons.XClose color='gray.500' />}
                />
              </Tooltip>
            </Flex>
          </CardHeader>
          <CardBody
            as={Flex}
            justify='space-between'
            px='6'
            py='0'
            overflowY='auto'
            maxH='calc(100vh - 4rem - 56px - 51px - 16px - 16px);'
          >
            <VStack w='full' align='flex-start' spacing='4'>
              <Flex flexDir='column'>
                <Text fontSize='sm' fontWeight='semibold' color='gray.700'>
                  When
                </Text>
                <Tooltip
                  label={`Organizer's time: ${zoned ?? 'unknown'}`}
                  fontSize='xs'
                  placement='bottom'
                >
                  <Text fontSize='sm' color='gray.700' w='full'>
                    {when}
                  </Text>
                </Tooltip>
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
                  {convert(event?.agenda ?? '')}
                </Text>
              </Flex>

              {/* Remove this when we have createNote mutation */}
              <Flex flexDir='column'>
                <Text fontSize='sm' fontWeight='semibold' color='gray.700'>
                  Note
                </Text>
                <Text fontSize='sm' color='gray.700'>
                  {convert(event?.note?.[0]?.html ?? '')}
                </Text>
              </Flex>

              {/* Uncomment this when we have createNote mutation */}
              {/* <FormAutoresizeTextarea
                size='sm'
                name='note'
                label='Notes'
                isLabelVisible
                formId='meeting-notes'
                placeholder='Write some notes from this meeting'
              /> */}
            </VStack>

            <Center minW='12' h='10' position='relative'>
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
          <CardFooter p='6' pt='0'>
            {externalSystem && externalUrl && (
              <Flex pt='4'>
                {ExternalSystemIcon}
                <Link
                  href={externalUrl}
                  fontSize='sm'
                  color='primary.700'
                  ml='2'
                  target='_blank'
                >
                  {`View in ${externalSystemLabel}`}
                </Link>
              </Flex>
            )}
          </CardFooter>
        </Card>
      </ScaleFade>
    </Flex>
  );
};

function getExternalSystemAssets(type: ExternalSystemType) {
  switch (type) {
    case ExternalSystemType.Calcom:
      return ['Calcom', <CalcomIcon key='calcom' />];
    case ExternalSystemType.Hubspot:
      return ['Hubspot', <HubspotIcon key='hubspot' />];
    default:
      return ['Unknown', <></>];
  }
}
