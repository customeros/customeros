import React from 'react';
import styles from './timeline-item.module.scss';
import ReactTimeAgo from 'react-time-ago';
import Moment from 'react-moment';

interface Props {
  children: React.ReactNode;
  createdAt?: string | number;
  fistOrLast?: boolean;
  style?: any;
}

export const TimelineItem: React.FC<Props> = ({
  children,
  createdAt,
  fistOrLast,
  ...rest
}) => {
  return (
    <div className={`${styles.timelineItem}`}>
      {!fistOrLast ? <span className={styles.timelineLine} /> : null}
      {createdAt ? (
        <>
          <ReactTimeAgo
            className='text-sm text-gray-500 mb-1'
            date={new Date(createdAt)}
            locale='en-US'
          />
          <Moment
            className='text-sm text-gray-500'
            date={createdAt}
            format={'D-M-YYYY h:mm A'}
          ></Moment>
        </>
      ) : (
        'Date not available'
      )}
      <span className={styles.timelineLine} />
      <div className={styles.content} {...rest}>
        {children}
      </div>
    </div>
  );
};
