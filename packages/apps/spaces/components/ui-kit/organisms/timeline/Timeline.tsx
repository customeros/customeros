import React, { useRef } from 'react';
import { useStickyScroll } from '../../../../hooks/useStickyScroll';

import {
  PhoneConversationTimelineItem,
  EmailTimelineItem,
  LiveConversationTimelineItem,
  NoteTimelineItem,
  WebActionTimelineItem,
} from '../../molecules';
import { TimelineItem } from '../../atoms/timeline-item';
import { TicketTimelineItem } from '../../molecules/ticket-timeline-item';
import styles from './timeline.module.scss';
import { InteractionTimelineItem } from '../../molecules/interaction-timeline-item';
import { ChatTimelineItem } from '../../molecules/conversation-timeline-item/ChatTimelineItem';
import { useInfiniteScroll } from './useInfiniteScroll';
import { Skeleton } from '../../atoms/skeleton';
import { TimelineStatus } from './timeline-status';

interface Props {
  loading: boolean;
  noActivity: boolean;
  id?: string;
  loggedActivities: Array<any>;
  notifyChange?: (id: any) => void;
  onLoadMore: (ref: any) => void;
  contactName?: string;
  mode: 'CONTACT' | 'ORGANIZATION';
}

export const Timeline = ({
  loading,
  noActivity,
  loggedActivities,
  id,
  onLoadMore,
  contactName = '',
  mode,
}: Props) => {
  const timelineContainerRef = useRef(null);
  const containerRef = useRef(null);
  const infiniteScrollElementRef = useRef(null);
  // @ts-expect-error revisit later
  useStickyScroll(containerRef, loggedActivities || []);
  useInfiniteScroll({
    element: infiniteScrollElementRef,
    isFetching: loading,
    callback: () => {
      if (loggedActivities.length > 10) {
        onLoadMore(containerRef);
      }
    },
  });

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
              contactId={id}
            />
          </TimelineItem>
        );
      case 'Conversation':
        // TODO move to interaction event once we have the data in backend
        if (data.channel === 'WEB_CHAT') {
          return (
            <ChatTimelineItem
              first={index == 0}
              feedId={data.id}
              source={data.source}
              createdAt={data?.startedAt}
              feedInitiator={{
                firstName: data.initiatorFirstName,
                lastName: data.initiatorLastName,
                phoneNumber: data.initiatorUsername.identifier,
                lastTimestamp: data.lastTimestamp,
              }}
            />
          );
        }
        if (data.channel === 'EMAIL') {
          return '';
        }
        // TODO move to interaction event once we have the data in backend
        if (data.channel === 'VOICE') {
          return (
            <PhoneConversationTimelineItem
              first={index == 0}
              feedId={data.id}
              source={data.source}
              createdAt={data?.startedAt}
              feedInitiator={{
                firstName: data.initiatorFirstName,
                lastName: data.initiatorLastName,
                phoneNumber: data.initiatorUsername.identifier,
                lastTimestamp: data.lastTimestamp,
              }}
            />
          );
        }
        return null;

      case 'LiveConversation':
        return (
          <LiveConversationTimelineItem
            first={index == 0}
            contactId={id}
            source={data.source}
          />
        );
      case 'PageView':
        return (
          <TimelineItem first={index == 0} createdAt={data?.startedAt}>
            <WebActionTimelineItem {...data} contactName={contactName} />
          </TimelineItem>
        );
      case 'InteractionSession':
        return (
          <TimelineItem first={index == 0} createdAt={data?.startedAt}>
            <InteractionTimelineItem
              {...data}
              contactId={contactName && id}
              organizationId={!contactName && id}
            />
          </TimelineItem>
        );
      case 'Ticket':
        return (
          <TimelineItem first={index == 0} createdAt={data?.createdAt}>
            <TicketTimelineItem {...data} />
          </TimelineItem>
        );

      case 'InteractionEvent':
        if (data.channel === 'EMAIL') {
          return (
            <TimelineItem first={index == 0} createdAt={data?.createdAt}>
              <EmailTimelineItem
                {...data}
                contactId={mode === 'CONTACT' && id}
              />
            </TimelineItem>
          );
        } else {
          return (
            <div>
              Sorry, looks like &apos;{type}&apos; activity type is not
              supported yet{' '}
            </div>
          );
        }
        return null;

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
    <div ref={timelineContainerRef} className={styles.timeline}>
      {!loading && noActivity && <TimelineStatus status='no-activity' />}
      <div className={styles.timelineContent} ref={containerRef}>
        {!!loggedActivities.length && (
          <div
            ref={infiniteScrollElementRef}
            style={{
              height: '6px',
              width: '6px',
            }}
          />
        )}
        {loading && (
          <div className='flex flex-column mt-4'>
            <Skeleton height={'40px'} className='mb-3' />
            <Skeleton height={'40px'} className='mb-3' />
            <Skeleton height={'40px'} className='mb-3' />
            <Skeleton height={'40px'} className='mb-3' />
            <Skeleton height={'40px'} className='mb-3' />
          </div>
        )}

        {loggedActivities.map((e: any, index) => (
          <React.Fragment key={e.id}>
            {getTimelineItemByType(e.__typename, e, index)}
          </React.Fragment>
        ))}
      </div>
    </div>
  );
};
