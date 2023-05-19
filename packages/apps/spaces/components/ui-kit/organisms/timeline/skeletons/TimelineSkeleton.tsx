import React from 'react';
import styles from '@spaces/organisms/timeline/timeline.module.scss';
import { TimelineItemSkeleton } from '@spaces/atoms/timeline-item';

export const TimelineSkeleton: React.FC = () => {
  return (
    <div className={styles.timeline}>
      <TimelineItemSkeleton />
      <TimelineItemSkeleton />
    </div>
  );
};
