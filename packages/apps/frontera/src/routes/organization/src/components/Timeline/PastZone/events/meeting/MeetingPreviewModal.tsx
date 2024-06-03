import { Link } from 'react-router-dom';
import { useForm } from 'react-inverted-form';

import { convert } from 'html-to-text';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Meeting, ExternalSystemType } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';
import {
  CardHeader,
  CardFooter,
  CardContent,
} from '@ui/presentation/Card/Card';
import { useUpdateMeetingMutation } from '@organization/graphql/updateMeeting.generated';
import { useAddMeetingNoteMutation } from '@organization/graphql/addMeetingNote.generated';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { CalcomIcon, MeetingIcon, HubspotIcon } from './icons';
import { getParticipantEmail } from '../../../../../hooks/utils';

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
      <CardHeader className='sticky top-0 pt-4 p-4 pb-1 rounded-xl flex justify-between gap-4 items-center'>
        <p className='text-gray-700 font-semibold text-lg line-clamp-1'>
          {event.name}
        </p>
        <div className='flex gap-1'>
          <Tooltip label='Copy link'>
            <IconButton
              size='xs'
              variant='ghost'
              aria-label='copy link'
              icon={<Link03 className='text-gray-500' />}
              onClick={() => copy(window.location.href)}
            />
          </Tooltip>
          <Tooltip label='Close'>
            <IconButton
              size='xs'
              variant='ghost'
              aria-label='close'
              onClick={closeModal}
              icon={<XClose className='text-gray-500' />}
            />
          </Tooltip>
        </div>
      </CardHeader>
      <CardContent className='flex justify-between px-6 py-0 overflow-auto max-h-[calc(100vh-4rem-56px-51px-16px-16px)]'>
        <div className='flex flex-col w-full items-start space-y-4'>
          <div className='flex flex-row items-center w-full'>
            <div className='flex flex-grow flex-col'>
              <p className='text-sm font-semibold text-gray-700'>When</p>
              <Tooltip
                label={`Organizer's time: ${
                  creatorTimeZone ? zoned : 'unknown'
                }`}
                side='bottom'
                className='text-xs'
              >
                <p className='text-sm text-gray-700 w-full'>{when}</p>
              </Tooltip>
            </div>

            <div className='flex items-center justify-center min-w-12 h-10 relative'>
              <MeetingIcon />
              <p className='absolute mt-1 text-xl font-semibold text-gray-700'>
                {new Date(event?.startedAt).getDate()}
              </p>
            </div>
          </div>

          <div className='flex flex-col'>
            <p className='text-sm font-semibold text-gray-700'>With</p>
            <div className='flex flex-col space-y-0 items-start'>
              {owner && (
                <p className='text-sm text-gray-700'>
                  {owner}
                  <span className='text-gray-500'>{` • Organizer`}</span>
                </p>
              )}
              {participants?.map((participant, i) => (
                <p
                  className='text-sm text-gray-700'
                  key={`${i}-${participant}`}
                >
                  {participant}
                </p>
              ))}
            </div>
          </div>

          {event?.agenda && (
            <div className='flex flex-col'>
              <p className='text-sm font-semibold text-gray-700'>Description</p>
              <p
                className={cn(
                  event?.agenda ? 'text-gray-700' : 'text-gray-500',
                  'text-sm',
                )}
              >
                {event?.agenda
                  ? convert(event?.agenda ?? '', { preserveNewlines: true })
                  : 'No description was added'}
              </p>
            </div>
          )}

          <FormAutoresizeTextarea
            size='sm'
            name='note'
            label='Notes'
            // isLabelVisible
            formId='meeting-notes'
            placeholder='Write some notes from this meeting'
          />
        </div>
      </CardContent>
      <CardFooter className='p-6 pt-0'>
        {externalSystem && externalUrl && (
          <div className='flex pt-4'>
            {ExternalSystemIcon}
            <Link
              className='text-sm text-primary-700 ml-2'
              to={externalUrl}
              target='_blank'
            >
              {`View in ${externalSystemLabel}`}
            </Link>
          </div>
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
