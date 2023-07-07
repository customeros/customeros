'use client';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import {
  buckets,
  events,
} from '@spaces/organization/organization-timeline/data';

import styles from './organization-timeline.module.scss';
import { DateTimeUtils } from '@spaces/utils/date';
import { Virtuoso, VirtuosoHandle } from 'react-virtuoso';
import { Skeleton } from '@spaces/atoms/skeleton';
import { Card, CardHeader, CardBody } from '@spaces/ui/layout/Card';
import { Text } from '@spaces/ui/typography/Text';
import { Heading } from '@spaces/ui/typography/Heading';

const itemSize = 160;
const START_INDEX = events.length - 1;
const INITIAL_ITEM_COUNT = events.length;

const Bucket = ({ e }) => (
  <div key={e.timeBucket} style={{ height: `${e.count * itemSize}px` }}>
    <span style={{ whiteSpace: 'nowrap' }} className={styles.date}>
      {DateTimeUtils.format(
        e.timeBucket,
        DateTimeUtils.defaultFormatShortString,
      )}
    </span>
    {Array(e.count)
      .fill('')
      .map((_, i) => (
        <div className={styles.marker} key={`timeline-type-marker-${i}`} />
      ))}
  </div>
);

const TimelineItem = ({ user, children }) => (
  <div style={{ padding: '22px 0 0' }}>{children}</div>
);

const EmailTimelineItem = ({ content }) => {
  return (
    <Card
      variant='outline'
      size='md'
      padding={'24px'}
      fontSize='14px'
      maxW='549px'
      maxH='138px'
      height={'138px'}
    >
      <CardHeader padding={0}>
        <Text>
          <Text as={'span'} fontWeight={500}>Jonty Knox</Text>{' '}
          <Text as={'span'} color='#6C757D'>emailed</Text>{' '}
          <Text as={'span'} fontWeight={500} color="#1FAA3D" marginRight={2}>Joan Doe</Text> <Text as={'span'} color='#6C757D'>CC:</Text>{' '}
          <Text as={'span'}>jane@doe.com</Text>
        </Text>
        <Text fontWeight='bold'>Annoying production error</Text>
      </CardHeader>
      <CardBody padding={0} overflow={'hidden '}>
        <Text noOfLines={[1, 2]}>
          Magic numbers: There are magic numbers in your component, consider
          declaring them as constants at the top of your component file. Utilize
          CSS Modules: I see that you are using CSS modules but then you are
          also using inline styles, Inline styles are harder to manage and orga
          nize. In this case, it would be better to use the styles object that
          CSS Modules provides. Magic numbers: There are magic numbers in your
          component, consider declaring them as constants at the top of your
          component file. Utilize CSS Modules: I see that you are using CSS
          modules but then you are also using inline styles, Inline styles are
          harder to manage and organize. In this case, it would be better to use
          the styles object that CSS Modules provides. Here is the improved
          version of your code
        </Text>
      </CardBody>
    </Card>
  );
};

export const OrganizationTimeline = ({ id }) => {
  const otherRef = useRef<any>(null);
  const ref = useRef<VirtuosoHandle>(null);
  const [firstItemIndex, setFirstItemIndex] = useState(START_INDEX);
  const [users, setUsers] = useState(events);

  useEffect(() => {
    if (otherRef.current) {
      otherRef.current.scrollTop = otherRef.current.scrollHeight;
    }
  }, [otherRef]);

  return (
    <div style={{ display: 'flex' }} className={styles.container}>
      <div
        ref={otherRef}
        style={{ maxHeight: '100%', overflowY: 'hidden' }}
        className={styles.line}
      >
        {buckets.map((e, i) => (
          <Bucket e={e} key={`bucket${i}`} />
        ))}
      </div>
      <Virtuoso
        ref={ref}
        onScroll={(e) => {
          if (otherRef?.current) {
            otherRef.current.scrollTop = e.target.scrollTop;
          }
        }}
        style={{ height: '100%', width: '100%' }}
        firstItemIndex={firstItemIndex}
        initialItemCount={events.length}
        initialTopMostItemIndex={INITIAL_ITEM_COUNT - 1}
        data={users}
        increaseViewportBy={300}
        overscan={10}
        atTopThreshold={100}
        itemContent={(index, user) => (
          <TimelineItem user={user}>
            <EmailTimelineItem content={user} />
          </TimelineItem>
        )} // render TimelineItem
      />
    </div>
  );
};
