import React from 'react';
import { Skeleton } from 'primereact/skeleton';
import {
  ConversationTimelineItem,
  NoteTimelineItem,
  WebActionTimelineItem,
} from '../../molecules';
import { TimelineItem } from '../../atoms/timeline-item';
import { uuidv4 } from '../../../../utils';

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
      case 'Note':
        return (
          <TimelineItem
            fistOrLast={loggedActivities.length - 1 === index}
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
      case 'Conversation':
        return (
          <ConversationTimelineItem
            fistOrLast={loggedActivities.length - 1 === index}
            feedId={data.id}
            source={data.source}
            createdAt={data?.startedAt}
          />
        );
      case 'PageViewAction':
        return (
          <TimelineItem
            fistOrLast={loggedActivities.length - 1 === index}
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
    <div className='mb-3'>
      {loggedActivities.map((e: any, index) => (
        <React.Fragment key={uuidv4()}>
          {getTimelineItemByTime(e.__typename, e, index)}
        </React.Fragment>
      ))}
    </div>
  );
};
