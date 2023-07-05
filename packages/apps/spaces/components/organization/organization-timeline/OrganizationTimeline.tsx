import React, { useCallback, useRef, useState } from 'react';
import {
  buckets,
  events,
} from '@spaces/organization/organization-timeline/data';
// @ts-expect-error
import { FixedSizeList as List } from 'react-window';
// @ts-expect-error
import InfiniteLoader from 'react-window-infinite-loader';
import styles from './organization-timeline.module.scss';
import { DateTimeUtils } from '@spaces/utils/date';
// Assuming all events have same height
const itemSize = 32;

function getEventByIndex(index) {
  // for simplicity, return the single event we have for each item
  // replace this with your actual logic to map an index to an event
  return events[0];
}

export const OrganizationTimeline = ({ id }: { id: string }) => {
  const listRef = useRef<any>(null);
  const otherRef = useRef<any>(null);
  const itemCount = events.length;
  const handleScrollToBottom = useCallback(() => {
    listRef?.current?.scrollTo({
      top: listRef?.current?.scrollHeight,
    });
  }, [listRef]);

  const handleScrubberChange = (event) => {
    const { value } = event.target;
    if (listRef.current) {
      listRef.current.scrollTo({
        top: (value / 100) * itemCount * itemSize,
      });
      otherRef.current.scrollTo({
        top: (value / 100) * itemCount * itemSize,
      });
    }
  };

  return (
    <div
      style={{
        display: 'flex',
      }}
    >
      <div
        ref={otherRef}
        style={{ maxHeight: '90vh', overflowX: 'hidden', overflowY: 'auto' }}
      >
        {buckets.map((e, i) => (
          <div key={e.timeBucket} style={{ height: `${e.count * itemSize}px` }}>
            <span
              style={{
                background: 'paleturquoise',
                fontSize: '10px',
                whiteSpace: 'nowrap',
              }}
            >
              {DateTimeUtils.format(e.timeBucket)}
            </span>
          </div>
        ))}
      </div>
      <div
        style={{ maxHeight: '90vh', overflowX: 'hidden', overflowY: 'auto' }}
        ref={listRef}
        className={styles.div}
      >
        {events
          .map((e, i) => (
            <div
              key={e.date}
              style={{ background: '#ccc', marginBottom: '8px', width: '300px' }}
            >
              <span style={{ marginRight: '8px' }}>{i}</span>

              Events
            </div>
          ))
          .reverse()}
      </div>
      <input
        className={styles.input}
        type='range'
        defaultValue={100}
        min={1}
        max={1}
        orient='vertical'
        aria-orientation={'vertical'}
        onChange={handleScrubberChange}
      />
    </div>
  );
};
