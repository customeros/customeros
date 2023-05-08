import React from 'react';
import Image from 'next/image';
import styles from './timeline-status.module.scss';

interface Props {
  status: 'no-activity' | 'timeline-error';
}
export const TimelineStatus: React.FC<Props> = ({ status }) => {
  return (
    <div className={styles.contentWrapper}>
      <div className={styles.noActivity}>
        <Image
          alt=''
          src={`/backgrounds/blueprint/${status}.webp`}
          fill
          style={{
            objectFit: 'cover',
          }}
          placeholder={'blur'}
          blurDataURL={`/backgrounds/blueprint/${status}-blur.webp`}
        />
      </div>

      <h1 className={styles.message}>
        {status === 'no-activity'
          ? 'No activity logged yet'
          : 'Oops! Something went wrong while loading the timeline'}
      </h1>
    </div>
  );
};
