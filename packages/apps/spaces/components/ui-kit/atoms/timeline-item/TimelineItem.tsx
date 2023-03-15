import React from 'react';
import styles from './timeline-item.module.scss';
import { DateTimeUtils } from '../../../../utils';

interface Props {
  children: React.ReactNode;
  createdAt?: string | number;
  first?: boolean;
  contentClassName?: any;
}

export const TimelineItem: React.FC<Props> = ({
  children,
  createdAt,
  first,
  contentClassName,
  ...rest
}) => {
  return (
    <div className={`${styles.timelineItem}`}>
      {createdAt ? (
        <div className={styles.when}>
          <div className={styles.timeAgo}>
            {DateTimeUtils.timeAgo(new Date(createdAt), { addSuffix: true })}
          </div>
          {DateTimeUtils.format(new Date(createdAt))}
        </div>
      ) : (
        'Date not available'
      )}
      <div className={`${styles.content} ${contentClassName}`} {...rest}>
        {children}
      </div>
    </div>
  );
};
