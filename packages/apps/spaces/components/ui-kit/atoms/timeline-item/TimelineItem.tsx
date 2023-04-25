import React from 'react';
import styles from './timeline-item.module.scss';
import { DateTimeUtils } from '../../../../utils';

interface Props {
  children: React.ReactNode;
  createdAt?: string | number;
  first?: boolean;
  contentClassName?: any;
  hideTimeTick?: boolean;
}

export const TimelineItem: React.FC<Props> = ({
  children,
  createdAt,
  first,
  contentClassName,
  hideTimeTick,
  ...rest
}) => {
  return (
    <div className={`${styles.timelineItem}`}>
      {!hideTimeTick && (
        <>
          {createdAt ? (
            <div className={styles.when}>
              <div className={styles.timeAgo}>
                {DateTimeUtils.timeAgo(createdAt, {
                  addSuffix: true,
                })}
              </div>
              {DateTimeUtils.format(createdAt)}
            </div>
          ) : (
            'Date not available'
          )}
        </>
      )}

      <div className={`${styles.content} ${contentClassName}`} {...rest}>
        {children}
      </div>
    </div>
  );
};
