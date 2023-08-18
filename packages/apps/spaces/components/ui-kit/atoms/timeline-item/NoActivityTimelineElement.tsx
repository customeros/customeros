import React from 'react';
import styles from './timeline-item.module.scss';
import Image from 'next/image';
import { motion } from 'framer-motion';

export const NoActivityTimelineElement: React.FC = () => {
  return (
    <motion.div
      initial='hidden'
      whileInView='visible'
      viewport={{ once: true }}
      className={`${styles.timelineItem}`}
      transition={{ duration: 0.3 }}
      variants={{
        visible: { opacity: 1 },
        hidden: { opacity: 0 },
      }}
    >
      <div className={styles.when}>
        <span className={styles.emptyTimeline}>No activity logged yet</span>
        <div className={styles.metadata}>
          <span className={styles.emptyPlaceholder} />
          <div className={styles.sourceLogo}>
            <Image
              className={styles.logo}
              src={`/logos/customer-os.png`}
              alt='Openline'
              height={16}
              width={16}
            />
          </div>
        </div>
      </div>
    </motion.div>
  );
};
