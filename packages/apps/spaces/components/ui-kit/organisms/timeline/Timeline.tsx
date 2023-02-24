import React from 'react';
import { Skeleton } from 'primereact/skeleton';
import {
  ConversationTimelineItem,
  NoteTimelineItem,
  WebActionTimelineItem,
} from '../../molecules';
import { TimelineItem } from '../../atoms/timeline-item';

interface Props {
  loading: boolean;
  noActivity: boolean;
  contactId?: string;
  loggedActivities: Array<any>;
  notifyChange?: (id: any) => void;
  notifyContactNotesUpdate?: (id: any) => void;
}
export const Timeline = ({
  loading,
  noActivity,
  loggedActivities,
  contactId,
  notifyChange = () => null,
  notifyContactNotesUpdate = () => null,
}: Props) => {
  if (loading) {
    return (
      <div className='flex flex-column mt-4'>
        <Skeleton className='mb-3' />
        <Skeleton className='mb-3' />
        <Skeleton className='mb-3' />
        <Skeleton className='mb-3' />
        <Skeleton />
      </div>
    );
  }

  if (!loading && noActivity) {
    return (
      <p className='text-gray-600 font-italic mt-4'>No activity logged yet</p>
    );
  }
  const getTimelineItemByTime = (type: string, data: any, index: number) => {
    switch (type) {
      case 'NOTE':
        return (
          <TimelineItem
            last={loggedActivities.length - 1 === index}
            createdAt={data?.createdAt}
          >
            <NoteTimelineItem
              noteContent={data.html}
              createdAt={data.createdAt}
              createdBy={data?.createdBy}
              id={data.id}
              source={data?.source}
              refreshNoteData={
                data?.contact ? notifyContactNotesUpdate : notifyChange
              }
              contactId={contactId}
            />
          </TimelineItem>
        );
      case 'CONVERSATION':
        return (
          <ConversationTimelineItem
            feedId={data.id}
            source={data.source}
            createdAt={data?.createdAt}
          />
        );
      case 'ACTION':
        return (
          <TimelineItem
            last={loggedActivities.length - 1 === index}
            createdAt={data?.createdAt}
          >
            <WebActionTimelineItem {...data} />
          </TimelineItem>
        );
      // case "CALL":
      //     return <PhoneCallTimelineItem phoneCallParties={data} duration={}/>
      default:
        return type ? (
          <div>
            Sorry, looks like &apos;{type}&apos; activity type is not supported
            yet{' '}
          </div>
        ) : (
          ''
        );
    }
  };

  return (
    <div className='mt-5 mb-3'>
      <div className='text-sm text-gray-500 flex justify-content-center mb-1'>
        {' '}
        Now{' '}
      </div>
      {loggedActivities.map((e: any, index) => (
        <React.Fragment key={e.id}>
          {getTimelineItemByTime(e.type, e, index)}
        </React.Fragment>
      ))}
    </div>
  );
};
