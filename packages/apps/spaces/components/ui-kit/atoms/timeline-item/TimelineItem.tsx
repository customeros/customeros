import React, { useCallback } from 'react';
import styles from './timeline-item.module.scss';
import { DateTimeUtils } from '../../../../utils';
import Image from 'next/image';
import { DataSource } from '../../../../graphQL/__generated__/generated';
import { motion } from 'framer-motion';

interface Props {
  children: React.ReactNode;
  createdAt?: string | number;
  first?: boolean;
  contentClassName?: any;
  hideTimeTick?: boolean;
  source: string;
}

export const TimelineItem: React.FC<Props> = ({
  children,
  createdAt,
  first,
  contentClassName,
  hideTimeTick,
  source = '',
  ...rest
}) => {
  const getSourceLogo = useCallback(() => {
    if (source === DataSource.ZendeskSupport) return 'zendesksupport';
    if (source === DataSource.Hubspot) return 'hubspot';
    return 'openline_small';
  }, [source]);

  return (
    <motion.div
      initial='hidden'
      whileInView='visible'
      viewport={{ once: true }}
      className={`${styles.timelineItem}`}
      transition={{ duration: 0.3, delay: 0.1 }}
      variants={{
        visible: { opacity: 1, scale: 1 },
        hidden: { opacity: 0, scale: 0 },
      }}
    >
      {!hideTimeTick && (
        <>
          {createdAt ? (
            <div className={styles.when}>
              <div className={styles.timeAgo}>
                {DateTimeUtils.timeAgo(createdAt, {
                  addSuffix: true,
                })}
              </div>
              <div className={styles.metadata}>
                {DateTimeUtils.format(createdAt)}{' '}
                {!!source.length && (
                  <div
                    className={styles.sourceLogo}
                    data-tooltip={`From ${source.toLowerCase()}`}
                  >
                    <Image
                      className={styles.logo}
                      src={`/logos/${getSourceLogo()}.svg`}
                      alt={source}
                      height={16}
                      width={16}
                    />
                  </div>
                )}
              </div>
            </div>
          ) : (
            'Date not available'
          )}
        </>
      )}

      <div className={`${styles.content} ${contentClassName}`} {...rest}>
        {children}
      </div>
    </motion.div>
  );
};
