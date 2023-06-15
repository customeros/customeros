import React, { useEffect, useRef, useState } from 'react';
import {
  NoActivityTimelineElement,
  TimelineItemSkeleton,
} from '@spaces/atoms/timeline-item';
import { useInfiniteScroll } from './useInfiniteScroll';
import classNames from 'classnames';
import styles from './timeline.module.scss';
import { useTimeline } from '@spaces/organisms/timeline/context/useTimeline';
import { TimelineItemByType } from '@spaces/organisms/timeline/TimelineItemByType';
import { LiveEventTimelineItem } from '@spaces/molecules/live-event-timeline-item';

interface Props {
  loading: boolean;
  noActivity: boolean;
  loggedActivities: Array<any>;
  onLoadMore: (ref: any) => void;
  id?: string;
  mode: 'CONTACT' | 'ORGANIZATION';
  contactName?: string;
}

export const Timeline = ({
  loading,
  noActivity,
  loggedActivities,
  onLoadMore,
  id,
  contactName,
}: Props) => {
  const [useAnchoring, setUseAnchoring] = useState(true);
  const anchor = useRef<HTMLDivElement>(null);
  const { timelineContainerRef, onScrollToBottom } = useTimeline();

  const infiniteScrollElementRef = useRef(null);
  useInfiniteScroll({
    element: infiniteScrollElementRef,
    isFetching: loading,
    callback: () => {
      if (loggedActivities.length > 10 && !loading && !useAnchoring) {
        onLoadMore(timelineContainerRef);
      }
    },
  });

  useEffect(() => {
    if (useAnchoring && !loading && loggedActivities.length) {
      setTimeout(() => {
        onScrollToBottom();
        return setUseAnchoring(false);
      }, 100);
    }
  }, [loading, loggedActivities, onScrollToBottom, useAnchoring]);
  return (
    <div ref={timelineContainerRef} className={styles.timeline}>
      <div className={classNames(styles.timelineContent, styles.scrollable)}>
        <div
          ref={infiniteScrollElementRef}
          style={{
            height: '1px',
            width: '100%',
            display: useAnchoring ? 'none' : 'block',
          }}
        />
        {loading && (
          <>
            <TimelineItemSkeleton key='timeline-element-skeleton-1' />
            <TimelineItemSkeleton key='timeline-element-skeleton-2' />
          </>
        )}
        {noActivity && (
          <NoActivityTimelineElement key='no-activity-timeline-item' />
        )}

        {loggedActivities.map((e: any, index) => {
          return (
            <TimelineItemByType
              key={`${e.__typename}-${e.id}-${index}-timeline-element`}
              type={e.__typename}
              data={e}
              index={index}
              loggedActivities={loggedActivities}
              mode='CONTACT'
              contactName={contactName || ''}
              id={e.id}
            />
          );
        })}
        <LiveEventTimelineItem
          key='live-stream-timeline-item'
          first={false}
          contactId={id}
          source={'LiveStream'}
        />
        <div
          className={styles.scrollAnchor}
          ref={anchor}
          key='chat-scroll-timeline-anchor'
        />
      </div>
    </div>
  );
};
