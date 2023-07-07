import React, { useCallback, useEffect, useRef, useState } from 'react';
import { buckets, events } from '@spaces/organization/organization-timeline/data';
import styles from './organization-timeline.module.scss';
import { DateTimeUtils } from '@spaces/utils/date';
import { Virtuoso, VirtuosoHandle } from "react-virtuoso";
import { Skeleton } from "@spaces/atoms/skeleton";

const itemSize = 102;
const START_INDEX = events.length - 1
const INITIAL_ITEM_COUNT = events.length

const Bucket = ({ e }) => (
    <div key={e.timeBucket} style={{ height: `${(e.count * itemSize)}px`,  }}>
        <span
            style={{ whiteSpace: 'nowrap', }}
            className={styles.date}
        >
          {DateTimeUtils.format(e.timeBucket, DateTimeUtils.defaultFormatShortString)}
        </span>
        {Array(e.count).fill('').map((_, i) => (
            <div className={styles.marker} key={`timeline-type-marker-${i}`} />
        ))}
    </div>
)

const TimelineItem = ({ user }) => (
    <div style={{ padding: '22px 0 0' }}>
        <div style={{ padding: '1rem 0.5rem', background: '#f2fafa' }} className={styles.abc}>
            <h4>
                {DateTimeUtils.format(user.date, DateTimeUtils.defaultFormatShortString)}
            </h4>
            <div>
                {user.content.from} -> {user.content.to}
            </div>
        </div>
    </div>
)

export const OrganizationTimeline = ({ id }) => {
    const otherRef = useRef<any>(null);
    const ref = useRef<VirtuosoHandle>(null)
    const [firstItemIndex, setFirstItemIndex] = useState(START_INDEX)
    const [users, setUsers] = useState(events)


    useEffect(() => {
        if (otherRef.current) {
            otherRef.current.scrollTop = otherRef.current.scrollHeight
        }
    }, [otherRef])

    return (
        <div style={{ display: 'flex' }}>
            <div
                ref={otherRef}
                style={{ maxHeight: 800, overflowY: 'hidden' }}
                className={styles.line}
            >
                {buckets.map((e, i) => <Bucket e={e} />)}
            </div>
            <Virtuoso
                ref={ref}
                onScroll={(e) => {
                    if (otherRef?.current) {
                        otherRef.current.scrollTop = e.target.scrollTop
                    }
                }}
                style={{ height: 800, width: '100%' }}
                firstItemIndex={firstItemIndex}
                initialItemCount={events.length}
                initialTopMostItemIndex={INITIAL_ITEM_COUNT - 1}
                data={users}
                increaseViewportBy={300}
                overscan={10}
                atTopThreshold={100}
                itemContent={(index, user) => <TimelineItem user={user} />} // render TimelineItem
            />
        </div>
    );
};