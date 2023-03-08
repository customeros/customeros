import React from 'react';
import styles from './timeline-item.module.scss';
import format from 'date-fns/format';

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
      {!first ? (
        <span className={`${styles.timelineLine} ${styles.first}`} />
      ) : null}
      {createdAt ? (
        <div className={styles.when}>
          {format(new Date(createdAt), 'dd/MM/yyyy h:mm a')}
        </div>
      ) : (
        'Date not available'
      )}
      <span className={`${styles.timelineLine} ${styles.second}`} />
      <div className={`${styles.content} ${contentClassName}`} {...rest}>
        {children}
      </div>
    </div>
  );
};
