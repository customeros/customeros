import { convert } from 'html-to-text';

import { cn } from '@ui/utils/cn';
import { Meeting } from '@graphql/types';
import { File02 } from '@ui/media/icons/File02';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { MeetingIcon } from './icons';
import {
  getParticipants,
  getParticipantName,
} from '../../../../../hooks/utils';

interface MeetingStubProps {
  data: Meeting;
}

export const MeetingStub = ({ data }: MeetingStubProps) => {
  const owner = getParticipantName(data.createdBy[0]);
  const isSentByTenantUser =
    data.createdBy[0]?.__typename === 'UserParticipant';
  const firstParticipant = getParticipantName(data.attendedBy?.[0]);
  const [participants, remaining] = getParticipants(data);
  const { openModal } = useTimelineEventPreviewMethodsContext();

  const note = convert(data?.note?.[0]?.content ?? '', {
    preserveNewlines: true,
  });
  const agenda = convert(data?.agenda ?? '', { preserveNewlines: true });

  return (
    <Card
      onClick={() => openModal(data.id)}
      className={cn(
        isSentByTenantUser ? 'ml-6' : 'ml-0',
        'bg-white max-w-[549px] border border-gray-200 rounded-lg cursor-pointer shadow-xs hover:shadow-md transition-all duration-200 ease-out',
      )}
    >
      <CardContent className='px-3 py-2'>
        <div className='flex w-full justify-between relative gap-3'>
          <div className='flex flex-col items-start max-w-[461px]'>
            <p className='text-sm font-semibold text-gray-700 line-clamp-1'>
              {data?.name ?? '(No title)'}
            </p>
            <div className='flex'>
              <p
                className={cn(
                  note || agenda ? 'line-clamp-1' : 'line-clamp-3',
                  'text-sm text-gray-700 max-w-[463px]',
                )}
              >
                {owner || firstParticipant}{' '}
                <span className='text-gray-500'>met</span> {participants}
              </p>
              {remaining && (
                <span className='text-sm text-gray-500 ml-1 whitespace-nowrap'>
                  {` + ${remaining}`}
                </span>
              )}
            </div>

            {(note || agenda) && (
              <div className='flex items-start max-w-[517px]'>
                {note && <File02 className='size-3 mt-1 mr-1 text-gray-500' />}
                <p className='text-sm text-gray-500 line-clamp-2'>
                  {note || agenda}
                </p>
              </div>
            )}
          </div>

          <div className='flex min-w-[48px] h-[40px] text-3xl items-center justify-center'>
            <MeetingIcon />
            <p className='absolute text-xl font-semibold text-gray-700'>
              {new Date(data?.startedAt).getDate()}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
