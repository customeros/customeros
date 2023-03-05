import React, { useEffect, useRef } from 'react';
import { Skeleton } from 'primereact/skeleton';
import {
  ConversationTimelineItem,
  LiveConversationTimelineItem,
  NoteTimelineItem,
  WebActionTimelineItem,
} from '../../molecules';
import { TimelineItem } from '../../atoms/timeline-item';
import { uuidv4 } from '../../../../utils';
import { TicketTimelineItem } from '../../molecules/ticket-timeline-item';

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
  const timelineContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (timelineContainerRef?.current && !loading) {
      timelineContainerRef?.current?.scroll({
        top: timelineContainerRef.current.scrollHeight,
        behavior: 'smooth',
      });
    }
  }, [timelineContainerRef?.current, loading]);

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
  const getTimelineItemByType = (type: string, data: any, index: number) => {
    switch (type) {
      case 'Note':
        return (
          <TimelineItem first={index == 0} createdAt={data?.createdAt}>
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
            first={index == 0}
            feedId={data.id}
            source={data.source}
            createdAt={data?.startedAt}
          />
        );
      case 'LiveConversation':
        return (
          <LiveConversationTimelineItem
            first={index == 0}
            contactId={contactId}
            source={data.source}
          />
        );
      case 'PageViewAction':
        return (
          <TimelineItem first={index == 0} createdAt={data?.startedAt}>
            <WebActionTimelineItem {...data} />
          </TimelineItem>
        );
      case 'Ticket':
        return (
          <TimelineItem first={index == 0} createdAt={data?.createdAt}>
            <TicketTimelineItem {...data} />
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
          'AAAA'
        );
    }
  };

  return (
    <div ref={timelineContainerRef}>
      {loggedActivities.map((e: any, index) => (
        <React.Fragment key={uuidv4()}>
          {getTimelineItemByType(e.__typename, e, index)}
        </React.Fragment>
      ))}
    </div>
  );
};
