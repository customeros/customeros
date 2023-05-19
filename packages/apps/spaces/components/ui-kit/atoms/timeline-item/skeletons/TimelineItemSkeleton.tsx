import React from 'react';
import styles from '@spaces/atoms/timeline-item/timeline-item.module.scss';
import { Skeleton } from '@spaces/atoms/skeleton';
const easing = [0.6, -0.05, 0.01, 0.99];
import { motion } from 'framer-motion';

const fadeInUp = {
  initial: {
    x: 100,
    opacity: 0,
    transition: { duration: 0.6, ease: easing },
  },
  exit: { opacity: 0, transition: { duration: 0.6, ease: easing } },
  animate: {
    x: 0,
    opacity: 1,
    transition: {
      duration: 0.6,
      ease: easing,
    },
  },
};
export const TimelineItemSkeleton: React.FC<{ key?: string }> = ({ key }) => {
  return (
    <motion.div
      variants={fadeInUp}
      className={`${styles.timelineItem}`}
      key={key}
    >
      <>
        <div className={styles.when}>
          <div className={styles.timeAgo}>
            <Skeleton height='12px' width='90px' />
          </div>
          <div className={styles.metadata}>
            <Skeleton height='12px' width='140px' />
            <div className={styles.sourceLogo}>
              <Skeleton
                height={'16px'}
                width={'16px'}
                className={styles.logo}
              />
            </div>
          </div>
        </div>
      </>

      <div className={`${styles.content}`}>
        <Skeleton
          height='120px'
          width='100%'
          className={styles.timelineItemSkeletonContent}
        />
      </div>
    </motion.div>
  );
};
