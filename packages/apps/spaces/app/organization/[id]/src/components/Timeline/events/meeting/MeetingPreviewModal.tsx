'use client';

import { convert } from 'html-to-text';
import { useForm } from 'react-inverted-form';

import { Icons } from '@ui/media/Icon';
import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Center } from '@ui/layout/Center';
import { Text } from '@ui/typography/Text';
import { Link } from '@ui/navigation/Link';
import { Tooltip } from '@ui/overlay/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { DateTimeUtils } from '@spaces/utils/date';
import { ExternalSystemType, Meeting } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { CardBody, CardHeader, CardFooter } from '@ui/presentation/Card';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';
import { useUpdateMeetingMutation } from '@organization/src/graphql/updateMeeting.generated';
import { useAddMeetingNoteMutation } from '@organization/src/graphql/addMeetingNote.generated';

import {
  useTimelineEventPreviewMethodsContext,
  useTimelineEventPreviewStateContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
import { getParticipantEmail } from '../utils';
import { MeetingIcon, HubspotIcon, CalcomIcon } from './icons';

interface MeetingPreviewModalProps {
  invalidateQuery: () => void;
}

export const MeetingPreviewModal = ({
  invalidateQuery,
}: MeetingPreviewModalProps) => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const [_, copy] = useCopyToClipboard();
  const client = getGraphQLClient();
  const updateMeeting = useUpdateMeetingMutation(client, {
    onSuccess: invalidateQuery,
  });
  const addMeetingNote = useAddMeetingNoteMutation(client, {
    onSuccess: invalidateQuery,
  });

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

  const zoned = (() => {
    if (!event?.startedAt) return null;

    const start = DateTimeUtils.convertToTimeZone(
      event?.startedAt,
      !event?.endedAt
        ? DateTimeUtils.timeWithGMT
        : DateTimeUtils.usaTimeFormatString,
      creatorTimeZoneCode,
    );

    if (!event?.endedAt) return start;

    const end = DateTimeUtils.convertToTimeZone(
      event?.endedAt,
      DateTimeUtils.timeWithGMT,
      creatorTimeZoneCode,
    );

    return `${start} - ${end}`;
  })();

  const owner = getParticipantEmail(event?.createdBy?.[0]);

  const participants = event?.attendedBy
    ?.map(getParticipantEmail)
    .filter((c) => c !== owner);

  const externalSystem = event?.externalSystem?.[0]?.type;
  const externalUrl = event?.externalSystem?.[0]?.externalUrl;
  const [externalSystemLabel, ExternalSystemIcon] =
    getExternalSystemAssets(externalSystem);

  useForm<{ note: string }>({
    formId: 'meeting-notes',
    defaultValues: {
      note: convert(event?.note?.[0]?.content ?? '', {
        preserveNewlines: true,
      }),
    },
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        const { name, value } = action.payload;
        if (name === 'note') {
          if (!event?.note.length) {
            addMeetingNote.mutate({
              meetingId: event?.id,
              note: { content: value },
            });
          } else {
            updateMeeting.mutate({
              meetingId: event?.id,
              meeting: {
                appSource: event?.appSource ?? '',
                note: { id: event?.note?.[0]?.id, content: value },
              },
            });
          }
        }
      }
      return next;
    },
  });

  return (
    <>
      <CardHeader
        position='sticky'
        top={0}
        pt={4}
        p='4'
        pb='1'
        borderRadius='xl'
        as={Flex}
        justify='space-between'
        gap='4'
        align='center'
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
          <Flex flexDir='row' alignItems={'center'} w={'100%'}>
            <Flex flexGrow={1} flexDir='column'>
              <Text fontSize='sm' fontWeight='semibold' color='gray.700'>
                When
              </Text>
              <Tooltip
                label={`Organizer's time: ${
                  creatorTimeZone ? zoned : 'unknown'
                }`}
                fontSize='xs'
                placement='bottom'
              >
                <Text fontSize='sm' color='gray.700' w='full'>
                  {when}
                </Text>
              </Tooltip>
            </Flex>

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

          {event?.agenda && (
            <Flex flexDir='column' w={'inherit'}>
              <Text fontSize='sm' fontWeight='semibold' color='gray.700'>
                Description
              </Text>
              <Text
                fontSize='sm'
                color={event?.agenda ? 'gray.700' : 'gray.500'}
              >
                {event?.agenda
                  ? convert(event?.agenda ?? '', { preserveNewlines: true })
                  : 'No description was added'}
              </Text>
            </Flex>
          )}

          <FormAutoresizeTextarea
            size='sm'
            name='note'
            label='Notes'
            isLabelVisible
            formId='meeting-notes'
            placeholder='Write some notes from this meeting'
          />
        </VStack>
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
    </>
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
